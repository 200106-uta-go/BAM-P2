// Package cluster is a package of functions the assist with setting up a cluster on kops on the linux machine.
// Requires the following packages
//   commander (github.com/200106-uta-go/BAM-P2/pkg/commander) - used to send commands to the OS using "os/exec" standard go package
//   userinputs (github.com/200106-uta-go/BAM-P2/pkg/userinputs) - used to get inputs from the user (stdin)
//   awsinit (github.com/200106-uta-go/BAM-P2/deployments/awsinit) - used to add new kops aws user to aws cli configuration
package cluster

import (
	"fmt"
	"strings"

	"github.com/200106-uta-go/BAM-P2/deployments/awsinit"
	"github.com/200106-uta-go/BAM-P2/pkg/commander"
	"github.com/200106-uta-go/BAM-P2/pkg/filewriter"
	"github.com/200106-uta-go/BAM-P2/pkg/userinputs"
)

// CreateKopsAWSUser creates an aws user with the correct permissions for kops to operate on aws
func CreateKopsAWSUser() {
	// check if user already exist
	test := commander.CmdRunOutSilent("aws iam get-user --user-name kops")
	if test == "" {
		// if not create it
		// create aws group and policies
		fmt.Println("Creating kops AWS group and policies...")
		commander.CmdRunSilent("aws iam create-group --group-name kops")
		commander.CmdRunSilent("aws iam attach-group-policy --policy-arn arn:aws:iam::aws:policy/AmazonEC2FullAccess --group-name kops")
		commander.CmdRunSilent("aws iam attach-group-policy --policy-arn arn:aws:iam::aws:policy/AmazonRoute53FullAccess --group-name kops")
		commander.CmdRunSilent("aws iam attach-group-policy --policy-arn arn:aws:iam::aws:policy/AmazonS3FullAccess --group-name kops")
		commander.CmdRunSilent("aws iam attach-group-policy --policy-arn arn:aws:iam::aws:policy/IAMFullAccess --group-name kops")
		commander.CmdRunSilent("aws iam attach-group-policy --policy-arn arn:aws:iam::aws:policy/AmazonVPCFullAccess --group-name kops")

		// create aws user
		fmt.Println("Creating kops AWS user...")
		commander.CmdRunSilent("aws iam create-user --user-name kops")
		commander.CmdRunSilent("aws iam add-user-to-group --user-name kops --group-name kops")
		out := commander.CmdRunOutSilent("aws iam create-access-key --user-name kops") // save key values here
		userinputs.RequestAnswer("Please take a moment to record the following in a safe place.\n" + out +
			"Please enter any value to continue:")

		// get keys from user creation
		// find keys
		var key, secret string
		tempSlice := strings.Split(out, " ")
		for i := 0; i < len(tempSlice); i++ {
			if tempSlice[i] == "\"AccessKeyId\":" {
				key = tempSlice[i+1]
			}
			if tempSlice[i] == "\"SecretAccessKey\":" {
				secret = tempSlice[i+1]
			}
		}

		// cleamup keys
		key = key[1:len(key)]           // remove first character "
		key = key[:len(key)-3]          // remove last 3 characters " | , | \n
		secret = secret[1:len(secret)]  // remove first character "
		secret = secret[:len(secret)-3] // remove last 3 characters " | , | \n

		// add user to aws cli
		awsinit.AddAWSUserAuto(key, secret, "", "")

		// export keys as env
		filewriter.AppendToEnv("export AWS_ACCESS_KEY_ID=" + key)
		filewriter.AppendToEnv("export AWS_SECRET_ACCESS_KEY=" + secret)

	} else {
		// else let user now it already exist
		fmt.Println("kops user already exist in AWS")
	}
}

// CreateKopsStateStore creates a kops state store in a s3 bucket. Asks user for name of bucket, and if bucket is present
// uses it, otherwise creates bucket
func CreateKopsStateStore() string {
	haveBucket := false
	// get bucket name from user
	bucket := userinputs.RequestAnswer("Please provide a name for your aws s3 bucket (all buckets must be unique):")
	out := commander.CmdRunOutSilent("aws s3 ls")
	outSlice := strings.Split(out, " ")
	for _, o := range outSlice {
		if o == bucket {
			haveBucket = true
			break
		}
	}
	if haveBucket {
		// already have bucket with that name
		fmt.Println("That bucket already exist in your AWS, using that bucket")
	} else {
		// creat bucket
		region := commander.CmdRunOutSilent("aws configure get region")
		region = region[:len(region)-1] // remove last character \n

		// no bucket with that name need to create one
		fmt.Println("Creating aws s3 bucket...")
		commander.CmdRun("aws s3api create-bucket --bucket " + bucket + " --region " + region)
		commander.CmdRun("aws s3api put-bucket-versioning --bucket " + bucket + " --versioning-configuration Status=Enabled")
		commander.CmdRun("aws s3api put-bucket-encryption --bucket " + bucket + " --server-side-encryption-configuration '{\"Rules\":[{\"ApplyServerSideEncryptionByDefault\":{\"SSEAlgorithm\":\"AES256\"}}]}'")
	}

	filewriter.AppendToEnv("export KOPS_STATE_STORE=s3://" + bucket)
	return "s3://" + bucket
}

// PrepairCluster perpairs a cluster for activation
func PrepairCluster() string {
	fmt.Println("Prepairing cluster")
	var clusterName string

	// clust.kopsStateStore
	kobStateStore := CreateKopsStateStore()
	// clust.clusterName
	for {
		clusterName = userinputs.RequestAnswer("Enter cluster name (must be followed with a .k8s.local, ex: yourCluster.k8s.local):")
		if CheckCluster(clusterName) {
			fmt.Println("That cluster already exist, please use a diffrent name")
		} else {
			break
		}
	}
	// clust.region
	region := commander.CmdRunOutSilent("aws configure get region")
	region = region[:len(region)-1] // remove last character \n
	// clust.masterType
	masterType := userinputs.RequestAnswer("What EC2 type do you want for your master (ex: t2.micro. m5.large, etc.)?")
	// clust.masterCount
	masterCount := userinputs.RequestAnswer("How many masters do you want to run in this cluster?")
	// clust.nodeType
	nodeType := userinputs.RequestAnswer("What EC2 type do you want for your nodes (ex: t2.micro. m5.large, etc.)?")
	// clust.nodeCount
	nodeCount := userinputs.RequestAnswer("How many nodes do you want to run in this cluster?")

	fmt.Println("Creating cluster with the following...\n" +
		"  cluster:\n" +
		"    name: " + clusterName + "\n" +
		"    s3 bucket: " + kobStateStore + "\n" +
		"    region: " + region + "\n" +
		"    master type: " + masterType + "\n" +
		"    master count: " + masterCount + "\n" +
		"    node type: " + nodeType + "\n" +
		"    node count: " + nodeCount + "\n")

	commander.CmdRun("kops create cluster --cloud=aws " +
		"--master-zones=" + region + "a --zones=" + region + "a," + region + "b," + region + "c " +
		"--node-count=" + nodeCount + " --node-size=" + nodeType + " " +
		"--master-count=" + masterCount + " --master-size=" + masterType + " " + clusterName)

	return clusterName
}

// Up calls PrepairCluster and then brings it up, activiating any required cloud resources
func Up() {
	// deploy cluster (up)
	CreateKopsAWSUser()
	name := PrepairCluster()
	commander.CmdRun("kops update cluster " + name + " --yes")
}

// Down removes the cluster also removes/deletes any cloud resources currently active
func Down() {
	// destroy cluster (down)
	var name string
	out := commander.CmdRunOut("kops get cluster")
	fmt.Println("Current running clusters...")
	fmt.Println(out)
	for {
		name = userinputs.RequestAnswer("Enter name of cluster you wish to remove:")
		if CheckCluster(name) {
			break
		} else {
			fmt.Println("That cluster does not exist, please use a diffrent name")
		}
	}
	commander.CmdRun("kops delete cluster " + name + " --yes")
}

// CheckCluster checks to see if cluster exists and returns true (bool) if present, otherwise false (bool)
func CheckCluster(name string) bool {
	out := commander.CmdRunOutSilentNoErr("kops get cluster")
	outSlice := strings.Split(out, " ")
	for _, o := range outSlice {
		if o == name {
			return true
		}
	}
	return false
}

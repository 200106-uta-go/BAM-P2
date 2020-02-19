package aws_cli

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/200106-uta-go/BAM-P2/pkg/userinputs"
)

func ReadyToInstall() bool {
	fmt.Println("We will not be able to complete your AWS cli install and setup if you do not have credeintuals.")
	fmt.Println("\n## To create access keys for an IAM user ##")
	fmt.Println("\n1. Sign in to the AWS Management Console and open the IAM console at https://console.aws.amazon.com/iam/")
	fmt.Println("2. In the navigation pane, choose Users.")
	fmt.Println("3. Choose the name of the user whose access keys you want to create, and then choose the Security credentials tab.")
	fmt.Println("4. In the Access keys section, choose Create access key.")
	fmt.Println("5. To view the new access key pair, choose Show. You will not have access to the secret access key again after this")
	fmt.Println("   dialog box closes. Your credentials will look something like this:")
	fmt.Println("      Access key ID: AKIAIOSFODNN7EXAMPLE")
	fmt.Println("      Secret access key: wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY")
	fmt.Println("6. To download the key pair, choose Download .csv file. Store the keys in a secure location. You will not have")
	fmt.Println("   access to the secret access key again after this dialog box closes. Keep the keys confidential in order to protect")
	fmt.Println("   your AWS account and never email them. Do not share them outside your organization, even if an inquiry appears to ")
	fmt.Println("   come from AWS or Amazon.com. No one who legitimately represents Amazon will ever ask you for your secret key.")
	fmt.Println("7. After you download the .csv file, choose Close. When you create an access key, the key pair is active by default,")
	fmt.Println("   and you can use the pair right away.")
	ans := userinputs.RequestAnswer("\nDid you complete the above steps or already have a IAM User key pair? (yes/no)")
	if ans == "yes" || ans == "y" {
		return true
	}
	return false
}

func CheckInstall() {
	// logic to check if asw cli is already installed
	cmdRun("aws --version")\


}

func InstallAWS() {
	cmdRun("curl https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip -o awscliv2.zip")
	cmdRun("unzip awscliv2.zip")
	cmdRun("sudo ./aws/install")

	// clean up
	cmdRun("rm -r aws")
	cmdRun("rm -r awscliv2.zip")
}

func SetupAWS() {
	/*
		logic to setup AWS
			$ aws configure
				AWS Access Key ID [None]: AKIAIOSFODNN7EXAMPLE
				AWS Secret Access Key [None]: wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY
				Default region name [None]: us-west-2
				Default output format [None]: json
	*/
	awsKey := userinputs.RequestAnswer("AWS Access Key ID (Eaxmple: AKIAIOSFODNN7EXAMPLE):")
	awsSecret := userinputs.RequestAnswer("AWS Secret Access Key (Example: wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY):")
	region := userinputs.RequestAnswer("Default region name (Example: us-east-1):")

	cmdRun("aws configure set aws_access_key_id "+awsKey)
	cmdRun("aws configure set aws_secret_access_key "+awsSecret)
	cmdRun("was configure set region "+region)	
}

func cmdRun(s string) {
	fmt.Println("Running command: " + s)
	mySlice := strings.Split(s, " ")
	cmd := exec.Command(mySlice[0], mySlice[1:]...)
	cmd.Stderr = os.Stderr
	out, err := cmd.Output()
	if err != nil {
		log.Println(err)
	}
	fmt.Println(string(out))
}



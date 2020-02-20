// Package awsinit is a package of functions the assist with installing aws cli v2 on the linux machine. Requires the following packages
//   commander (github.com/200106-uta-go/BAM-P2/pkg/commander) - used to send commands to the OS using "os/exec" standard go package
//   userinputs (github.com/200106-uta-go/BAM-P2/pkg/userinputs) - used to get inputs from the user (stdin)
package awsinit

import (
	"fmt"
	"strings"

	"github.com/200106-uta-go/BAM-P2/pkg/commander"
	"github.com/200106-uta-go/BAM-P2/pkg/userinputs"
)

// ReadyToConfig displays instructions to the user about what they need before hand to install and setup aws cli v2
// confirms the user is ready returns true of the user responds with "yes" or "y", otherwise false.
func ReadyToConfig() bool {
	fmt.Println("To be able to complete setup you will need to have IAM credeintuals.")
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
	ans := userinputs.RequestAnswer("\nDid you complete the above steps or already have a IAM User key pair?\n" +
		"Alternativelly you can run \"aws configure\" at anytime to configure your aws cli install.\n" +
		"Continue (y/n): ")
	if ans == "yes" || ans == "y" {
		return true
	}
	return false
}

// CheckInstall simply checks to see if aws cli v2 is installed, returns true if installled, otherwise false.
func CheckInstall() bool {
	// logic to check if asw cli is already installed
	version := commander.CmdRunOutSilent("aws --version")
	v := strings.Split(version, " ")
	if v[0] == "aws-cli/2.0.0" {
		return true
	}
	return false
}

// InstallAWS downloads aws cli v.2 zip from the amazon's server, exstracts it and installs it. After installation the
// installer files are removed for a clean environment.
func InstallAWS() {
	commander.CmdRun("curl https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip -o awscliv2.zip")
	commander.CmdRun("unzip awscliv2.zip")
	commander.CmdRun("sudo ./aws/install")

	// clean up
	commander.CmdRun("rm -r aws")
	commander.CmdRun("rm -r awscliv2.zip")
}

// AddAWSUserM helps configure the aws cli with user input
func AddAWSUserM() {
	awsKey := userinputs.RequestAnswer("AWS Access Key ID (Eaxmple: AKIAIOSFODNN7EXAMPLE):")
	awsSecret := userinputs.RequestAnswer("AWS Secret Access Key (Example: wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY):")
	region := userinputs.RequestAnswer("Default region name (Example: us-east-1):")
	output := userinputs.RequestAnswer("Default output format (json/yaml/text/table):")

	commander.CmdRunSilent("aws configure set aws_access_key_id " + awsKey)
	commander.CmdRunSilent("aws configure set aws_secret_access_key " + awsSecret)
	commander.CmdRunSilent("aws configure set region " + region)
	commander.CmdRunSilent("aws configure set aws_default_output " + output)

	fmt.Println("AWS CLI Configured using provided information!")
}

// AddAWSUserAuto helps configure the aws cli with passed variables
func AddAWSUserAuto(key string, secret string, region string, output string) {
	commander.CmdRunSilent("aws configure set aws_access_key_id " + key)
	commander.CmdRunSilent("aws configure set aws_secret_access_key " + secret)
	commander.CmdRunSilent("aws configure set region " + region)
	commander.CmdRunSilent("aws configure set aws_default_output " + output)

	fmt.Println("AWS CLI Configured using provided information!")
}

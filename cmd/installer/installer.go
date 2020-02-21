package main

import (
	"fmt"
	"strings"

	"github.com/200106-uta-go/BAM-P2/deployments/awsinit"
	"github.com/200106-uta-go/BAM-P2/deployments/kopsinit"
	"github.com/200106-uta-go/BAM-P2/deployments/kubectlinit"
	"github.com/200106-uta-go/BAM-P2/pkg/commander"
)

func main() {
	// check for curl and install if not present
	fmt.Println("Installing curl...")
	out := commander.CmdRunOutSilentNoErr("curl --version")
	outSlice := strings.Split(out, " ")
	if outSlice[0] != "curl" {
		commander.CmdRunSilent("sudo apt install curl -y")
	} else {
		fmt.Println("Curl is already installed!")
	}

	// install and setup aws
	if awsinit.CheckInstall() {
		fmt.Println("AWS CLI is already installed!")
	} else {
		awsinit.InstallAWS()
		if awsinit.ReadyToConfig() {
			awsinit.AddAWSUserM()
		}
	}

	// install kops
	if kopsinit.CheckInstall() {
		fmt.Println("kops is already installed!")
	} else {
		kopsinit.KopsInstall()
	}

	// install kubectl
	if kubectlinit.CheckInstall() {
		fmt.Println("kubectl is already installed!")
	} else {
		kubectlinit.KubectlInstall()
	}
}

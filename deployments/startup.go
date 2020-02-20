package main

import (
	"flag"
	"fmt"

	"github.com/200106-uta-go/BAM-P2/deployments/awsinit"
	"github.com/200106-uta-go/BAM-P2/deployments/cluster"
	"github.com/200106-uta-go/BAM-P2/deployments/kopsinit"
	"github.com/200106-uta-go/BAM-P2/deployments/kubectlinit"
)

var install bool
var deploy bool

func init() {
	// sets flag options
	flag.BoolVar(&install, "install", false, "Launches installer for aws | kops | kubectl") // install prerequisites
	flag.BoolVar(&deploy, "deploy", false, "Starts a cluster")                              // deploy cluster
}

func main() {
	switch {
	case install:
		// run installers
		installer()
	case deploy:
		// deploy cluster
		cluster.Up()
	default:
		installer()
		cluster.Up()
	}
}

// installers (aws | kops | kubectl)
func installer() {
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

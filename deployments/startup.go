package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func main() {
	// install aws
	installAWS()
	// install kops

	// install kubectl
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

func installAWS() {
	cmdRun("curl https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip -o awscliv2.zip")
	cmdRun("unzip awscliv2.zip")
	cmdRun("sudo ./aws/install")

	// clean up
	cmdRun("rm -r aws")
	cmdRun("rm -r awscliv2.zip")
}

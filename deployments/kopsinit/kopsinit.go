// Package kopsinit is a package of functions the assist with installing kops for kubernetes on the linux machine. Requires the following packages
//   commander (github.com/200106-uta-go/BAM-P2/pkg/commander) - used to send commands to the OS using "os/exec" standard go package
package kopsinit

import (
	"strings"

	"github.com/200106-uta-go/BAM-P2/pkg/commander"
)

// CheckInstall simply checks to see if kops is installed, returns true if installled, otherwise false.
func CheckInstall() bool {
	// logic to check if kops is already installed
	version := commander.CmdRunOutSilent("kops version")
	v := strings.Split(version, " ")
	if v[0] == "Version" {
		return true
	}
	return false
}

// KopsInstall downloads a version of kops (currently hard coded to v1.15.2) and makes it execatable then places
// it in "user/local/bin" (sudo permissions required)
func KopsInstall() {
	commander.CmdRun("curl -Lo kops https://github.com/kubernetes/kops/releases/download/v1.15.2/kops-linux-amd64")
	commander.CmdRun("chmod +x ./kops")
	commander.CmdRun("sudo mv ./kops /usr/local/bin/")
}

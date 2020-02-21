// Package kubectlinit is a package of functions the assist with installing kubectl for kubernetes on the linux machine. Requires the following packages
//   commander (github.com/200106-uta-go/BAM-P2/pkg/commander) - used to send commands to the OS using "os/exec" standard go package
package kubectlinit

import (
	"strings"

	"github.com/200106-uta-go/BAM-P2/pkg/commander"
)

// CheckInstall simply checks to see if kops is installed, returns true if installled, otherwise false.
func CheckInstall() bool {
	// logic to check if kops is already installed
	version := commander.CmdRunOutSilentNoErr("kubectl version --client --short")
	v := strings.Split(version, " ")
	if v[0] == "Client" {
		return true
	}
	return false
}

// KubectlInstall downloads a version of kubectl (currently hard coded to v1.17.3) and makes it execatable then places
// it in "user/local/bin" (sudo permissions required)
func KubectlInstall() {
	commander.CmdRun("curl -Lo kubectl https://storage.googleapis.com/kubernetes-release/release/v1.17.3/bin/linux/amd64/kubectl")
	commander.CmdRun("chmod +x ./kubectl")
	commander.CmdRun("sudo mv ./kubectl /usr/local/bin/")
}

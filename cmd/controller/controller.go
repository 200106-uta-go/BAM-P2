package main

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

// var subcommand string

/*
// Control struct outline
type Control struct {
	Apply -f -- filepath string
	Delete ${object} name  -- object string, name string
	Get ${object} name  -- object string, name string
	Describe   -- object string, name string
	Scale -- replicascount string, deployment filepath string
	Logs -- podname string
	Cluster-info -- no arguments
}
*/

// func init() {
// 	//gets the subcommand passed
// 	if len(os.Args) > 0 {
// 		subcommand = os.Args[1]
// 	}
// }

// //calls controller functions based on cli args
// func main() {
// 	switch subcommand {
// 	case "scale":
// 		scale(os.Args[2], "deploy.yaml")
// 	}
// }

// func cmdRun(s string) {
// 	mySlice := strings.Split(s, " ")
// 	cmd := exec.Command(mySlice[0], mySlice[1:]...)
// 	fmt.Println(cmd)
// 	err := cmd.Run()
// 	if err != nil {
// 		log.Println(err)
// 	}
// }

// //scales the deployed pods
// func scale(num string, deploy string) {
// 	fmt.Println("Scaling to", num)
// 	cmdRun(fmt.Sprintf("kubectl scale --replicas %s -f %s", num, deploy))
// }

// KubeApply creates string for kubectl apply command
func KubeApply(filepath string) string {
	var outputstring string
	if filepath == "" {
		outputstring = ""
	} else {
		outputstring = fmt.Sprintf("apply -f %s", filepath)
	}
	return outputstring
}

// KubeDelete creates string for kubectl delete command
func KubeDelete(object string, name string) string {
	var outputstring string
	if object == "" || name == "" {
		outputstring = ""
	} else {
		outputstring = fmt.Sprintf("delete %s %s", object, name)
	}
	return outputstring
}

// KubeGet creates string for kubectl get command
func KubeGet(object string, name string) string {
	var outputstring string
	if object == "" || name == "" {
		outputstring = ""
	} else {
		outputstring = fmt.Sprintf("get %s %s", object, name)
	}
	return outputstring
}

// KubeDescribe creates string for kubectl describe command
func KubeDescribe(object string, name string) string {
	var outputstring string
	if object == "" || name == "" {
		outputstring = ""
	} else {
		outputstring = fmt.Sprintf("describe %s %s", object, name)
	}
	return outputstring
}

// KubeScale creates string for kubectl scale command
func KubeScale(count string, deployment string) string {
	var outputstring string
	if count == "" || deployment == "" {
		outputstring = ""
	} else if strings.HasSuffix(deployment, ".yaml") {
		outputstring = fmt.Sprintf("scale --replicas=%s -f %s", count, deployment)
	} else {
		outputstring = fmt.Sprintf("scale --replicas=%s %s", count, deployment)
	}
	return outputstring
}

// KubeLogs creates string for kubectl logs command
func KubeLogs(podname string) string {
	var outputstring string
	if podname == "" {
		outputstring = ""
	} else {
		outputstring = fmt.Sprintf("logs %s", podname)
	}
	return outputstring
}

// KubeClusterInfo creates string for kubectl cluster-info command
func KubeClusterInfo() string {
	return "cluster-info"
}

// KubeCommand runs the kubectl CLI with a command
func KubeCommand(command string) {
	cmd := exec.Command(fmt.Sprintf("kubectl %s", command))
	err := cmd.Run()
	if err != nil {
		log.Println(err)
	}
}

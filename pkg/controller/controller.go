package controller

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

// KubeApply creates string for kubectl apply command
func KubeApply(filepath string) string {
	var outputstring string
	if filepath == "" {
		outputstring = ""
	} else {
		outputstring = fmt.Sprintf("apply -f %s", filepath)
	}
	return KubeCommand(outputstring)
}

// KubeDelete creates string for kubectl delete command
func KubeDelete(object string, name string) string {
	var outputstring string
	if object == "" || name == "" {
		outputstring = ""
	} else {
		outputstring = fmt.Sprintf("delete %s %s", object, name)
	}
	return KubeCommand(outputstring)
}

// KubeGet creates string for kubectl get command
func KubeGet(object string, name string) string {
	var outputstring string
	if object == "" {
		outputstring = ""
	} else {
		outputstring = fmt.Sprintf("get %s -o=json %s", object, name)
	}
	return KubeCommand(outputstring)
}

// KubeDescribe creates string for kubectl describe command
func KubeDescribe(object string, name string) string {
	var outputstring string
	if object == "" {
		outputstring = ""
	} else {
		outputstring = fmt.Sprintf("describe %s -o=json %s", object, name)
	}
	return KubeCommand(outputstring)
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
	return KubeCommand(outputstring)
}

// KubeLogs creates string for kubectl logs command
func KubeLogs(podname string) string {
	var outputstring string
	if podname == "" {
		outputstring = ""
	} else {
		outputstring = fmt.Sprintf("logs %s", podname)
	}
	return KubeCommand(outputstring)
}

// KubeRun starts a container in a pod using a dockerhub image
func KubeRun(image string) string {
	var outputstring string
	if image == "" {
		outputstring = ""
	} else {
		name := strings.Split(image, "/")
		if len(name) > 1 {
			outputstring = fmt.Sprintf("run %s --image %s", name[1], image)
		} else {
			outputstring = fmt.Sprintf("run %s --image %s", image, image)
		}
	}
	return KubeCommand(outputstring)
}

// KubeClusterInfo creates string for kubectl cluster-info command
func KubeClusterInfo() string {
	return KubeCommand("cluster-info")
}

// KubeCommand runs the kubectl CLI with a command
func KubeCommand(command string) string {

	//remove whitespace if emtry string was given as argument
	command = strings.TrimSpace(command)

	//create and run command
	cmd := exec.Command("kubectl", strings.Split(command, " ")...)
	//log.Println("Running", cmd)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	bytes, err := cmd.Output()
	if err != nil {
		log.Println(stderr.String())
	}
	return string(bytes)
}

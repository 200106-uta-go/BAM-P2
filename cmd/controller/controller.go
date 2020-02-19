package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

var subcommand string

func init() {
	//gets the subcommand passed
	if len(os.Args) > 0 {
		subcommand = os.Args[1]
	}
}

//calls controller functions based on cli args
func main() {
	switch subcommand {
	case "scale":
		scale(os.Args[2], "deploy.yaml")
	}
}

func cmdRun(s string) {
	mySlice := strings.Split(s, " ")
	cmd := exec.Command(mySlice[0], mySlice[1:]...)
	fmt.Println(cmd)
	err := cmd.Run()
	if err != nil {
		log.Println(err)
	}
}

//scales the deployed pods
func scale(num string, deploy string) {
	fmt.Println("Scaling to", num)
	cmdRun(fmt.Sprintf("kubectl scale --replicas %s -f %s", num, deploy))
}

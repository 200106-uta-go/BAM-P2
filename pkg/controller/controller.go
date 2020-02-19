package controller

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

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

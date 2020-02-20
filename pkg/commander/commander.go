package commander

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func cmd(s string) ([]byte, error) {
	cmdSlice := strings.Split(s, " ")
	command := exec.Command(cmdSlice[0], cmdSlice[1:]...)
	command.Stderr = os.Stderr
	out, err := command.Output()
	return out, err
}

func CmdRun(s string) {
	fmt.Println("Running command: " + s)
	out, err := cmd(s)
	if err != nil {
		log.Println("Command: " + s + ": Failed :: " + err.Error())
	}
	fmt.Println(string(out))
}

func CmdRunSilent(s string) {
	_, err := cmd(s)
	if err != nil {
		log.Println("Command: " + s + ": Failed :: " + err.Error())
	}
}

func CmdRunOut(s string) string {
	fmt.Println("Running command: " + s)
	out, err := cmd(s)
	if err != nil {
		log.Println("Command: " + s + ": Failed :: " + err.Error())
	}
	fmt.Println(string(out))
	return string(out)
}

func CmdRunOutSilent(s string) string {
	out, err := cmd(s)
	if err != nil {
		log.Println("Command: " + s + ": Failed :: " + err.Error())
	}
	return string(out)
}

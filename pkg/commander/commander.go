// Package commander is a wrapper package for the "os/exec" .Command function allowing diffrent types of returns and user controls
package commander

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

// cmd is a private function that is aonly available to this pkg but is used to send the command to the OS
func cmd(s string) ([]byte, error) {
	cmdSlice := strings.Split(s, " ")
	command := exec.Command(cmdSlice[0], cmdSlice[1:]...)
	//command.Stderr = os.Stderr
	out, err := command.Output()
	return out, err
}

// CmdRun runs a command provided by user (string) and provides feedback on what it is doing
func CmdRun(s string) {
	fmt.Println("Running command: " + s)
	out, err := cmd(s)
	if err != nil {
		log.Println("Command: " + s + ": Failed :: " + err.Error())
	}
	fmt.Println(string(out))
}

// CmdRunNoErr runs a command provided by user (string) and provides feedback on what it is doing but without showing errors
func CmdRunNoErr(s string) {
	fmt.Println("Running command: " + s)
	out, _ := cmd(s)
	fmt.Println(string(out))
}

// CmdRunSilent runs a command provided by user (string) and does not provide feedback on what it is doing
func CmdRunSilent(s string) {
	_, err := cmd(s)
	if err != nil {
		log.Println("Command: " + s + ": Failed :: " + err.Error())
	}
}

// CmdRunSilentNoErr runs a command provided by user (string) and does not provide feedback on what it is doing
// and without showing any errors
func CmdRunSilentNoErr(s string) {
	cmd(s)
}

// CmdRunOut runs a command provided by user (string) and returns the output in a string
func CmdRunOut(s string) string {
	fmt.Println("Running command: " + s)
	out, err := cmd(s)
	if err != nil {
		log.Println("Command: " + s + ": Failed :: " + err.Error())
	}
	fmt.Println(string(out))
	return string(out)
}

// CmdRunOutNoErr runs a command provided by user (string) and returns the output in a string but with out showing errors
func CmdRunOutNoErr(s string) string {
	fmt.Println("Running command: " + s)
	out, _ := cmd(s)
	fmt.Println(string(out))
	return string(out)
}

// CmdRunOutSilent runs a command provided by user (string) and returns the output in a string and does not provide feedback
// on what it is doing
func CmdRunOutSilent(s string) string {
	out, err := cmd(s)
	if err != nil {
		log.Println("Command: " + s + ": Failed :: " + err.Error())
	}
	return string(out)
}

// CmdRunOutSilentNoErr runs a command provided by user (string) and returns the output in a string and does not provide feedback
// on what it is doing but with out show
func CmdRunOutSilentNoErr(s string) string {
	out, _ := cmd(s)
	return string(out)
}

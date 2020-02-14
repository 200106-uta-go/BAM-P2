// GO-SECURE (LOGGER)
/*
	GO-SECURE (LOGGER) is part of a suite of security applications built in Go.
	This modual is a light-weight tool for handling and logging messages from
	other applications in this suite. (Academic concept)

	Revature: Brandon Locker (GameMasterTwig)
*/
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"

	"github.com/200106-uta-go/BAM-P2/pkg/filewriter"
)

// connection to logger server
type logger struct {
	LogAddress string `json:"logger"`
}

var connSig chan string
var logs []logger
var logAddr string

func init() {
	getLogAdder()
}

func getLogAdder() {
	if filewriter.CheckForFile("logConfig.json") {
		// file is present
		f, err := ioutil.ReadFile("logConfig.json")
		if err != nil {
			log.Fatalf("Unable to open logConfig: %+v", err)
		}

		// decode config (json)
		err = json.Unmarshal(f, &logs)
		if err != nil {
			log.Fatalf("Unable to decode logConfig: %+v", err)
		}
		logAddr = logs[0].LogAddress
	}
}

func main() {
	connSig = make(chan string)

	listn, _ := net.Listen("tcp", logAddr)

	fmt.Println("Logging Server listening on " + logAddr)

	for {
		go logSession(listn)
		<-connSig
	}
}

func logSession(listn net.Listener) {
	conn, _ := listn.Accept()
	defer conn.Close()

	// fmt.Println("New Connection On " + conn.LocalAddr().String())
	connSig <- "Done"

	for {
		buffer := make([]byte, 1024)

		// Attempt read
		_, err := conn.Read(buffer)
		if err != nil {
			break
		}

		cleanBuf := bytes.Trim(buffer, "\x00")

		// display log information here
		// just display to STDOUT for now

		if string(cleanBuf) != "" {
			fmt.Println(string(cleanBuf))

			/* if strings.Contains(string(cleanBuf), "Packet sent") || strings.Contains(string(cleanBuf), "Packet read") {
				// only print packets to consoul
				fmt.Println(string(cleanBuf))
			} else {
				// print and save to file
				fmt.Println(string(cleanBuf))
				if !filewriter.WriteAppend("log.txt", cleanBuf) {
					fmt.Println("Failed to write to log.txt")
				}
			} */
		}
	}
}

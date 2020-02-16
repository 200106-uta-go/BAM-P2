// GO-SECURE (RPROXY)
/*
	GO-SECURE (RPROXY) is part of a suite of security applications built in Go.
	This modual is a tool for forwarding requests to a server (reverse proxy).
	You can currently forward using HTTP or TCP protocals.
	(Academic concept)

	Revature: Brandon Locker (GameMasterTwig)
*/
package main

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"

	t "github.com/200106-uta-go/BAM-P2/pkg/tcprproxy"

	// used to send log messages to multiple writers
	_ "github.com/200106-uta-go/BAM-P2/pkg/logger"
)

type connection struct {
	FrontAddr string `json:"frontend"`
	BackAddr  string `json:"backend"`
	LogAddr   string `json:"logger"`
}

var connections []connection

func init() {
	// // code to pull from json
	// // open config file (json) at path
	// f, err := os.Open("config.json")
	// if err != nil {
	// 	log.Fatalf("Unable to open config: %+v", err)
	// }
	// // defer close
	// defer f.Close()

	// // decode config (json)
	// err = json.NewDecoder(f).Decode(&connections)
	// if err != nil {
	// 	log.Fatalf("Unable to decode config: %+v", err)
	// }

	//load in environment variables from .env
	//will print error message when running from docker image
	//because env file is passed into docker run command
	envErr := godotenv.Load("/home/ubuntu/go/src/github.com/200106-uta-go/BAM-P2/.env")
	if envErr != nil {
		if !strings.Contains(envErr.Error(), "no such file or directory") {
			log.Println("Error loading .env: ", envErr)
		}
	}

	connections = append(connections, connection{
		FrontAddr: ":" + os.Getenv("REV_FRONT"),
		BackAddr:  os.Getenv("REV_BACK"),
		LogAddr:   os.Getenv("LOG_ADDR"),
	})
}

func main() {
	connect := connections[0]
	t.TCPForward(connect.FrontAddr, connect.BackAddr, connect.LogAddr, "RPROXY")

	// var wg sync.WaitGroup
	// wg.Add(len(connections))
	// for _, c := range connections {
	// 	go func() {
	//		t.TCPForward(c.FrontAddr, c.BackAddr, c.LogAddr)
	// 		wg.Done()
	// 	}()
	// }
	// wg.Wait()
}

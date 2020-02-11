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
	"encoding/json"
	"log"
	"os"

	t "github.com/Gamemastertwig/go-secure/rproxy/tcprproxy"
)

type connection struct {
	FrontAddr string `json:"frontend"`
	BackAddr  string `json:"backend"`
	LogAddr   string `json:"logger"`
}

var connections []connection

func init() {
	// code to pull from json
	// open config file (json) at path
	f, err := os.Open("config.json")
	if err != nil {
		log.Fatalf("Unable to open config: %+v", err)
	}
	// defer close
	defer f.Close()

	// decode config (json)
	err = json.NewDecoder(f).Decode(&connections)
	if err != nil {
		log.Fatalf("Unable to decode config: %+v", err)
	}
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

package main

import (
	"flag"
	"fmt"

	"github.com/200106-uta-go/BAM-P2/deployments/cluster"
)

var create bool
var destroy bool
var up bool

func init() {
	// sets flag options
	flag.BoolVar(&create, "create", false, "Creates a cluster")          // create cluster
	flag.BoolVar(&destroy, "destroy", false, "Destroy/Delete a cluster") // destroy cluster
	flag.BoolVar(&up, "up", false, "Brings a cluster up")                // up cluster

	flag.Parse()
}

func main() {
	switch {
	case create:
		// create cluster
		cluster.PrepairCluster()
	case destroy:
		cluster.Destroy()
	case up:
		// bring cluster up
		cluster.Up()
	default:
		fmt.Println("Please add flag (-create | -destroy | -up)")
	}
}

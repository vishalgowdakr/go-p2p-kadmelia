package main

import (
	"fmt"
	"go-p2p/bootstrapserver"
	"go-p2p/cli"
	"os"
)

func main() {
	args := os.Getenv("P2P_CONFIG")
	if args == "bootstrapserver" {
		// Start the RPC server
		fmt.Print("Bootstrap server listening on port : 2233")
		bootstrapserver.StartBootstrapServer()
	} else if args == "test" {
		// Run the test
		TestBootstrapServer()
	} else {
		cli.Start()
	}
}

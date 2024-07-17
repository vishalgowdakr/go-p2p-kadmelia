package main

import (
	"fmt"
	"go-p2p/bootstrapserver"
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
		fmt.Println("Invalid argument")
	}
}

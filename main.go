package main

import (
	"fmt"
	"go-p2p/bootstrapserver"
	"go-p2p/client"
	"os"

	"github.com/libp2p/go-libp2p"
)

func initialise() {
	priv, err := client.LoadOrCreatePrivateKey()
	if err != nil {
		panic(err)
	}
	node, err := libp2p.New(
		libp2p.Identity(priv),
	)
	if err != nil {
		panic(err)
	}
	rt := client.NewRoutingTable(node.ID().String())
	client.InitializeGlobals(node.ID().String(), rt, node.Addrs())
	go client.StartRPCServer()
}

func main() {
	args := os.Getenv("P2P_CONFIG")
	if args == "bootstrapserver" {
		// Start the RPC server
		fmt.Print("Bootstrap server listening on port : 2233")
		bootstrapserver.StartBootstrapServer()
	} else if args == "node" {
		initialise()
	} else {
		initialise()
	}
}

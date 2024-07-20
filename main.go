package main

import (
	"fmt"
	"go-p2p/bootstrapserver"
	"go-p2p/cli"
	"go-p2p/client"
	"go-p2p/tree"
	"net"
	"os"
)

func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println("Error getting IP addresses:", err)
		return ""
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
			return ipnet.IP.String()
		}
	}
	return "unknown"
}

func initialise() {
	priv, err := client.LoadOrCreatePrivateKey()
	if err != nil {
		panic(err)
	}
	nodeId, err := client.CreateIDFromPrivateKey(priv)
	if err != nil {
		panic(err)
	}
	node := tree.NodeAddr{
		Id:            nodeId,
		ListenAddress: getLocalIP() + ":2233",
	}
	rt := client.NewRoutingTable(node.Id)
	nodes, err := client.RegisterNewNodeRPC(node)
	if err != nil {
		panic(err)
	}
	if len(nodes) == 1 {
		fmt.Println("This is the first node in the network")
	} else {
		for _, n := range nodes[1:] {
			client.GetRoutingTable(n)
		}
	}
	client.InitializeGlobals(rt, node)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go [bootstrapserver|node|cli]")
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "bootstrapserver":
		bootstrapserver.StartBootstrapServer()
	case "node":
		initialise()
		client.StartRPCServer()
	case "cli":
		initialise()
		go client.StartRPCServer()
		cli.Start()

	default:
		fmt.Printf("Unknown command: %s\n", command)
		fmt.Println("Usage: go run main.go [bootstrapserver|node|cli]")
		os.Exit(1)
	}
}

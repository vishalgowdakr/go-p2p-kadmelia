package main

import (
	"fmt"
	"go-p2p/client"
	"go-p2p/tree"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func generateRandomBinaryString(length int) string {
	rand.NewSource(time.Now().UnixNano())
	binary := ""
	for i := 0; i < length; i++ {
		binary += strconv.Itoa(rand.Intn(2))
	}
	return binary
}

func TestBootstrapServer() {
	// Create a file to save results
	file, err := os.OpenFile("test_bootstrap_server.txt", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)

	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()
	var registeredNodes []string

	// Generate and register 10 sample nodes
	for i := 0; i < 10; i++ {
		node := tree.NodeAddr{
			Id: generateRandomBinaryString(tree.IdLength),
		}

		nodes, err := client.RegisterNewNodeRPC(node)
		if err != nil {
			fmt.Fprintf(file, "Error registering node %d: %v\n", i, err)
		} else {
			fmt.Fprintf(file, "Registered node: %s\n", node.Id)
			registeredNodes = append(registeredNodes, node.Id)
		}
		fmt.Fprintf(file, "Response:\n")
		for _, n := range nodes {
			fmt.Fprintf(file, "\t%s\n", n.Id)
		}
	}

	// Get K nearest nodes for 5 random node IDs
	for _, n := range registeredNodes {
		fmt.Fprintf(file, "\nGetting K nearest nodes for ID: %s\n", n)
		nodes, err := client.GetKNearestNodesRPC(n)
		if err != nil {
			fmt.Fprintf(file, "Error getting K nearest nodes: %v\n", err)
		}
		fmt.Fprintf(file, "Response:\n")
		for _, n := range nodes {
			fmt.Fprintf(file, "\t%s\n", n.Id)
		}
	}

	testnode := "1111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111"
	fmt.Fprintf(file, "\nGetting K nearest nodes for ID: %s\n", testnode)
	nodes, err := client.GetKNearestNodesRPC(testnode)
	if err != nil {
		fmt.Fprintf(file, "Error getting K nearest nodes: %v\n", err)
	}
	fmt.Fprintf(file, "Response:\n")
	for _, n := range nodes {
		fmt.Fprintf(file, "\t%s\n", n.Id)
	}

	fmt.Println("Test completed. Results saved in test_bootstrap_server.txt")
}

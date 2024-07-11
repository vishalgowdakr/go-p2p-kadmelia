package main

import (
	"fmt"
	"go-p2p/bootstrapserver"
	"go-p2p/client"
	"go-p2p/tree"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func generateRandomBinaryString(length int) string {
	rand.Seed(time.Now().UnixNano())
	binary := ""
	for i := 0; i < length; i++ {
		binary += strconv.Itoa(rand.Intn(2))
	}
	return binary
}

func testWithSampleData() {
	// Start the RPC server
	go bootstrapserver.StartRpcServer()

	// Create a file to save results
	file, err := os.Create("test_results.txt")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()
	var registeredNodes []string

	// Generate and register 10 sample nodes
	for i := 0; i < 10; i++ {
		node := tree.NodeAddr{
			Id: generateRandomBinaryString(20),
			Ip: fmt.Sprintf("node%d", i),
		}

		nodes, err := client.RegisterNewNode(node)
		if err != nil {
			fmt.Fprintf(file, "Error registering node %d: %v\n", i, err)
		} else {
			fmt.Fprintf(file, "Registered node: %s - %s\n", node.Id, node.Ip)
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
		nodes, err := client.GetKNearestNodes(n)
		if err != nil {
			fmt.Fprintf(file, "Error getting K nearest nodes: %v\n", err)
		}
		fmt.Fprintf(file, "Response:\n")
		for _, n := range nodes {
			fmt.Fprintf(file, "\t%s\n", n.Id)
		}
		// Note: The actual results of GetKNearestNodes are not captured here.
		// You may need to modify the client.GetKNearestNodes function to return the results
		// so they can be written to the file.
	}
	testnode := "11111111111111111111"
	fmt.Fprintf(file, "\nGetting K nearest nodes for ID: %s\n", testnode)
	nodes, err := client.GetKNearestNodes(testnode)
	if err != nil {
		fmt.Fprintf(file, "Error getting K nearest nodes: %v\n", err)
	}
	fmt.Fprintf(file, "Response:\n")
	for _, n := range nodes {
		fmt.Fprintf(file, "\t%s\n", n.Id)
	}

	fmt.Println("Test completed. Results saved in test_results.txt")
}

func main() {
	testWithSampleData()
}

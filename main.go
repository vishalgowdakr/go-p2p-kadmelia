package main

import (
	"bufio"
	"fmt"
	"go-p2p/bootstrapserver"
	"go-p2p/client"
	"go-p2p/tree"
	"os"
	"strings"
)

func main() {
	go bootstrapserver.StartRpcServer()
	for {
		fmt.Println("Options:\n1. Register New Node\n2. Get K Nearest Nodes\n-1. Exit")
		reader := bufio.NewReader(os.Stdin)
		options, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			continue
		}
		options = strings.TrimSpace(options)

		switch options {
		case "-1":
			return
		case "1":
			fmt.Println("========== Register New Node ==========")
			fmt.Print("Enter Node Id: ")
			id, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("Error reading input:", err)
				continue
			}
			id = strings.TrimSpace(id)

			fmt.Print("Enter Node Ip: ")
			ip, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("Error reading input:", err)
				continue
			}
			ip = strings.TrimSpace(ip)

			node := tree.NodeAddr{
				Id: id,
				Ip: ip,
			}
			err = client.RegisterNewNode(node)
			if err != nil {
				fmt.Println("RPC Error:", err)
			}
		case "2":
			fmt.Println("========== Get K Nearest Nodes ==========")
			fmt.Print("Enter Node Id: ")
			id, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("Error reading input:", err)
				continue
			}
			id = strings.TrimSpace(id)

			err = client.GetKNearestNodes(id)
			if err != nil {
				fmt.Println("RPC Error:", err)
			}
		default:
			fmt.Println("Invalid choice")
		}
	}
}

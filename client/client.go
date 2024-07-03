package client

import (
	"fmt"
	t "go-p2p/tree"
	"log"
	"net/rpc"
)

func RegisterNewNode(addr t.NodeAddr) error {
	client, err := rpc.DialHTTP("tcp", "localhost"+":2233")
	if err != nil {
		log.Fatal("dialing:", err)
	}
	nodes := []t.NodeAddr{}
	err = client.Call("BootstrapServer.RegisterNewNode", &addr, &nodes)
	if err != nil {
		log.Fatal("register error:", err)
	}
	for _, node := range nodes {
		fmt.Println(node.Id)
	}
	return nil
}

func GetKNearestNodes(id string) error {
	client, err := rpc.DialHTTP("tcp", "localhost"+":2233")
	if err != nil {
		log.Fatal("dialing:", err)
	}
	nodes := []t.NodeAddr{}
	err = client.Call("BootstrapServer.GetKNearestNodes", &id, &nodes)
	if err != nil {
		log.Fatal("error:", err)
	}
	for _, node := range nodes {
		fmt.Println(node.Id)
	}
	return nil
}

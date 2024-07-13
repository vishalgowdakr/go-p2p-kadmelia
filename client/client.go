package client

import (
	"fmt"
	t "go-p2p/tree"
	"log"
	"net"
	"net/rpc"
)

var hostname string = "bootstrapserver"

func lookUp(hostname string) []net.IP {
	addr, err := net.LookupIP(hostname)
	if err != nil {
		fmt.Printf("Failed to lookup ip address for bootstrapserver")
		fmt.Printf(err.Error())
		return nil
	}
	for _, address := range addr {
		fmt.Printf(address.String())
	}
	return addr
}

func RegisterNewNode(addr t.NodeAddr) ([]t.NodeAddr, error) {
	address := lookUp(hostname)
	if address == nil {
		panic("Lookup failed")
	}
	client, err := rpc.DialHTTP("tcp", address[0].String()+":2233")
	if err != nil {
		log.Fatal("dialing:", err)
	}
	nodes := []t.NodeAddr{}
	err = client.Call("BootstrapServer.RegisterNewNode", &addr, &nodes)
	if err != nil {
		log.Fatal("register error:", err)
	}
	return nodes, nil
}

func GetKNearestNodes(id string) ([]t.NodeAddr, error) {
	address := lookUp(hostname)
	if address == nil {
		panic("Lookup failed")
	}
	client, err := rpc.DialHTTP("tcp", address[0].String()+":2233")
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
	return nodes, nil
}

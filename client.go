package main

import (
	"fmt"
	"go-p2p/bootstrapper"
	"log"
	"net/rpc"
)

type Args struct{}

func RpcClient() {
	var reply []bootstrapper.NodeAddr
	args := Args{}
	client, err := rpc.DialHTTP("tcp", "localhost"+":2233")
	if err != nil {
		log.Fatal("dialing:", err)
	}
	err = client.Call("BootStrapServer.RegisterNewNode", &args, &reply)
	fmt.Println("hello from  client")
	if err != nil {
		log.Fatal("arith error:", err)
	} else {
		fmt.Print("good")
	}
}

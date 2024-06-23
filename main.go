package main

import (
	"fmt"
	"go-p2p/bootstrapper"
	"log"
	"net"
	"net/http"
	"net/rpc"
)

func main() {
	bootstrapserver := new(bootstrapper.BootStrapServer)
	rpc.Register(bootstrapserver)
	rpc.HandleHTTP()
	fmt.Println("listening on port 2233")
	l, e := net.Listen("tcp", ":2233")
	if e != nil {
		log.Fatal("listen error:", e)
	}
	http.Serve(l, nil)
}

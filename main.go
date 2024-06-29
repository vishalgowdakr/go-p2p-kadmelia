package main

import (
	b "go-p2p/bootstrapserver"
)

func main() {
	go b.StartRpcServer()
}

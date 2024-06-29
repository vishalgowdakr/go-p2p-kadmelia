package bootstrapserver

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
)

type BootstrapServer struct{}

var tree = NewTree()

func (server BootstrapServer) RegisterNewNode(nodeAddr *NodeAddr, nodes *[]NodeAddr) error {
	tree.Insert(&Node{Addr: nodeAddr})
	*nodes = append(*nodes, tree.GetKNearestNodes(nodeAddr.Id)...)
	return nil
}

func (server BootstrapServer) GetKNearestNodes(id string, nodes *[]NodeAddr) error {
	*nodes = append(*nodes, tree.GetKNearestNodes(id)...)
	return nil
}

func StartRpcServer() {
	server := new(BootstrapServer)
	rpc.Register(server)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":2233")
	if e != nil {
		log.Fatal("listen error:", e)
	}
	http.Serve(l, nil)
}

package bootstrapserver

import (
	t "go-p2p/tree"
	"log"
	"net"
	"net/http"
	"net/rpc"
)

type BootstrapServer struct{}

var tree = t.NewTree()

func (server BootstrapServer) RegisterNewNode(nodeAddr *t.NodeAddr, nodes *[]t.NodeAddr) error {
	tree.Insert(&t.Node{Addr: nodeAddr})
	*nodes = append(*nodes, tree.GetKNearestNodes(nodeAddr.Id)...)
	return nil
}

func (server BootstrapServer) GetKNearestNodes(id *string, nodes *[]t.NodeAddr) error {
	*nodes = append(*nodes, tree.GetKNearestNodes(*id)...)
	return nil
}

func StartBootstrapServer() {
	server := new(BootstrapServer)
	rpc.Register(server)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":2233")
	if e != nil {
		log.Fatal("listen error:", e)
	}
	http.Serve(l, nil)
}

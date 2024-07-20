package bootstrapserver

import (
	t "go-p2p/tree"
	"log"
	"net"
	"net/http"
	"net/rpc"
)

// BootstrapServer represents the RPC server for bootstrap functionality
type BootstrapServer struct{}

// Initialize a new tree for storing node addresses
var tree = t.NewTree()

// RegisterNewNode registers a new node with the bootstrap server and returns the k-nearest nodes
func (server BootstrapServer) RegisterNewNode(nodeAddr *t.NodeAddr, nodes *[]t.NodeAddr) error {
	// Insert the new node into the tree
	tree.Insert(&t.Node{Addr: nodeAddr})
	// Get the k-nearest nodes and append them to the nodes slice
	*nodes = append(*nodes, tree.GetKNearestNodes(nodeAddr.Id)...)
	log.Println("New node registered:", nodeAddr.Id)
	return nil
}

// GetKNearestNodes returns the k-nearest nodes for a given node ID
func (server BootstrapServer) GetKNearestNodes(id *string, nodes *[]t.NodeAddr) error {
	// Get the k-nearest nodes and append them to the nodes slice
	log.Println("Getting k-nearest nodes for:", id)
	*nodes = append(*nodes, tree.GetKNearestNodes(*id)...)
	return nil
}

// StartBootstrapServer starts the bootstrap server, listening for RPC calls
func StartBootstrapServer() {
	// Create a new BootstrapServer instance
	server := new(BootstrapServer)
	// Register the BootstrapServer instance with the RPC package
	err := rpc.Register(server)
	if err != nil {
		log.Fatalf("Error registering RPC server: %v", err)
	}
	// Register an HTTP handler for RPC messages to be received via HTTP
	rpc.HandleHTTP()
	// Listen on TCP port 2233
	listener, err := net.Listen("tcp", ":2233")
	if err != nil {
		log.Fatalf("Listen error: %v", err)
	}
	// Serve HTTP requests on the listener
	log.Println("Bootstrap server started on port 2233")
	err = http.Serve(listener, nil)
	if err != nil {
		log.Fatalf("Error serving: %v", err)
	}
}

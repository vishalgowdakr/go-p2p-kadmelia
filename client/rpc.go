package client

import (
	"fmt"
	"go-p2p/tree"
	"log"
	"net"
	"net/http"
	"net/rpc"
)

type Client struct {
}

// New returns a new instance of the Kademlia client
func New() *Client {
	return &Client{}
}

// StartRPCServer starts the RPC server for the Kademlia client
func StartRPCServer() {
	server := New()
	err := rpc.Register(server)
	if err != nil {
		log.Fatal("Error registering RPC server:", err)
	}
	rpc.HandleHTTP()
	listener, err := net.Listen("tcp", ":2233")
	if err != nil {
		log.Fatal("Listen error:", err)
	}
	log.Println("RPC server listening on port 2233")
	err = http.Serve(listener, nil)
	if err != nil {
		log.Fatal("HTTP serve error:", err)
	}
}

// GetChunk retrieves a chunk of data by its ID
func (client *Client) GetChunk(chunkID *string, data *[]byte) error {
	store, err := NewChunkStore("file_chunks.db")
	if err != nil {
		fmt.Println("Error creating store:", err)
		return err
	}
	defer store.Close()

	retrievedChunk, err := store.Retrieve(*chunkID)
	if err != nil {
		fmt.Println("Error retrieving chunk:", err)
		return err
	}

	*data = retrievedChunk.Data
	return nil
}

// GetRoutingTable retrieves the current routing table and inserts a new node address
func (client *Client) GetRoutingTable(a *tree.NodeAddr, rt *RoutingTable) error {
	err := myRoutingTable.InsertIntoRoutingTable(*a)
	if err != nil {
		return fmt.Errorf("error inserting into routing table: %w", err)
	}
	myrt := GetMyRoutingTable()
	*rt = myrt
	return nil
}

type FindNodeArgs struct {
	peer  *tree.NodeAddr
	peers *[]tree.NodeAddr
}

// FindNode finds a node by its ID and updates the routing table
func (client *Client) FindNode(nodeID *string, a *FindNodeArgs) error {
	if len(*a.peers) == 0 {
		return fmt.Errorf("no peers provided")
	}

	peer := (*a.peers)[0]
	err := myRoutingTable.InsertIntoRoutingTable(peer)
	if err != nil {
		return fmt.Errorf("error inserting into routing table: %w", err)
	}

	err = Findnode(*nodeID, a.peer, a.peers)
	if err != nil {
		return fmt.Errorf("error finding node: %w", err)
	}

	return nil
}

// Store stores a chunk of data and updates the routing table
func (client *Client) Store(chunk *FileChunk, a *tree.NodeAddr) error {
	err := myRoutingTable.InsertIntoRoutingTable(*a)
	if err != nil {
		return fmt.Errorf("error inserting into routing table: %w", err)
	}

	err = Store(chunk)
	if err != nil {
		return fmt.Errorf("error storing chunk: %w", err)
	}

	return nil
}

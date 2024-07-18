package client

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"time"

	peerstore "github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
)

type Client struct {
	pingTimeout time.Duration
}

// New returns a new instance of the Kademlia client
func New() *Client {
	c := &Client{
		pingTimeout: 3 * time.Second,
	}

	return c
}

// TODO:
type getargs struct {
	peer *peerstore.AddrInfo
	data *[]byte
}

func StartRPCServer() {
	server := new(Client)
	rpc.Register(server)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":2233")
	if e != nil {
		log.Fatal("listen error:", e)
	}
	http.Serve(l, nil)
}

func GetChunk(chunkID *string, data *[]byte) error {
	store, err := NewChunkStore("file_chunks.db")
	if err != nil {
		fmt.Println("Error creating store:", err)
		return nil
	}
	defer store.Close()
	retrievedChunk, err := store.Retrieve(*chunkID)
	if err != nil {
		fmt.Println("Error retrieving chunk:", err)
		return err
	}
	*data = retrievedChunk.data
	return nil
}

func (client *Client) GetRoutingTable(a *[]multiaddr.Multiaddr, rt *RoutingTable) error {
	addrs := *a
	peer := addrs[0]
	myRoutingTable.InsertIntoRoutingTable(peer)
	_, myrt := GetMyNodeID()
	*rt = myrt
	return nil
}

type args struct {
	peer  *peerstore.AddrInfo
	peers *[]peerstore.AddrInfo
	addrs *[]multiaddr.Multiaddr
}

func (client *Client) FindNode(nodeID *string, a *args) error {
	addrs := *a.addrs
	peer := addrs[0]
	myRoutingTable.InsertIntoRoutingTable(peer)
	err := Findnode(*nodeID, a.peer, a.peers)
	return err
}

func (client *Client) Store(chunk *FileChunk, a *[]multiaddr.Multiaddr) error {
	addrs := *a
	peer := addrs[0]
	myRoutingTable.InsertIntoRoutingTable(peer)
	err := Store(chunk)
	if err != nil {
		return err
	}
	return nil
}

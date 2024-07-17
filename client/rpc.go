package client

import (
	"fmt"
	"time"

	peerstore "github.com/libp2p/go-libp2p/core/peer"
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

type args struct {
	peer  *peerstore.AddrInfo
	peers *[]peerstore.AddrInfo
}

func (client *Client) FindNode(nodeID *string, a *args) error {
	err := Findnode(*nodeID, a.peer, a.peers)
	return err
}

func (client *Client) Store(chunk *FileChunk, reply *string) error {
	err := Store(chunk, reply)
	if err != nil {
		return err
	}
	return nil
}

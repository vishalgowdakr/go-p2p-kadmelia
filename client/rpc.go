package client

import (
	"context"
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
func (client *Client) FindNode(ctx *context.Context, nodeID *string) (peerstore.AddrInfo, []peerstore.AddrInfo, error)

func (client *Client) Store(chunk []byte) (peerstore.AddrInfo, error)

func (client *Client) FindChunk(chunkID *string) ([]byte, error)

func (client *Client) Ping(ctx *context.Context, peer *peerstore.AddrInfo) error

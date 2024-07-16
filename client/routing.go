package client

import (
	"context"
	"log"

	peerstore "github.com/libp2p/go-libp2p/core/peer"
)

func Findnode(ctx *context.Context, nodeID *string) (peerstore.AddrInfo, []peerstore.AddrInfo, error) {
	clientId := context.Background().Value("id").(*string)
	bucketIndex, err := xorStrings(*clientId, *nodeID)
	if err != nil {
		log.Fatal(err)
	}
	rt := context.Background().Value("rt").(RoutingTable)
	list := rt[bucketIndex].List
	curr := list.Back()

	// Assume bucket_array is a predefined array of peerstore.AddrInfo
	var bucket_array []peerstore.AddrInfo

	// Iterate through the list to find the nodeID
	for curr != nil {
		node := curr.Value.(peerstore.AddrInfo) // Adjust this line based on the actual type stored in the list
		if node.ID.String() == *nodeID {
			return node, bucket_array, nil
		}
		bucket_array = append(bucket_array, curr.Value.(peerstore.AddrInfo))
		curr = curr.Prev()
	}

	// If nodeID is not found
	return peerstore.AddrInfo{}, bucket_array, nil
}


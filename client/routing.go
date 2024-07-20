package client

import (
	"go-p2p/tree"
	"log"
)

func Findnode(nodeID string, peer *tree.NodeAddr, peers *[]tree.NodeAddr) error {
	client := GetMyAddrInfo()
	bucketIndex, err := xorStrings(client.Id, nodeID)
	if err != nil {
		log.Fatal(err)
	}
	rt := GetMyRoutingTable()
	list := rt[bucketIndex].Queue
	i := 0

	// Assume bucket_array is a predefined array of peerstore.AddrInfo
	//
	var bucket_array []tree.NodeAddr

	// Iterate through the list to find the nodeID
	for i < len(list) {
		curr := list[i]
		if list[i].Id == nodeID {
			*peer = curr
			return nil
		}
		bucket_array = append(bucket_array, curr)
		i++
	}

	// If nodeID is not found
	*peers = bucket_array
	return nil
}

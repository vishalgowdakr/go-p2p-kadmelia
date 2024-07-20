package client

import (
	"fmt"
	t "go-p2p/tree"
	"log"
	"net"
	"net/rpc"
	"os"
)

var hostname string = "bootstrapserver"

type DataArr []struct {
	index int
	data  *[]byte
}

func GetRoutingTable(peer t.NodeAddr) {
	client, err := rpc.DialHTTP("tcp", peer.ListenAddress)
	if err != nil {
		log.Fatal("dialing:", err)
	}
	me := GetMyAddrInfo()
	rt := GetMyRoutingTable()
	err = client.Call("Client.GetRoutingTable", &me, &rt)
	if err != nil {
		log.Fatal("error:", err)
	}
	UpdateRoutingTable(rt)
}

func DownloadFile(torrentFilePath string, downloadChunk func(string, *t.NodeAddr) ([]byte, error)) error {
	torrentFile, err := deserializeTorrentFile(torrentFilePath)
	if err != nil {
		return err
	}

	file, err := os.Create("downloads/" + torrentFile.Filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Pre-allocate slice with the correct size
	chunks := make([][]byte, len(torrentFile.FileChunkIndex))

	for chunkID, index := range torrentFile.FileChunkIndex {
		peer := torrentFile.FilechunksId[chunkID]
		data, err := downloadChunk(chunkID, &peer)
		if err != nil {
			return fmt.Errorf("failed to download chunk %s: %w", chunkID, err)
		}
		chunks[index] = data
	}

	// Write chunks in order
	for _, chunk := range chunks {
		_, err := file.Write(chunk)
		if err != nil {
			return fmt.Errorf("failed to write chunk: %w", err)
		}
	}

	// Verify file size
	info, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file info: %w", err)
	}
	if info.Size() != int64(torrentFile.Filesize) {
		return fmt.Errorf("file size mismatch: expected %d, got %d", torrentFile.Filesize, info.Size())
	}

	return nil
}

func GetRPC(cid string, peer *t.NodeAddr) ([]byte, error) {
	client, err := rpc.DialHTTP("tcp", peer.ListenAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to dial: %w", err)
	}
	defer client.Close()

	var data []byte
	err = client.Call("Client.GetChunk", &cid, &data)
	if err != nil {
		return nil, fmt.Errorf("RPC call failed: %w", err)
	}

	return data, nil
}

func StoreRPC(chunk FileChunk, peer *t.NodeAddr) error {
	client, err := rpc.DialHTTP("tcp", peer.ListenAddress)
	if err != nil {
		log.Fatal("dialing:", err)
	}
	err = client.Call("Client.Store", &chunk, peer)
	if err != nil {
		log.Fatal("error:", err)
	}
	return err
}

func FindNodeRPC(nodeID string, peer *t.NodeAddr, peers *[]t.NodeAddr) error {
	err := Findnode(nodeID, peer, peers)
	if err != nil {
		return fmt.Errorf("Findnode error: %v", err)
	}
	for peer == nil && len(*peers) > 0 {
		alternatePeer := (*peers)[0]
		*peers = (*peers)[1:] // Remove the used peer from the list
		la := alternatePeer.ListenAddress

		client, err := rpc.DialHTTP("tcp", la)
		if err != nil {
			log.Printf("dialing error: %v", err)
			continue // Try the next peer
		}

		nodeID = alternatePeer.Id // Assuming Id is a string

		args := FindNodeArgs{peer: peer, peers: peers}
		err = client.Call("Client.FindNode", &nodeID, &args)
		if err != nil {
			log.Printf("RPC call error: %v", err)
			continue // Try the next peer
		}
	}
	if peer == nil {
		return fmt.Errorf("failed to find node %s", nodeID)
	}
	return nil
}
func RegisterNewNodeRPC(addr t.NodeAddr) ([]t.NodeAddr, error) {
	addrs, err := net.LookupHost(hostname)
	if err != nil {
		log.Fatal("lookup error:", err)
	}
	client, err := rpc.DialHTTP("tcp", addrs[0]+":2233")
	if err != nil {
		log.Fatal("dialing:", err)
	}
	nodes := []t.NodeAddr{}
	err = client.Call("BootstrapServer.RegisterNewNode", &addr, &nodes)
	if err != nil {
		log.Fatal("register error:", err)
	}
	return nodes, nil
}

func GetKNearestNodesRPC(id string) ([]t.NodeAddr, error) {
	addrs, err := net.LookupHost(hostname)
	if err != nil {
		log.Fatal("lookup error:", err)
	}
	client, err := rpc.DialHTTP("tcp", addrs[0]+":2233")
	if err != nil {
		log.Fatal("dialing:", err)
	}
	nodes := []t.NodeAddr{}
	err = client.Call("BootstrapServer.GetKNearestNodes", &id, &nodes)
	if err != nil {
		log.Fatal("error:", err)
	}
	return nodes, nil
}

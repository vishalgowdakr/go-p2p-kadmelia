package client

import (
	"fmt"
	t "go-p2p/tree"
	"log"
	"net"
	"net/rpc"
	"os"

	peerstore "github.com/libp2p/go-libp2p/core/peer"
)

var hostname string = "bootstrapserver"

func GetRoutingTable(peer *peerstore.AddrInfo) {
	ip, port := getIpAndPort(*peer)
	client, err := rpc.DialHTTP("tcp", ip+":"+port)
	if err != nil {
		log.Fatal("dialing:", err)
	}
	_, rt := GetMyNodeID()
	err = client.Call("Client.GetRoutingTable", &myAddrInfo, &rt)
	if err != nil {
		log.Fatal("error:", err)
	}
	UpdateRoutingTable(rt)
}

func DownloadFile(torrentFilePath string, downloadChunk func(string, *peerstore.AddrInfo) ([]byte, error)) error {
	torrentFile, err := deserializeTorrentFile(torrentFilePath)
	if err != nil {
		return err
	}
	file, err := os.Create("downloads/" + torrentFile.filename)
	if err != nil {
		return err
	}
	defer file.Close()

	for chunkID := range torrentFile.filechunksId {
		peer := torrentFile.filechunksId[chunkID]
		data, err := downloadChunk(chunkID, &peer)
		if err != nil {
			return err
		}
		file.Write(data)
	}
	return nil
}

func GetRPC(cid string, peer *peerstore.AddrInfo) ([]byte, error) {
	ip, port := getIpAndPort(*peer)
	client, err := rpc.DialHTTP("tcp", ip+":"+port)
	if err != nil {
		log.Fatal("dialing:", err)
	}
	data := []byte{}
	err = client.Call("Client.GetChunk", &cid, &data)
	if err != nil {
		log.Fatal("error:", err)
	}
	return data, err
}

func StoreRPC(chunk *FileChunk, reply *string, peer *peerstore.AddrInfo) error {
	ip, port := getIpAndPort(*peer)
	client, err := rpc.DialHTTP("tcp", ip+":"+port)
	if err != nil {
		log.Fatal("dialing:", err)
	}
	err = client.Call("Client.Store", chunk, reply)
	if err != nil {
		log.Fatal("error:", err)
	}
	return err
}

func FindNodeRPC(nodeID *string, peer *peerstore.AddrInfo, peers *[]peerstore.AddrInfo) error {
	err := Findnode(*nodeID, peer, peers)
	if err != nil {
		log.Fatal(err)
	}
	for peer == nil {
		alternatePeer := (*peers)[0]
		ip, port := getIpAndPort(alternatePeer)
		client, err := rpc.DialHTTP("tcp", ip+":"+port)
		if err != nil {
			log.Fatal("dialing:", err)
		}
		nodeID = (*string)(&alternatePeer.ID)
		type args struct {
			peer  *peerstore.AddrInfo
			peers *[]peerstore.AddrInfo
		}
		reply := args{peer: peer, peers: peers}
		err = client.Call("Client.FindNode", nodeID, &reply)
	}
	return err
}

func RegisterNewNodeRPC(addr t.NodeAddr) ([]t.NodeAddr, error) {
	client, err := rpc.DialHTTP("tcp", "localhost"+":2233")
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
	client, err := rpc.DialHTTP("tcp", "localhost"+":2233")
	if err != nil {
		log.Fatal("dialing:", err)
	}
	nodes := []t.NodeAddr{}
	err = client.Call("BootstrapServer.GetKNearestNodes", &id, &nodes)
	if err != nil {
		log.Fatal("error:", err)
	}
	for _, node := range nodes {
		fmt.Println(node.Id)
	}
	return nodes, nil
}

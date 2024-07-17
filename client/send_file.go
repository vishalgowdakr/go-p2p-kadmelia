package client

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/rpc"
	"os"
	"strings"

	"github.com/ipfs/go-cid"
	peerstore "github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multihash"
)

type TorrentFile struct {
	filename     string
	filesize     int
	filechunksId map[string]peerstore.AddrInfo
}

type FileChunk struct {
	id    string
	index int
	data  []byte
	owner peerstore.AddrInfo
}

const chunkSize = 1 * 1024 * 1024 // 1 MB chunks

func sendChunk(filechunk FileChunk) (peerstore.AddrInfo, error) {
	peers, err := GetKNearestNodesRPC(filechunk.id)
	if err != nil {
		return peerstore.AddrInfo{}, err
	}

	for _, peer := range peers {
		ip, port := getIpAndPort(*peer.Host)
		if ip == "" || port == "" {
			continue
		}
		client, err := rpc.DialHTTP("tcp", ip+port)
		if err != nil {
			log.Fatal("dialing:", err)
		}
		reply := ""
		error := client.Call("Client.Store", &filechunk, &reply)
		if error != nil && reply != "OK" {
			log.Fatal("error:", error)
		}
		return *peer.Host, nil
	}
	return peerstore.AddrInfo{}, fmt.Errorf("No peers available to store chunk")
}

func SendFile(filepath string, sendChunk func(FileChunk) (peerstore.AddrInfo, error)) error {
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}

	filename := strings.Split(filepath, "/")[len(strings.Split(filepath, "/"))-1]
	totalSize := fileInfo.Size()
	torrentFile := TorrentFile{
		filename: filename,
		filesize: int(totalSize),
	}
	fileChunksId := make([]string, 0)

	/* 	sentBytes := int64(0) */

	reader := bufio.NewReader(file)
	buffer := make([]byte, chunkSize)
	index := 0

	for {
		n, err := reader.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		chunk := buffer[:n]
		cid, err := getContentId(chunk)
		if err != nil {
			return err
		}

		fileChunk := FileChunk{
			id:    cid,
			index: index,
			data:  chunk,
		}
		index++
		fileChunksId = append(fileChunksId, cid)
		owner, err := sendChunk(fileChunk)
		if err != nil {
			return err
		}
		torrentFile.filechunksId[cid] = owner

		//sentBytes += int64(n)
		//progress := float64(sentBytes) / float64(totalSize) * 100
		// Log or update progress here
		// log.Printf("Progress: %.2f%%", progress)
	}
	//serialize torrent file and save it
	serializeTorrentFile(torrentFile)

	return nil
}

func serializeTorrentFile(torrentFile TorrentFile) error {
	// use marshalling to serialize the torrent file
	serializedData, err := json.Marshal(torrentFile)
	if err != nil {
		return fmt.Errorf("error serializing torrent file: %w", err)
	}
	//save the serialized data to a file
	err = os.WriteFile(torrentFile.filename+".torrent", serializedData, 0644)
	if err != nil {
		return fmt.Errorf("error saving torrent file: %w", err)
	}
	return nil
}

func deserializeTorrentFile(filepath string) (TorrentFile, error) {
	// read the serialized data from the file
	serializedData, err := os.ReadFile(filepath)
	if err != nil {
		return TorrentFile{}, fmt.Errorf("error reading torrent file: %w", err)
	}
	// use unmarshalling to deserialize the torrent file
	var torrentFile TorrentFile
	err = json.Unmarshal(serializedData, &torrentFile)
	if err != nil {
		return TorrentFile{}, fmt.Errorf("error deserializing torrent file: %w", err)
	}
	return torrentFile, nil
}

func getContentId(data []byte) (string, error) {
	mh, err := multihash.Sum(data, multihash.SHA2_256, -1)
	if err != nil {
		return "", err
	}
	c := cid.NewCidV1(cid.Raw, mh)
	return c.String(), nil
}

func getIpAndPort(addr peerstore.AddrInfo) (string, string) {
	// Use the first multiaddr in the list
	if len(addr.Addrs) == 0 {
		return "", ""
	}

	multiaddr := addr.Addrs[0]

	// Extract the IP address and port
	var ip, port string

	// Split the multiaddr into its components
	parts := strings.Split(multiaddr.String(), "/")

	for i := 1; i < len(parts); i += 2 {
		switch parts[i] {
		case "ip4", "ip6":
			ip = parts[i+1]
		case "tcp", "udp":
			port = parts[i+1]
		}
	}

	return ip, port
}

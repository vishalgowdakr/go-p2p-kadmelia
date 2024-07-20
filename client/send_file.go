package client

import (
	"bufio"
	"encoding/json"
	"fmt"
	"go-p2p/tree"
	"io"
	"log"
	"net/rpc"
	"os"
	"strings"
	"sync"

	"github.com/ipfs/go-cid"
	"github.com/multiformats/go-multihash"
)

type TorrentFile struct {
	Filename       string
	Filesize       int
	FileChunkIndex map[string]int
	FilechunksId   map[string]tree.NodeAddr
}

type FileChunk struct {
	Id    string
	Index int
	Data  []byte
}

const chunkSize = 1 * 1024 * 1024 // 1 MB chunks

func SendChunk(filechunk FileChunk) (tree.NodeAddr, error) {
	peers, err := GetKNearestNodesRPC(filechunk.Id)
	if err != nil {
		return tree.NodeAddr{}, err
	}

	for _, peer := range peers {

		client, err := rpc.DialHTTP("tcp", peer.ListenAddress)
		if err != nil {
			log.Fatal("dialing:", err)
		}
		error := client.Call("Client.Store", &filechunk, &tree.NodeAddr{})
		if error != nil {
			log.Fatal("error:", error)
			continue
		}
		return peer, nil
	}
	return tree.NodeAddr{}, fmt.Errorf("No peers available to store chunk")
}

func SendFile(filepath string, sendChunk func(FileChunk) (tree.NodeAddr, error)) error {
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
		Filename:       filename,
		Filesize:       int(totalSize),
		FileChunkIndex: make(map[string]int),
		FilechunksId:   make(map[string]tree.NodeAddr),
	}

	var wg sync.WaitGroup
	errChan := make(chan error, int(totalSize/chunkSize)+1)
	resultChan := make(chan struct {
		cid   string
		index int
		owner tree.NodeAddr
	}, int(totalSize/chunkSize)+1)

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

		chunk := make([]byte, n) // Create a new slice to avoid data race
		copy(chunk, buffer[:n])

		wg.Add(1)
		go func(index int, chunk []byte) {
			defer wg.Done()

			cid, err := getContentId(chunk)
			if err != nil {
				errChan <- err
				return
			}

			fileChunk := FileChunk{
				Id:    cid,
				Index: index,
				Data:  chunk,
			}

			owner, err := sendChunk(fileChunk)
			if err != nil {
				errChan <- err
				return
			}

			resultChan <- struct {
				cid   string
				index int
				owner tree.NodeAddr
			}{cid, index, owner}
		}(index, chunk)

		index++
	}

	go func() {
		wg.Wait()
		close(resultChan)
		close(errChan)
	}()

	for result := range resultChan {
		torrentFile.FilechunksId[result.cid] = result.owner
		torrentFile.FileChunkIndex[result.cid] = result.index
	}

	if len(errChan) > 0 {
		return <-errChan // Return the first error encountered
	}

	return serializeTorrentFile(torrentFile)
}

func serializeTorrentFile(torrentFile TorrentFile) error {
	// use marshalling to serialize the torrent file
	serializedData, err := json.Marshal(torrentFile)
	if err != nil {
		return fmt.Errorf("error serializing torrent file: %w", err)
	}
	//save the serialized data to a file
	err = os.WriteFile("torrent/"+torrentFile.Filename+".torrent", serializedData, 0644)
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

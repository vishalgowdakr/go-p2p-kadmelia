package file

import (
	"bufio"
	"go-p2p/client"
	"go-p2p/tree"
	"io"
	"os"

	"github.com/ipfs/go-cid"
	peerstore "github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multihash"
)

type TorrentFile struct {
	filename   string
	filesize   int
	filechunks []FileChunk
}

type FileChunk struct {
	id    string
	index int
	data  []byte
	owner peerstore.AddrInfo
}

const chunkSize = 1 * 1024 * 1024 // 1 MB chunks

func sendChunk(filechunk FileChunk) error {
	// Call Store rpc after implementing
	peers, err := client.GetKNearestNodes(filechunk.id)
	if err != nil {
		return err
	}
	for _, peer := range peers {
	}
	return nil
}

func SendFile(filepath string, sendChunk func(FileChunk) error) error {
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}

	totalSize := fileInfo.Size()
	sentBytes := int64(0)

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
		err = sendChunk(fileChunk)
		if err != nil {
			return err
		}

		//sentBytes += int64(n)
		//progress := float64(sentBytes) / float64(totalSize) * 100
		// Log or update progress here
		// log.Printf("Progress: %.2f%%", progress)
	}

	return nil
}

func getContentId(data []byte) (string, error) {
	mh, err := multihash.Sum(data, multihash.SHA2_256, -1)
	if err != nil {
		return "", err
	}
	c := cid.NewCidV1(cid.Raw, mh)
	return c.String(), nil
}

func chunkFile(filepath string) ([][]byte, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var chunks [][]byte
	buf := make([]byte, chunkSize)

	for {
		n, err := file.Read(buf)
		if err != nil && err != io.EOF {
			return nil, err
		}
		if n == 0 {
			break
		}
		chunk := make([]byte, n)
		copy(chunk, buf[:n])
		chunks = append(chunks, chunk)
	}

	return chunks, nil
}

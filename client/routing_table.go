package client

import (
	"container/list"
	"fmt"
	"go-p2p/tree"
	"strconv"

	peerstore "github.com/libp2p/go-libp2p/core/peer"
	multiaddr "github.com/multiformats/go-multiaddr"
)

type Bucket struct {
	ID   string
	List list.List
}

type RoutingTable []Bucket

// XOR performs the XOR operation on two byte slices of equal length.
func xor(a, b []byte) ([]byte, error) {
	if len(a) != len(b) {
		return nil, fmt.Errorf("inputs must have the same length")
	}
	result := make([]byte, len(a))
	for i := range a {
		result[i] = a[i] ^ b[i]
	}
	return result, nil
}

// XORStrings performs XOR on two strings and returns the number of leading zeros in the binary representation of the result until the first 1 is encountered.
func xorStrings(s1, s2 string) (int, error) {
	b1 := []byte(s1)
	b2 := []byte(s2)

	// Pad the shorter string with zero bytes to match the length of the longer string.
	maxLen := len(b1)
	if len(b2) > maxLen {
		maxLen = len(b2)
	}
	paddedB1 := make([]byte, maxLen)
	copy(paddedB1, b1)
	paddedB2 := make([]byte, maxLen)
	copy(paddedB2, b2)

	result, err := xor(paddedB1, paddedB2)
	if err != nil {
		return 0, err
	}

	// Convert the result to a binary string and count leading zeros.
	binaryStr := ""
	for _, b := range result {
		binaryStr += fmt.Sprintf("%08b", b) // Convert each byte to an 8-bit binary string.
	}

	leadingZeros := 0
	for _, ch := range binaryStr {
		if ch == '0' {
			leadingZeros++
		} else {
			break
		}
	}

	return leadingZeros, nil
}

// NewRoutingTable creates a new routing table based on the given ID.
func NewRoutingTable(ID string) RoutingTable {
	tempID := ""
	rt := RoutingTable{}
	for _, char := range ID {
		bit := int(char) - int('0')
		bitStr := strconv.Itoa(bit)
		tempID += bitStr
		bucket := Bucket{
			ID: tempID,
		}
		rt = append(rt, bucket)
	}
	return rt
}

// ConstructRoutingTable constructs a routing table by merging the given routing tables.
func ConstructRoutingTable(rt, prt RoutingTable) RoutingTable {
	for i, bucket := range prt {
		if bucket.ID == rt[i].ID {
			rt[i].List = bucket.List
		}
	}
	return rt
}

// Ping checks if a peer is reachable.
func Ping(peer *peerstore.AddrInfo) bool {
	// TODO: implement this
	return false
}

// function to check if the bucket already contains the peer
func (bkt *Bucket) contains(peer *peerstore.AddrInfo) bool {
	for e := bkt.List.Front(); e != nil; e = e.Next() {
		addrInfo, ok := e.Value.(*peerstore.AddrInfo)
		if !ok {
			return false
		}
		if addrInfo.ID == peer.ID {
			return true
		}
	}
	return false
}

// InsertIntoBucket inserts a peer into the bucket.
func (bkt *Bucket) insertIntoBucket(peer *peerstore.AddrInfo) {
	frontElement := bkt.List.Front()
	if frontElement != nil {
		if bkt.contains(peer) {
			return
		}
		addrInfo, ok := frontElement.Value.(*peerstore.AddrInfo)
		if !ok {
			return
		}
		if Ping(addrInfo) {
			bkt.List.PushBack(peer)
			bkt.List.Remove(frontElement)
		} else if bkt.List.Len() < tree.K {
			bkt.List.PushBack(peer)
		}
	}
}

// InsertIntoRoutingTable inserts a peer into the routing table.
func (rt *RoutingTable) InsertIntoRoutingTable(ma multiaddr.Multiaddr) error {
	peer, err := getPeerInfo(ma)
	if err != nil {
		return err
	}

	myID, err := getNodeID()
	if err != nil {
		return err
	}

	peerID := peer.ID
	bucketIndex, err := xorStrings(myID, peerID.String())
	if err != nil {
		return err
	}

	(*rt)[bucketIndex].insertIntoBucket(peer)
	return nil
}

// getPeerInfo extracts peer information from the multiaddress.
func getPeerInfo(ma multiaddr.Multiaddr) (*peerstore.AddrInfo, error) {
	addr, err := multiaddr.NewMultiaddr(ma.String())
	if err != nil {
		return nil, err
	}
	peer, err := peerstore.AddrInfoFromP2pAddr(addr)
	if err != nil {
		return nil, err
	}
	return peer, nil
}

// getNodeID reads and decodes the node ID from a file.
func getNodeID() (string, error) {
	encodedID, err := ReadFromFile("../client/node_id.txt")
	if err != nil {
		return "", err
	}
	return DecodeNodeID(encodedID)
}

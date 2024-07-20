package client

import (
	"fmt"
	"go-p2p/tree"
	"strconv"
)

type Bucket struct {
	ID     string
	MaxLen int
	Queue  []tree.NodeAddr
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
	binaryStr := ""
	for _, b := range result {
		binaryStr += fmt.Sprintf("%08b", b)
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

// NewRoutingTable creates a new routing table based on the given ID(hex format).
func NewRoutingTable(ID string) RoutingTable {
	binaryID := tree.BinaryString(ID)
	str1 := ""
	str2 := ""
	rt := RoutingTable{}
	for _, char := range binaryID {
		bit := int(char) - int('0')
		bit = 1 - bit
		bitStr := strconv.Itoa(bit)
		str2 = str1 + bitStr
		bucket := Bucket{
			ID:     str2,
			MaxLen: tree.K,
			Queue:  []tree.NodeAddr{},
		}
		str1 += string(char)
		rt = append(rt, bucket)
	}
	return rt
}

// ConstructRoutingTable constructs a routing table by merging the given routing tables.
func ConstructRoutingTable(rt, prt RoutingTable) RoutingTable {
	for i, bucket := range prt {
		if i < len(rt) && (bucket.ID == rt[i].ID || rt[i].ID == bucket.ID[:len(rt[i].ID)] || bucket.ID == rt[i].ID[:len(bucket.ID)]) {
			rt[i].Queue = bucket.Queue
		}
	}
	return rt
}

// Ping checks if a peer is reachable.
func Ping(peer *tree.NodeAddr) bool {
	// TODO: implement this
	return false
}

// contains checks if the bucket already contains the peer
func (bkt *Bucket) contains(peer *tree.NodeAddr) bool {
	for _, addrInfo := range bkt.Queue {
		if addrInfo.Id == peer.Id {
			return true
		}
	}
	return false
}

// InsertIntoBucket inserts a peer into the bucket.
func (bkt *Bucket) insertIntoBucket(peer *tree.NodeAddr) {
	if len(bkt.Queue) < bkt.MaxLen {
		if !bkt.contains(peer) {
			bkt.Queue = append(bkt.Queue, *peer)
		}
	} else {
		if Ping(&bkt.Queue[0]) {
			return
		} else {
			bkt.Queue = append(bkt.Queue[1:], *peer)
		}
	}
}

// InsertIntoRoutingTable inserts a peer into the routing table.
func (rt *RoutingTable) InsertIntoRoutingTable(ma tree.NodeAddr) error {
	me := GetMyAddrInfo()
	myID := me.Id
	peerID := ma.Id
	bucketIndex, err := xorStrings(myID, peerID)
	if err != nil {
		return err
	}
	if bucketIndex < len(*rt) {
		(*rt)[bucketIndex].insertIntoBucket(&ma)
	}
	return nil
}

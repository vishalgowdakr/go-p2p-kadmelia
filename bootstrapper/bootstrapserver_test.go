package bootstrapper_test

import (
	"crypto/rand"
	"crypto/sha1"
	"fmt"
	b "go-p2p/bootstrapper"
	"math/big"
	"testing"
)

func generateSHA1Hash(data []byte) string {
	hasher := sha1.New()
	hasher.Write(data)
	hash := hasher.Sum(nil)
	return fmt.Sprintf("%x", hash)
}

func generateRandomIP() (string, error) {
	// Local IP addresses range from 192.168.0.0 to 192.168.255.255
	max := big.NewInt(256)

	// Generate random values for the third and fourth octets
	bigInt1, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}
	bigInt2, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("192.168.%d.%d", bigInt1.Int64(), bigInt2.Int64()), nil
}

func TestRegisterNewNode(t *testing.T) {
	btsvr := b.BootStrapServer{}
	node1 := b.NodeAddr{Id: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 14, 15, 16, 17, 18, 19, 20}, Ip: "192.168.100.1"}
	var reply1 []b.NodeAddr
	btsvr.RegisterNewNode(&node1, &reply1)
	if len(reply1) != 1 {
		t.Fatalf("Node 1 Register failed")
	}
	node2 := b.NodeAddr{Id: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 21, 12, 13, 14, 14, 15, 16, 17, 18, 19, 20}, Ip: "192.168.100.2"}
	var reply2 []b.NodeAddr
	btsvr.RegisterNewNode(&node2, &reply2)
	if len(reply2) != 2 {
		t.Fatalf("Node 2 Register failed")
	}
	node3 := b.NodeAddr{Id: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 31, 12, 13, 14, 14, 15, 16, 17, 18, 19, 20}, Ip: "192.168.100.3"}
	var reply3 []b.NodeAddr
	btsvr.RegisterNewNode(&node3, &reply3)
	if len(reply3) != 3 {
		t.Fatalf("Node 3 Register failed")
	}
	node4 := b.NodeAddr{Id: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 41, 12, 13, 14, 14, 15, 16, 17, 18, 19, 20}, Ip: "192.168.100.4"}
	var reply4 []b.NodeAddr
	btsvr.RegisterNewNode(&node4, &reply4)
	if len(reply4) != 4 {
		t.Fatalf("Node 4 Register failed")
	}
	node5 := b.NodeAddr{Id: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 51, 12, 13, 14, 14, 15, 16, 17, 18, 19, 20}, Ip: "192.168.100.5"}
	var reply5 []b.NodeAddr
	btsvr.RegisterNewNode(&node5, &reply5)
	if len(reply5) != 5 {
		t.Fatalf("Node 5 Register failed")
	}
	node6 := b.NodeAddr{Id: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 61, 12, 13, 14, 14, 15, 16, 17, 18, 19, 20}, Ip: "192.168.100.6"}
	var reply6 []b.NodeAddr
	btsvr.RegisterNewNode(&node6, &reply6)
	if len(reply6) != 6 {
		t.Fatalf("Node 6 Register failed")
	}
	node7 := b.NodeAddr{Id: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 71, 12, 13, 14, 14, 15, 16, 17, 18, 19, 20}, Ip: "192.168.100.7"}
	var reply7 []b.NodeAddr
	btsvr.RegisterNewNode(&node7, &reply7)
	if len(reply7) != 7 {
		t.Fatalf("Node 7 Register failed")
	}
}

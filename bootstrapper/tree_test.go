package bootstrapper

import (
	"bufio"
	"bytes"
	"os"
	"testing"
)

func TestSerializationDeserialization(t *testing.T) {
	// Create a sample binary tree
	root := NewNode(&NodeAddr{id: []byte{1, 0, 1, 0}, ip: "192.168.1.1"})
	root.left = NewNode(&NodeAddr{id: []byte{1, 1, 0, 0}, ip: "192.168.1.2"})
	root.right = NewNode(&NodeAddr{id: []byte{0, 1, 0, 1}, ip: "192.168.1.3"})
	root.left.left = NewNode(&NodeAddr{id: []byte{0, 0, 0, 1}, ip: "192.168.1.4"})
	root.left.right = NewNode(&NodeAddr{id: []byte{1, 0, 1, 1}, ip: "192.168.1.5"})

	// Serialize the tree to a buffer
	var buffer bytes.Buffer
	writer := bufio.NewWriter(&buffer)
	Serialize(root, writer)
	writer.Flush()

	// Deserialize the tree from the buffer
	scanner := bufio.NewScanner(&buffer)
	scanner.Split(bufio.ScanWords)
	deserializedRoot := deSerialize(scanner)

	// Verify the tree structure
	if !compareTrees(root, deserializedRoot) {
		t.Error("Deserialized tree does not match original tree")
	}
}

func TestGetKNearestNodes(t *testing.T) {
	// Construct a sample binary tree
	root := NewNode(&NodeAddr{id: []byte{1, 0, 1, 0}, ip: "192.168.1.1"})
	root.left = NewNode(&NodeAddr{id: []byte{1, 1, 0, 0}, ip: "192.168.1.2"})
	root.right = NewNode(&NodeAddr{id: []byte{0, 1, 0, 1}, ip: "192.168.1.3"})
	root.left.left = NewNode(&NodeAddr{id: []byte{0, 0, 0, 1}, ip: "192.168.1.4"})
	root.left.right = NewNode(&NodeAddr{id: []byte{1, 0, 1, 1}, ip: "192.168.1.5"})

	btree := Btree{root: root}

	// Find the 3 nearest nodes to a given ID
	id := []byte{1, 0, 1, 1}
	k := 3
	nearestNodes := GetKNearestNodes(btree, id, k)

	// Verify the number of nearest nodes found
	if len(nearestNodes) != k {
		t.Errorf("Expected %d nearest nodes, but got %d", k, len(nearestNodes))
	}

	// Verify the nearest nodes are correct
	expectedIPs := []string{"192.168.1.5", "192.168.1.1", "192.168.1.2"}
	for i, node := range nearestNodes {
		if node.ip != expectedIPs[i] {
			t.Errorf("Expected IP %s, but got %s", expectedIPs[i], node.ip)
		}
	}
}

func compareTrees(node1, node2 *Node) bool {
	if node1 == nil && node2 == nil {
		return true
	}
	if node1 == nil || node2 == nil {
		return false
	}
	if !bytes.Equal(node1.key.id, node2.key.id) || node1.key.ip != node2.key.ip {
		return false
	}
	return compareTrees(node1.left, node2.left) && compareTrees(node1.right, node2.right)
}

func TestMain(m *testing.M) {
	// Set up test environment
	code := m.Run()

	// Clean up test environment
	os.Remove("tree.txt")

	// Exit
	os.Exit(code)
}

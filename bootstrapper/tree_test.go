package bootstrapper_test

import (
	"bufio"
	"fmt"
	b "go-p2p/bootstrapper"
	"os"
	"testing"
)

func TestBtreeSerialization(t *testing.T) {
	// Create a sample binary tree
	root := b.NewNode(&b.NodeAddr{Id: []byte{1, 2, 3}, Ip: "192.168.1.1"})
	root.Left = b.NewNode(&b.NodeAddr{Id: []byte{4, 5}, Ip: "192.168.1.2"})
	root.Right = b.NewNode(&b.NodeAddr{Id: []byte{6, 7}, Ip: "192.168.1.3"})
	root.Left.Left = b.NewNode(&b.NodeAddr{Id: []byte{8}, Ip: "192.168.1.4"})

	// Serialize the binary tree to a file
	file, err := os.Create("tree_test.txt")
	if err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}
	writer := bufio.NewWriter(file)
	b.Serialize(root, writer)
	writer.Flush()
	file.Close()

	// DeSerialize the binary tree from the file
	file, err = os.Open("tree_test.txt")
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)
	deSerializedRoot := b.DeSerialize(scanner)
	file.Close()

	// Perform inorder traversal of the deSerialized tree and compare
	expectedInorder := "ID: [8], IP: 192.168.1.4\nID: [4 5], IP: 192.168.1.2\nID: [1 2 3], IP: 192.168.1.1\nID: [6 7], IP: 192.168.1.3\n"
	actualInorder := captureInorder(deSerializedRoot)
	if expectedInorder != actualInorder {
		t.Fatalf("Expected inorder traversal: %s, but got: %s", expectedInorder, actualInorder)
	}
}

func captureInorder(root *b.Node) string {
	if root == nil {
		return ""
	}
	result := ""
	if root.Left != nil {
		result += captureInorder(root.Left)
	}
	result += "ID: " + fmt.Sprint(root.Key.Id) + ", IP: " + root.Key.Ip + "\n"
	if root.Right != nil {
		result += captureInorder(root.Right)
	}
	return result
}

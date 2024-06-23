package bootstrapper

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

const MARKER = -1

type Btree struct {
	Root *Node
}

type NodeAddr struct {
	Id []byte
	Ip string
}

type Node struct {
	Key    *NodeAddr
	Left   *Node
	Right  *Node
	Parent *Node
}

func NewNode(key *NodeAddr) *Node {
	return &Node{Key: key, Left: nil, Right: nil}
}

// Function to Serialize the binary tree to a file
func Serialize(root *Node, writer *bufio.Writer) {
	if root == nil {
		writer.WriteString(fmt.Sprintf("%d ", MARKER))
		return
	}
	// Serialize node address: id length, id bytes, and ip string
	writer.WriteString(fmt.Sprintf("%d %s ", len(root.Key.Id), root.Key.Ip))
	for _, b := range root.Key.Id {
		writer.WriteString(fmt.Sprintf("%d ", b))
	}
	Serialize(root.Left, writer)
	Serialize(root.Right, writer)
}

// Function to deSerialize the binary tree from a file
func DeSerialize(scanner *bufio.Scanner) *Node {
	if !scanner.Scan() {
		return nil
	}
	val, _ := strconv.Atoi(scanner.Text())
	if val == MARKER {
		return nil
	}
	// DeSerialize node address
	idLen := val
	ip := ""
	if scanner.Scan() {
		ip = scanner.Text()
	}
	id := make([]byte, idLen)
	for i := 0; i < idLen; i++ {
		if scanner.Scan() {
			b, _ := strconv.Atoi(scanner.Text())
			id[i] = byte(b)
		}
	}
	root := NewNode(&NodeAddr{Id: id, Ip: ip})
	root.Left = DeSerialize(scanner)
	root.Right = DeSerialize(scanner)
	return root
}

// Inorder traversal of the binary tree
func Inorder(root *Node) {
	if root != nil {
		Inorder(root.Left)
		fmt.Printf("ID: %v, IP: %s\n", root.Key.Id, root.Key.Ip)
		Inorder(root.Right)
	}
}

// Function to get the binary tree from a file or return an empty tree
func GetBTree() Btree {
	var btree Btree
	file, err := os.Open("tree.txt")
	if err != nil {
		btree.Root = &Node{}
		return btree
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)
	btree.Root = DeSerialize(scanner)
	return btree
}

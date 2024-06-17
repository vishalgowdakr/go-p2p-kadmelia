package bootstrapper

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
)

const MARKER = -1

type Btree struct {
	root *Node
}

type NodeAddr struct {
	id []byte
	ip string
}

type Node struct {
	key   *NodeAddr
	left  *Node
	right *Node
}

func NewNode(key *NodeAddr) *Node {
	return &Node{key: key, left: nil, right: nil}
}

// Function to Serialize the binary tree to a file
func Serialize(root *Node, writer *bufio.Writer) {
	if root == nil {
		writer.WriteString(fmt.Sprintf("%d ", MARKER))
		return
	}
	// Serialize node address: id length, id bytes, and ip string
	writer.WriteString(fmt.Sprintf("%d %s ", len(root.key.id), root.key.ip))
	for _, b := range root.key.id {
		writer.WriteString(fmt.Sprintf("%d ", b))
	}
	Serialize(root.left, writer)
	Serialize(root.right, writer)
}

// Function to deSerialize the binary tree from a file
func deSerialize(scanner *bufio.Scanner) *Node {
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
	root := NewNode(&NodeAddr{id: id, ip: ip})
	root.left = deSerialize(scanner)
	root.right = deSerialize(scanner)
	return root
}

// Inorder traversal of the binary tree
func inorder(root *Node) {
	if root != nil {
		inorder(root.left)
		fmt.Printf("ID: %v, IP: %s\n", root.key.id, root.key.ip)
		inorder(root.right)
	}
}

// Function to get the binary tree from a file or return an empty tree
func GetBTree() Btree {
	var btree Btree
	file, err := os.Open("tree.txt")
	if err != nil {
		btree.root = nil
		return btree
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)
	btree.root = deSerialize(scanner)
	return btree
}

// Function to find the K nearest nodes based on the given ID
func GetKNearestNodes(btree Btree, id []byte, k int) []*NodeAddr {
	var nearestNodes []*NodeAddr

	// Helper function to compute Hamming distance
	hammingDistance := func(id1, id2 []byte) int {
		if len(id1) != len(id2) {
			return -1
		}
		distance := 0
		for i := 0; i < len(id1); i++ {
			if id1[i] != id2[i] {
				distance++
			}
		}
		return distance
	}

	// Recursive function to find the K nearest nodes
	var findNearest func(*Node)
	findNearest = func(node *Node) {
		if node == nil {
			return
		}
		// Calculate the distance from the current node to the target ID
		distance := hammingDistance(node.key.id, id)
		// Insert the node into the list of nearest nodes if within K limit
		if len(nearestNodes) < k {
			nearestNodes = append(nearestNodes, node.key)
		} else {
			// Check if the current node is closer than the farthest in the nearest nodes
			farthestDistance := hammingDistance(nearestNodes[len(nearestNodes)-1].id, id)
			if distance < farthestDistance {
				nearestNodes[len(nearestNodes)-1] = node.key
			}
		}
		// Sort the nodes by distance and keep only the K closest ones
		sort.Slice(nearestNodes, func(i, j int) bool {
			return hammingDistance(nearestNodes[i].id, id) < hammingDistance(nearestNodes[j].id, id)
		})
		// Recur for left and right subtrees
		findNearest(node.left)
		findNearest(node.right)
	}

	findNearest(btree.root)
	return nearestNodes
}

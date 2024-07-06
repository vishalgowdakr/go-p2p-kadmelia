package tree

import (
	"fmt"
	"strings"
)

const K = 20

type NodeAddr struct {
	Id string
	Ip string
}

type Node struct {
	Addr   *NodeAddr
	Parent *Node
	Left   *Node
	Right  *Node
}

type Tree struct {
	Head *Node
}

func NewTree() Tree {
	return Tree{
		Head: NewNode(NodeAddr{}),
	}
}

func NewNode(addr NodeAddr) *Node {
	return &Node{Addr: &addr}
}

func StrToIntArr(str string) []int {
	intArr := make([]int, len(str))
	for i, char := range str {
		intArr[i] = int(char) - int('0')
	}
	return intArr
}

func (tree Tree) SaveToDisk() error {
	return nil
}

func (tree *Tree) Insert(node *Node) bool {
	// Check if tree or node is nil
	if tree == nil || tree.Head == nil || node == nil {
		return false
	}

	id := StrToIntArr(node.Addr.Id)
	if len(id) == 0 {
		return false
	}

	curr := tree.Head
	for _, bit := range id {
		if curr.Addr != nil {
			temp := curr.Addr
			curr.Addr = nil
			tempNode := NewNode(*temp)
			tree.Insert(tempNode)
		}
		if bit == 0 {
			if curr.Left == nil {
				curr.Left = node
				node.Parent = curr
				return true
			} else {
				curr = curr.Left
			}
		} else {
			if curr.Right == nil {
				curr.Right = node
				node.Parent = curr
				return true
			} else {
				curr = curr.Right
			}
		}
	}
	return true
}

func (tree Tree) Print() {
	if tree.Head == nil {
		fmt.Println("Empty tree")
		return
	}
	printNode(tree.Head, 0, 0)
}

func printNode(node *Node, level int, dir int) {
	if node == nil {
		return
	}
	indent := strings.Repeat("  ", level)
	if node.Addr != nil {
		fmt.Printf("%s- ID: %s, IP: %s\n", indent, node.Addr.Id, node.Addr.Ip)
	} else {
		fmt.Printf("%s- %d\n", indent, dir)
	}
	printNode(node.Left, level+1, 0)
	printNode(node.Right, level+1, 1)
}

func TreeSearch(n *Node, visited_map map[*Node]bool, nodes *[]*Node) {
	if n == nil || visited_map == nil || nodes == nil || visited_map[n] == true || len(*nodes) == K {
		return
	}
	visited_map[n] = true
	if n.Addr != nil {
		*nodes = append(*nodes, n)
	}
	TreeSearch(n.Left, visited_map, nodes)
	TreeSearch(n.Right, visited_map, nodes)
	TreeSearch(n.Parent, visited_map, nodes)
}

func (tree Tree) FindNode(id string) *Node {
	if tree.Head == nil {
		return nil
	}
	curr := tree.Head
	bits := StrToIntArr(id)
	for _, bit := range bits {
		if curr == nil {
			return nil
		}
		if curr.Addr != nil && curr.Addr.Id == id {
			return curr
		}
		if bit == 0 {
			curr = curr.Left
		} else {
			curr = curr.Right
		}
	}
	return nil
}

func (tree Tree) GetKNearestNodes(id string) []NodeAddr {
	nodes := []*Node{}
	n := tree.FindNode(id)
	if n == nil {
		return []NodeAddr{}
	}
	TreeSearch(n, map[*Node]bool{}, &nodes)
	knodes := []NodeAddr{}
	for _, node := range nodes {
		if node != nil && node.Addr != nil {
			knodes = append(knodes, *node.Addr)
		}
	}
	return knodes
}

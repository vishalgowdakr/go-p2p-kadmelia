package bootstrapper

import (
	"log"
)

// global constants
const K = 20

// user defined types
type ID []byte

type BootStrapServer struct{}

// global variable
var btree Btree = GetBTree()

// helper functions
func ByteToBits(b byte) []int {
	bits := make([]int, 8)
	for i := 0; i < 8; i++ {
		// Extract the i-th bit from the right
		bits[7-i] = int((b >> i) & 1)
	}
	return bits
}

// takes dir direction 1 for right and 0 for left
func (n *Node) getChild(dir int) *Node {
	if n == nil {
		return nil
	}
	if dir == 1 {
		if n.Right == nil {
			n.Parent = n
			n.Right = &Node{}
		}
		return n.Right
	} else {
		if n.Left == nil {
			n.Parent = n
			n.Left = &Node{}
		}
		return n.Left
	}
}

func ModBFS(n *Node, visited_map map[*Node]bool, nodes []*Node) {
	if n == nil || visited_map[n] == true || len(nodes) == K {
		return
	}
	visited_map[n] = true
	if n.Key != nil {
		nodes = append(nodes, n)
	}
	ModBFS(n.Left, visited_map, nodes)
	ModBFS(n.Right, visited_map, nodes)
	ModBFS(n.Parent, visited_map, nodes)
}

// RPC functions
func (b *BootStrapServer) RegisterNewNode(args *NodeAddr, reply *[]NodeAddr) error {
	var bits []int
	for _, value := range args.Id {
		bits = append(bits, ByteToBits(value)...)
	}
	curr := btree.Root
	for _, bit := range bits {
		curr = curr.getChild(bit)
		if curr == nil {
			break
		} else if curr.Left == nil && curr.Right == nil && curr.Key != nil {
			var key NodeAddr = NodeAddr{
				Id: curr.Key.Id,
				Ip: curr.Key.Ip,
			}
			b.RegisterNewNode(&NodeAddr{Id: key.Id, Ip: key.Ip}, reply)
		} else if curr.Right == nil && curr.Left == nil && curr.Key == nil {
			log.Print("key added")
			curr.Key = args
			return nil
		}
	}
	err := b.GetKNearestNodes(args, reply)
	if err != nil {
		return err
	}
	log.Print(len(*reply))
	return nil
}

func (b *BootStrapServer) GetKNearestNodes(args *NodeAddr, reply *[]NodeAddr) error {
	var bits []int
	for _, value := range args.Id {
		bits = append(bits, ByteToBits(value)...)
	}
	nodes := []*Node{}
	curr := btree.Root
	for _, bit := range bits {
		curr = curr.getChild(bit)
		if curr == nil {
			log.Print("curr == nil")
			break
		}
		if curr.Left == nil && curr.Right == nil && curr.Key != nil {
			ModBFS(curr, map[*Node]bool{}, nodes)
		}
	}
	var kNodes []NodeAddr = []NodeAddr{}
	log.Print("hello")
	for _, node := range nodes {
		kNodes = append(kNodes, *node.Key)
	}
	*reply = kNodes
	return nil
}

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

func SetBTreeRoot(root *Node) {
	btree.Root = root
}

func ByteToBits(b byte) []int {
	bits := make([]int, 8)
	for i := 0; i < 8; i++ {
		// Extract the i-th bit from the right
		bits[7-i] = int((b >> i) & 1)
	}
	return bits
}

// takes dir direction 1 for right and 0 for left
func (n *Node) GetChild(dir int) *Node {
	if n == nil {
		return nil
	}
	if dir == 1 {
		if n.Right == nil {
			n.Right = &Node{Parent: n}
		}
		return n.Right
	} else {
		if n.Left == nil {
			n.Left = &Node{Parent: n}
		}
		return n.Left
	}
}

func ModBFS(n *Node, visited_map map[*Node]bool, nodes *[]*Node) {
	if n == nil || visited_map[n] || len(*nodes) == K {
		return
	}
	visited_map[n] = true

	// Check if this node has a valid NodeAddr
	if n.Key != nil && n.Key.Id != nil && len(n.Key.Id) > 0 {
		*nodes = append(*nodes, n)
		log.Printf("ModBFS added node: %v", n.Key)
		if len(*nodes) == K {
			return
		}
	}

	// Continue searching in all directions
	if n.Left != nil {
		ModBFS(n.Left, visited_map, nodes)
	}
	if n.Right != nil {
		ModBFS(n.Right, visited_map, nodes)
	}
	if n.Parent != nil {
		ModBFS(n.Parent, visited_map, nodes)
	}
}

// RPC functions
func (b *BootStrapServer) RegisterNewNode(args *NodeAddr, reply *string) error {
	var bits []int
	for _, value := range args.Id {
		bits = append(bits, ByteToBits(value)...)
	}
	curr := btree.Root
	for _, bit := range bits {
		curr = curr.GetChild(bit)
		if curr == nil {
			break
		} else if curr.Left == nil && curr.Right == nil && curr.Key != nil {
			var key NodeAddr = NodeAddr{
				Id: curr.Key.Id,
				Ip: curr.Key.Ip,
			}
			b.RegisterNewNode(&NodeAddr{Id: key.Id, Ip: key.Ip}, reply)
			*reply = "Success"
			return nil // Avoid infinite recursion
		} else if curr.Right == nil && curr.Left == nil && curr.Key == nil {
			log.Print("key added")
			curr.Key = args
			*reply = "Success"
			return nil
		}
	}
	return nil
}

func (b *BootStrapServer) GetKNearestNodes(args *NodeAddr, reply *[]NodeAddr) error {
	var bits []int
	for _, value := range args.Id {
		bits = append(bits, ByteToBits(value)...)
	}
	nodes := []*Node{}
	visited := make(map[*Node]bool)
	curr := btree.Root

	log.Printf("Searching for node with ID: %v", args.Id)

	// Traverse the tree to find the closest matching node
	for _, bit := range bits {
		if curr == nil {
			log.Print("Reached nil node during traversal")
			break
		}
		next := curr.GetChild(bit)
		if next == nil {
			log.Print("Next node is nil, stopping traversal")
			break
		}
		curr = next
	}

	log.Printf("Starting ModBFS from node: %v", curr.Key)

	// Use ModBFS to find K nearest nodes
	ModBFS(curr, visited, &nodes)

	log.Printf("ModBFS found %d nodes", len(nodes))

	// Convert Node pointers to NodeAddr and append to reply
	for _, node := range nodes {
		if node != nil && node.Key != nil {
			*reply = append(*reply, *node.Key)
			log.Printf("Added node to reply: %v", node.Key)
		}
		if len(*reply) >= K {
			log.Print("Reached K nodes, stopping")
			break
		}
	}

	log.Printf("Returning %d nodes", len(*reply))

	return nil
}

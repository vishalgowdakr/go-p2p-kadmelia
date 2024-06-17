package bootstrapper

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
)

type ID []byte

type BootStrapServer struct{}

var btree Btree = GetBTree()

func byteToBits(b byte) []int {
	bits := make([]int, 8)
	for i := 0; i < 8; i++ {
		// Extract the i-th bit from the right
		bits[7-i] = int((b >> i) & 1)
	}
	return bits
}

// takes dir direction 0 for right and 1 for left
func (n *Node) getChild(dir int) *Node {
	if dir == 0 {
		if n.right == nil {
			n.right = &Node{}
		}
		return n.right
	} else {
		if n.left == nil {
			n.left = &Node{}
		}
		return n.left
	}
}

func (b *BootStrapServer) RegisterNewNode(args *NodeAddr) {
	var bits []int
	for _, value := range args.id {
		bits = append(bits, byteToBits(value)...)
	}
	curr := btree.root
	for _, bit := range bits {
		curr = curr.getChild(bit)
		if curr.left == nil && curr.right == nil && curr.key.id != nil {
			var key NodeAddr = NodeAddr{
				id: curr.key.id,
				ip: curr.key.ip,
			}
			b.RegisterNewNode(&NodeAddr{id: key.id, ip: key.ip})
		} else if curr.right == nil && curr.left == nil && curr.key.id == nil {
			curr.key = args
			return
		}
	}
}

func (b *BootStrapServer) GetKNearestNodes(args *NodeAddr) {

}

func main() {
	bootstrapserver := new(BootStrapServer)
	rpc.Register(bootstrapserver)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":2233")
	if e != nil {
		log.Fatal("listen error:", e)
	}
	http.Serve(l, nil)
}

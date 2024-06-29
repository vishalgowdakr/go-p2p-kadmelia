package bootstrapper_test

import (
	"bytes"
	// "fmt"
	b "go-p2p/bootstrapper"
	// "reflect"
	"testing"
)

func TestModBFS(t *testing.T) {
	// Create sample nodes
	node1 := &b.Node{Key: &b.NodeAddr{}}
	node2 := &b.Node{Key: &b.NodeAddr{}, Parent: node1}
	node3 := &b.Node{Key: &b.NodeAddr{}, Parent: node1}
	node1.Left = node2
	node1.Right = node3
	node4 := &b.Node{Key: &b.NodeAddr{Id: []byte{4}, Ip: "4.4.4.4"}, Parent: node2}
	node2.Left = node4
	node5 := &b.Node{Key: &b.NodeAddr{Id: []byte{5}, Ip: "5.5.5.5"}, Parent: node3}
	node3.Right = node5

	testCases := []struct {
		name        string
		startNode   *b.Node
		expectedIDs [][]byte
	}{
		{
			name:        "Start from root",
			startNode:   node1,
			expectedIDs: [][]byte{{4}, {5}},
		},
		{
			name:        "Start from leaf",
			startNode:   node4,
			expectedIDs: [][]byte{{4}, {5}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Logf("Starting test case: %s", tc.name)

			visited_map := make(map[*b.Node]bool)
			nodes := []*b.Node{}

			t.Logf("Calling ModBFS with start node: %v", tc.startNode.Key)
			b.ModBFS(tc.startNode, visited_map, &nodes)

			t.Logf("ModBFS completed. Found %d nodes", len(nodes))

			if len(nodes) != len(tc.expectedIDs) {
				t.Errorf("Expected %d nodes, got %d", len(tc.expectedIDs), len(nodes))
			}

			for i, node := range nodes {
				t.Logf("Node %d: ID %v, IP %s", i, node.Key.Id, node.Key.Ip)

				if i >= len(tc.expectedIDs) {
					t.Errorf("Unexpected extra node at position %d: ID %v", i, node.Key.Id)
					continue
				}

				if !bytes.Equal(node.Key.Id, tc.expectedIDs[i]) {
					t.Errorf("Node at position %d: expected ID %v, got %v", i, tc.expectedIDs[i], node.Key.Id)
				}

				// Additional assertions
				assert(t, node.Left == nil, "Expected leaf node, but node %d has left child", i)
				assert(t, node.Right == nil, "Expected leaf node, but node %d has right child", i)
				assert(t, node.Key != nil, "Node %d has nil Key", i)
				assert(t, len(node.Key.Ip) > 0, "Node %d has empty IP address", i)
			}

			t.Logf("Visited map contains %d entries", len(visited_map))
			for node, visited := range visited_map {
				t.Logf("Node %v: visited = %v", node.Key, visited)
			}
		})
	}
}

// Helper function for assertions
func assert(t *testing.T, condition bool, format string, args ...interface{}) {
	if !condition {
		t.Errorf(format, args...)
	}
}

/* func TestGetKNearestNodes(t *testing.T) {
	// Create a BootStrapServer instance
	server := &b.BootStrapServer{}

	// Initialize the btree
	node1 := &b.Node{Key: &b.NodeAddr{}}
	node2 := &b.Node{Key: &b.NodeAddr{}, Parent: node1}
	node3 := &b.Node{Key: &b.NodeAddr{}, Parent: node1}
	node1.Left = node2
	node1.Right = node3
	node4 := &b.Node{Key: &b.NodeAddr{Id: []byte{1, 2, 3}, Ip: "192.168.1.1"}, Parent: node2}
	node2.Left = node4
	node5 := &b.Node{Key: &b.NodeAddr{Id: []byte{5}, Ip: "5.5.5.5"}, Parent: node3}
	node3.Right = node5

	b.SetBTreeRoot(node1)

	testCases := []struct {
		name     string
		input    *b.NodeAddr
		expected []b.NodeAddr
	}{
		{
			name:  "Get nearest nodes for existing node",
			input: node4.Key,
			expected: []b.NodeAddr{
				{Id: []byte{1, 2, 3}, Ip: "192.168.1.1"},
				{Id: []byte{5}, Ip: "5.5.5.5"},
			},
		},
		{
			name:  "Get nearest nodes for non-existing node",
			input: &b.NodeAddr{Id: []byte{5, 6, 7}, Ip: "192.168.1.5"},
			expected: []b.NodeAddr{
				{Id: []byte{5}, Ip: "5.5.5.5"},
				{Id: []byte{1, 2, 3}, Ip: "192.168.1.1"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var reply []b.NodeAddr
			err := server.GetKNearestNodes(tc.input, &reply)

			if err != nil {
				t.Fatalf("GetKNearestNodes returned an error: %v", err)
			}

			t.Logf("Reply: %v", reply)

			if len(reply) != len(tc.expected) {
				t.Errorf("Expected %d nodes in reply, got %d", len(tc.expected), len(reply))
			}
			err = server.GetKNearestNodes(tc.input, &reply)

			if err != nil {
				t.Fatalf("GetKNearestNodes returned an error: %v", err)
			}

			if len(reply) != len(tc.expected) {
				t.Errorf("Expected %d nodes in reply, got %d", len(tc.expected), len(reply))
			}

			// Check if all expected nodes are in the reply
			for _, expectedNode := range tc.expected {
				found := false
				for _, replyNode := range reply {
					if reflect.DeepEqual(expectedNode, replyNode) {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected node %v not found in reply", expectedNode)
				}
			}
			fmt.Println("Tree structure:")
		})
	}
} */

/* func TestRegisterNewNode(t *testing.T) {
	server := b.BootStrapServer{}
	var reply string
	test_cases := []struct {
		name     string
		input    *b.NodeAddr
		expected error
	}{
		{
			name:     "Register first node",
			input:    &b.NodeAddr{Id: []byte{1, 2, 3}, Ip: "192.168.1.1"},
			expected: nil,
		},
		{
			name:     "Register second node",
			input:    &b.NodeAddr{Id: []byte{1, 2, 4}, Ip: "192.168.1.2"},
			expected: nil,
		},
		{
			name:     "Register third node",
			input:    &b.NodeAddr{Id: []byte{1, 2, 3}, Ip: "192.168.1.1"},
			expected: nil,
		},
	}
	for _, tc := range test_cases {
		t.Run(tc.name, func(t *testing.T) {
			reply = ""
			err := server.RegisterNewNode(tc.input, &reply)
			if err != nil {
				t.Fatalf("RegisterNewNode returned an error: %v", err)
			}
			assert(t, reply == "Success", "Node registration failed")
		})
	}
} */

func TestGetChild(t *testing.T) {
	// Test cases
	tests := []struct {
		name     string
		node     *b.Node
		dir      int
		expected *b.Node
	}{
		{
			name:     "Nil node",
			node:     nil,
			dir:      0,
			expected: nil,
		},
		{
			name:     "Get left child, doesn't exist",
			node:     &b.Node{},
			dir:      0,
			expected: &b.Node{},
		},
		{
			name:     "Get right child, doesn't exist",
			node:     &b.Node{},
			dir:      1,
			expected: &b.Node{},
		},
		{
			name:     "Get existing left child",
			node:     &b.Node{Left: &b.Node{}},
			dir:      0,
			expected: &b.Node{},
		},
		{
			name:     "Get existing right child",
			node:     &b.Node{Right: &b.Node{}},
			dir:      1,
			expected: &b.Node{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.node.GetChild(tt.dir)

			if tt.node == nil {
				if result != nil {
					t.Errorf("Expected nil, got non-nil result")
				}
				return
			}

			if result == nil {
				t.Errorf("Expected non-nil result, got nil")
				return
			}

			if tt.dir == 0 {
				if tt.node.Left != result {
					t.Errorf("Expected left child to be set")
				}
				if result.Parent != tt.node {
					t.Errorf("Expected parent to be set for left child")
				}
			} else {
				if tt.node.Right != result {
					t.Errorf("Expected right child to be set")
				}
				if result.Parent != tt.node {
					t.Errorf("Expected parent to be set for right child")
				}
			}
		})
	}
}

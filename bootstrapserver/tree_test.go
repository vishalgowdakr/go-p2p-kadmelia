package bootstrapserver_test

import (
	"fmt"
	b "go-p2p/bootstrapserver"
	"testing"
)

func CompareSlice(actual, expected []int) bool {
	for i, value := range expected {
		if actual[i] != value {
			return false
		}
	}
	return true
}

func TestStrToInt(t *testing.T) {
	TestCases := []struct {
		name     string
		input    string
		expected []int
	}{
		{
			name:     "tc1",
			input:    "101010",
			expected: []int{1, 0, 1, 0, 1, 0},
		},
		{
			name:     "tc2",
			input:    "1011010",
			expected: []int{1, 0, 1, 1, 0, 1, 0},
		},
	}
	for _, tc := range TestCases {
		t.Run(tc.name, func(t *testing.T) {
			expected := b.StrToIntArr(tc.input)
			if !CompareSlice(expected, tc.expected) {
				t.Errorf("Actual output did not match expected")
			}
		})
	}
}

func TestTreeInsert(t *testing.T) {
	nodeAddr1 := b.NodeAddr{
		Id: "101010",
		Ip: "192.168.0.1",
	}
	node1 := b.NewNode(nodeAddr1)
	nodeAddr2 := b.NodeAddr{
		Id: "1011010",
		Ip: "192.168.0.2",
	}
	node2 := b.NewNode(nodeAddr2)
	TestCases := []struct {
		name     string
		id       []int
		node     *b.Node
		expected bool
	}{
		{
			name:     "tc1",
			id:       []int{1, 0, 1, 0, 1, 0},
			node:     node1,
			expected: true,
		},
		{
			name:     "tc2",
			id:       []int{1, 0, 1, 1, 0, 1, 0},
			node:     node2,
			expected: true,
		},
	}
	tree := b.NewTree()
	for _, tc := range TestCases {
		t.Run(tc.name, func(t *testing.T) {
			if !tree.Insert(tc.node) {
				t.Fatalf("insertion failed")
			}
		})
	}
	tree.Print()
}

func TestGetKNearestNodes(t *testing.T) {
	nodeAddr1 := b.NodeAddr{
		Id: "101010",
		Ip: "192.168.0.1",
	}
	node1 := b.NewNode(nodeAddr1)
	nodeAddr2 := b.NodeAddr{
		Id: "1011010",
		Ip: "192.168.0.2",
	}
	node2 := b.NewNode(nodeAddr2)
	tree := b.NewTree()
	tree.Insert(node1)
	tree.Insert(node2)

	TestCases := []struct {
		name     string
		input    string
		expected []b.NodeAddr
	}{
		{
			name:  "tc1",
			input: "101010",
			expected: []b.NodeAddr{
				{
					Id: "101010",
					Ip: "192.168.0.1",
				},
				{
					Id: "1011010",
					Ip: "192.168.0.2",
				},
			},
		},
		{
			name:  "tc2",
			input: "1011010",
			expected: []b.NodeAddr{
				{
					Id: "1011010",
					Ip: "192.168.0.2",
				},
				{
					Id: "101010",
					Ip: "192.168.0.1",
				},
			},
		},
	}

	for _, tc := range TestCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tree.GetKNearestNodes(tc.input)
			if len(result) != len(tc.expected) {
				t.Errorf("Expected %d nodes, got %d", len(tc.expected), len(result))
			}
			for i, node := range result {
				if node != tc.expected[i] {
					t.Errorf("Expected node %v, got %v", tc.expected[i], node)
				}
				fmt.Print(node)
			}
		})
	}
}

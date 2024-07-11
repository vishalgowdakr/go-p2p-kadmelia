package tree_test

import (
	"fmt"
	t "go-p2p/tree"
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

func TestStrToInt(test *testing.T) {
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
		test.Run(tc.name, func(test *testing.T) {
			expected := t.StrToIntArr(tc.input)
			if !CompareSlice(expected, tc.expected) {
				test.Errorf("Actual output did not match expected")
			}
		})
	}
}

func TestTreeInsert(test *testing.T) {
	nodeAddr1 := t.NodeAddr{
		Id: "101010",
		Ip: "192.168.0.1",
	}
	node1 := t.NewNode(nodeAddr1)

	nodeAddr2 := t.NodeAddr{
		Id: "1011010",
		Ip: "192.168.0.2",
	}
	node2 := t.NewNode(nodeAddr2)

	nodeAddr3 := t.NodeAddr{
		Id: "11111",
		Ip: "192.168.0.3",
	}
	node3 := t.NewNode(nodeAddr3)

	nodeAddr4 := t.NodeAddr{
		Id: "00000",
		Ip: "192.168.0.4",
	}
	node4 := t.NewNode(nodeAddr4)

	TestCases := []struct {
		name     string
		id       []int
		node     *t.Node
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
		{
			name:     "tc4_all_ones",
			id:       []int{1, 1, 1, 1, 1},
			node:     node3,
			expected: true,
		},
		{
			name:     "tc5_all_zeros",
			id:       []int{0, 0, 0, 0, 0},
			node:     node4,
			expected: true,
		},
	}
	tree := t.NewTree()
	for _, tc := range TestCases {
		test.Run(tc.name, func(test *testing.T) {
			if !tree.Insert(tc.node) {
				test.Fatalf("insertion failed")
			}
		})
	}
	tree.Print()
}

func TestGetKNearestNodes(test *testing.T) {
	nodeAddr1 := t.NodeAddr{
		Id: "101010",
		Ip: "192.168.0.1",
	}
	node1 := t.NewNode(nodeAddr1)
	nodeAddr2 := t.NodeAddr{
		Id: "1011010",
		Ip: "192.168.0.2",
	}
	node2 := t.NewNode(nodeAddr2)
	tree := t.NewTree()
	tree.Insert(node1)
	tree.Insert(node2)

	TestCases := []struct {
		name     string
		input    string
		expected []t.NodeAddr
	}{
		{
			name:  "tc1",
			input: "101010",
			expected: []t.NodeAddr{
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
			expected: []t.NodeAddr{
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
		test.Run(tc.name, func(test *testing.T) {
			result := tree.GetKNearestNodes(tc.input)
			if len(result) != len(tc.expected) {
				test.Errorf("Expected %d nodes, got %d", len(tc.expected), len(result))
			}
			for i, node := range result {
				if node != tc.expected[i] {
					test.Errorf("Expected node %v, got %v", tc.expected[i], node)
				}
				fmt.Print(node)
			}
		})
	}
}

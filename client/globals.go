package client

import (
	"go-p2p/tree"
	"sync"
)

var (
	myNodeID       string
	myRoutingTable RoutingTable
	myAddrInfo     tree.NodeAddr
	initOnce       sync.Once
	mu             sync.RWMutex
)

// InitializeGlobals should be called once when your node starts up
func InitializeGlobals(routingTable RoutingTable, node tree.NodeAddr) {
	initOnce.Do(func() {
		mu.Lock()
		defer mu.Unlock()
		myRoutingTable = routingTable
		myAddrInfo = node
	})
}

func GetMyAddrInfo() tree.NodeAddr {
	mu.RLock()
	defer mu.RUnlock()
	return myAddrInfo
}

func GetMyRoutingTable() RoutingTable {
	mu.RLock()
	defer mu.RUnlock()
	return myRoutingTable
}

func UpdateRoutingTable(newRoutingTable RoutingTable) {
	mu.Lock()
	defer mu.Unlock()
	myRoutingTable = newRoutingTable
}

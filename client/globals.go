package client

import (
	"sync"

	"github.com/multiformats/go-multiaddr"
)

var (
	myNodeID       string
	myRoutingTable RoutingTable
	myAddrInfo     []multiaddr.Multiaddr
	initOnce       sync.Once
	mu             sync.RWMutex
)

// InitializeGlobals should be called once when your node starts up
func InitializeGlobals(nodeID string, routingTable RoutingTable, addrInfo []multiaddr.Multiaddr) {
	initOnce.Do(func() {
		mu.Lock()
		defer mu.Unlock()
		myNodeID = nodeID
		myRoutingTable = routingTable
		myAddrInfo = addrInfo
	})
}

func GetMyAddrInfo() []multiaddr.Multiaddr {
	mu.RLock()
	defer mu.RUnlock()
	return myAddrInfo
}

func GetMyNodeID() (string, RoutingTable) {
	mu.RLock()
	defer mu.RUnlock()
	return myNodeID, myRoutingTable
}

func UpdateRoutingTable(newRoutingTable RoutingTable) {
	mu.Lock()
	defer mu.Unlock()
	myRoutingTable = newRoutingTable
}

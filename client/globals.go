package client

import (
	"sync"
)

var (
	myNodeID       string
	myRoutingTable RoutingTable
	initOnce       sync.Once
	mu             sync.RWMutex
)

// InitializeGlobals should be called once when your node starts up
func InitializeGlobals(nodeID string, routingTable RoutingTable) {
	initOnce.Do(func() {
		mu.Lock()
		defer mu.Unlock()
		myNodeID = nodeID
		myRoutingTable = routingTable
	})
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

# go-p2p-kadmelia

Implementation roadmap

Bootstrap Server
	[x] Binary Tree based on ID of node (Tree Insertion, Node searching and K nearest nodes searching)
	[x] Testing of Tree Insertion, Node searching and K nearest nodes searching
	[x] RPC endpoints ( RegisterNewNode, GetKNearestNodes )
	[x] Testing of RPC Endpoints

Node with Routing Table
	[ ] Node implementation using lib-p2p
	[ ] Routing table using K Bucket algorithm
	[ ] Refreshing routing table periodically to keep it healthy
	[ ] Leaving the network gracefully

File Sharing using the overlay network
	[ ] File sharding
	[ ] Get the table of target nodes with the help of bootstrap server
	[ ] Share the file with target nodes with acknowledgment
	
Torrent File Building
	[ ]
	
CLI or Android app

build:
	docker stop node
	docker rm node
	docker stop bootstrapserver
	docker rm bootstrapserver
	docker build -f docker/Dockerfile.bootstrapserver -t bootstrapserver .
	docker build -f docker/Dockerfile.node -t node .
	docker run -itd --name bootstrapserver --network p2p bootstrapserver
	docker run -itd --name node --network p2p node

build_nc:
	docker build -f docker/Dockerfile.bootstrapserver -t bootstrapserver .
	docker build -f docker/Dockerfile.node -t node .
	docker run -itd --name bootstrapserver --network p2p bootstrapserver
	docker run -itd --name node --network p2p node

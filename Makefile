build:
	docker stop node
	docker rm node
	docker stop bootstrapserver
	docker rm bootstrapserver
	docker build -f docker/Dockerfile.bootstrapserver -t bootstrapserver .
	docker build -f docker/Dockerfile.node -t node .
	docker run -itd --name bootstrapserver --network p2p docker_bootstrap
	docker run -itd --name cli1 --network p2p docker_cli1
	docker run -itd --name cli2 --network p2p docker_cli2
	docker exec -it node sh

build_nc:
	docker build -f docker/Dockerfile.bootstrapserver -t bootstrapserver .
	docker build -f docker/Dockerfile.node -t node .
	docker run -itd --name bootstrapserver --network p2p bootstrapserver
	docker run -itd --name node --network p2p node
	docker exec -it node sh

NETWORK_NAME := p2p

build:
	@docker network inspect $(NETWORK_NAME) >/dev/null 2>&1 || docker network create $(NETWORK_NAME)
	@docker rmi $$(docker images -f "dangling=true" -q) --force 2>/dev/null || true
	@docker stop bootstrapserver 2>/dev/null || true
	@docker stop cli1 2>/dev/null || true
	@docker stop cli2 2>/dev/null || true
	@docker stop $$(docker ps -a -q --filter name='^(node[0-9]+)$$' --format="{{.Names}}") 2>/dev/null || true
	@docker rm bootstrapserver 2>/dev/null || true
	@docker rm cli1 2>/dev/null || true
	@docker rm cli2 2>/dev/null || true
	@docker rm $$(docker ps -a -q --filter name='^(node[0-9]+)$$' --format="{{.Names}}") 2>/dev/null || true
	docker build -f docker/Dockerfile.bootstrapserver -t bootstrapserver .
	docker build -f docker/Dockerfile.node -t node .
	docker build -f docker/Dockerfile.cli -t cli .
	docker run -d --name bootstrapserver --network $(NETWORK_NAME) bootstrapserver
	docker run -d --name cli1 --network $(NETWORK_NAME) cli
	docker run -d --name cli2 --network $(NETWORK_NAME) cli
	@for i in 1 2 3 4 5 6 7 8 9 10; do \
		docker run -d --name node$$i --network $(NETWORK_NAME) node; \
	done
	@echo "All containers started. Use 'docker exec -it node1 sh' to access a node."

stop:
	@docker stop $$(docker ps -a -q --filter name='^(bootstrapserver|node[0-9]|cli[1-2]+)$$' --format="{{.Names}}") 2>/dev/null || true
	@docker rm $$(docker ps -a -q --filter name='^(bootstrapserver|node[0-9]|cli[1-2]+)$$' --format="{{.Names}}") 2>/dev/null || true
	@docker network rm $(NETWORK_NAME) 2>/dev/null || true

.PHONY: build stop

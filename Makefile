# Raspberry Pi's reserved local IP address
SERVER_ADDRESS ?= 0.0.0.0
PORT_NUMBER ?= 50051

# Local testing
build-local:
	go mod download

build-proto:
	protoc --go-grpc_out=. --go-grpc_opt=paths=source_relative --go_out=. --go_opt=paths=source_relative api/types/api.proto

test-local:
	./script/start-server-local.sh

# Build Docker images for all servers
build-home-server:
	docker build -t home-server:latest -f cmd/home-server/Dockerfile .

build-pi-agent:
	docker build -t pi-agent:latest -f cmd/pi-agent/Dockerfile .

build-pi-agent-controller:
	docker build -t pi-agent-controller:latest -f cmd/pi-agent-controller/Dockerfile .

build-all: build-home-server build-pi-server build-pi-agent-controller

# Run Docker containers for testing
run-home-server:
	docker run --rm -p 5005:5005 home-server:latest

run-pi-agent:
	docker run --rm -p 50051:50051 pi-agent:latest

run-pi-agent-controller:
	docker run --rm -p 50052:50052 pi-agent-controller:latest

run-all:
	docker container prune -f; docker-compose up --build --remove-orphans

# Stop all services
stop-all:
	docker compose down

# Clean up dangling Docker images
clean:
	docker system prune -af --volumes
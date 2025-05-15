# Raspberry Pi's reserved local IP address
SERVER_ADDRESS ?= 0.0.0.0
PORT_NUMBER ?= 50051

rebuild-proto:
	protoc --go-grpc_out=. --go-grpc_opt=paths=source_relative --go_out=. --go_opt=paths=source_relative api/types/api.proto

test-local:
	./script/start-server-local.sh

# Build Docker images for all servers
build-home-server:
	docker build -t home-server:latest -f cmd/home-server/Dockerfile .

build-pi-server:
	docker build -t pi-agent:latest -f cmd/pi-agent/Dockerfile .

build-queue-ingestor-server:
	docker build -t queue-ingestor-server:latest -f cmd/queue-ingestor-server/Dockerfile .

build-all: build-home-server build-pi-server build-queue-ingestor-server

# Run Docker containers for testing
run-home-server:
	docker run --rm -p 5005:5005 home-server:latest

run-pi-server:
	docker run --rm -p 50051:50051 pi-agent:latest

run-queue-ingestor-server:
	docker run --rm -p 50052:50052 queue-ingestor-server:latest

run-all:
	docker-compose up --build

# Stop all services
stop-all:
	docker compose down

# Clean up dangling Docker images
clean:
	docker rmi -f $(docker images -f "dangling=true" -q)
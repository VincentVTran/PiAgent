# Raspberry Pi's reserved local IP address
SERVER_ADDRESS ?= 0.0.0.0
PORT_NUMBER ?= 50051

rebuild-proto:
	protoc --go-grpc_out=. --go-grpc_opt=paths=source_relative --go_out=. --go_opt=paths=source_relative api/types/api.proto

test-local:
	./script/start-server-local.sh

start-client:
	go run example/test/client_main.go --server-address $(SERVER_ADDRESS):$(PORT_NUMBER)

# Build the Docker image
build:
	docker build -t homeserver:latest .

# Run the Docker container for testing
run:
	docker run --rm -p 50051:50051 homeserver:latest

# Clean up dangling Docker images
clean:
	docker rmi -f $(docker images -f "dangling=true" -q)
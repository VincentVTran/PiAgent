# Raspberry Pi's reserved local IP address
SERVER_ADDRESS ?= 0.0.0.0
PORT_NUMBER ?= 50051

# [Dev] Work environment setup
build-local:
	go mod download

build-proto:
	protoc --go-grpc_out=. --go-grpc_opt=paths=source_relative --go_out=. --go_opt=paths=source_relative api/types/api.proto

# [Local w/o containers] Local testing commands
test-local:
	./script/start-server-local.sh

# [Local w/ containers] Local image building commands
build-pi-controller-websocket:
	docker build -t pi-controller-websocket:latest -f cmd/pi-controller-websocket/Dockerfile .

build-pi-controller:
	docker build -t pi-controller:latest -f cmd/pi-controller/Dockerfile .

build-pi-controller-processor:
	docker build -t pi-controller-processor:latest -f cmd/pi-controller-processor/Dockerfile .

build-all: build-pi-controller-websocket build-pi-server build-pi-controller-processor

# [Local w/ containers] Local testing with containers
run-pi-controller-websocket:
	docker run --rm -p 5005:5005 pi-controller-websocket:latest

run-pi-controller:
	docker run --rm -p 50051:50051 pi-controller:latest

run-pi-controller-processor:
	docker run --rm -p 50052:50052 pi-controller-processor:latest

run-all:
	docker container prune -f; docker-compose up --build --remove-orphans

# [Prod] Installation commands
deploy-agent: 
	./script/install-pi-controller.sh

# Stop all services
stop-all:
	docker compose down

# Clean up dangling Docker images
clean:
	docker system prune -af --volumes
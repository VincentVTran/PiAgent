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
build-pi-agent-controller:
	docker build -t pi-agent-controller:latest -f cmd/pi-agent-controller/Dockerfile .

build-pi-agent:
	docker build -t pi-agent:latest -f cmd/pi-agent/Dockerfile .

build-pi-agent-controller-processor:
	docker build -t pi-agent-controller-processor:latest -f cmd/pi-agent-controller-processor/Dockerfile .

build-all: build-pi-agent-controller build-pi-server build-pi-agent-controller-processor

# [Local w/ containers] Local testing with containers
run-pi-agent-controller:
	docker run --rm -p 5005:5005 pi-agent-controller:latest

run-pi-agent:
	docker run --rm -p 50051:50051 pi-agent:latest

run-pi-agent-controller-processor:
	docker run --rm -p 50052:50052 pi-agent-controller-processor:latest

run-all:
	docker container prune -f; docker-compose up --build --remove-orphans

# [Prod] Installation commands
deploy-agent: 
	./script/install-pi-agent.sh

# Stop all services
stop-all:
	docker compose down

# Clean up dangling Docker images
clean:
	docker system prune -af --volumes
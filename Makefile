# Raspberry Pi's reserved local IP address
SERVER_ADDRESS ?= 192.168.1.124
# Default port number
PORT_NUMBER ?= 50051

rebuild-proto:
	protoc --go-grpc_out=. --go-grpc_opt=paths=source_relative --go_out=. --go_opt=paths=source_relative proto/api.proto

start-server:
	go run server/main.go --port $(PORT_NUMBER)

start-client:
	go run client/main.go --server-address $(SERVER_ADDRESS):$(PORT_NUMBER)
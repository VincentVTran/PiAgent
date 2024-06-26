## HomeServer

#### To init mod
```go mod init github.com/vincentvtran/genagent```

#### Installing dependencies
- Used ```go get```
- Followed the installation url for [grpc](https://grpc.io/docs/languages/go/quickstart/)
    - Installed gRPC compiler
    - Installed golang
    - Added client and server directory
    - Added .proto file in proto directory

#### Generating proto into go file
- For gRPC framework proto: 
```
protoc --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/agent.proto
```
- For response/request object setter and getter proto: ```protoc --go-grpc_out=. --go-grpc_opt=paths=source_relative --go_out=. --go_opt=paths=source_relative proto/api.proto```


#### Running gRPC server
- Running server: ```go run server/main.go```
- Running client: ```go run client/main.go```
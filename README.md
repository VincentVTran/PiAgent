## HomeServer
#### Component Overview
- pi-gateway = gRPC server hosted on a k8 cluster that interacts with pi-controller
- pi-controller = gRPC client that is installed onto raspberry pi. Purpose is to listen to controller and handles on-system functionalities

#### To init mod (Fresh build)
```go mod init github.com/vincentvtran/genagent```

#### To build from pre-existing mod file
```make build-local```

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
make build-proto
```

#### Running application locally
- Running with local compiler: ```make test-local```
- Running with local container: ```make run-all```
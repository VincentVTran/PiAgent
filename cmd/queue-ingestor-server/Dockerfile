# Use the official Golang image as a base image
FROM golang:1.22

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules manifests and download dependencies
COPY ../../../go.mod ../../../go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY ../../../ ./

# Build the gRPC server
RUN go build -o queue-ingestor-server ./cmd/queue-ingestor-server/main.go

# Expose the gRPC server port
EXPOSE 50052

# Command to run the gRPC server
CMD ["./queue-ingestor-server"]

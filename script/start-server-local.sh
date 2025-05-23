#!/usr/bin/env bash
# Function to find an available port
make_port_available() {
    PORT=$1
    PID=$(lsof -t -i:$PORT)
    if [ -n "$PID" ]; then
        echo "Killing process $PID on port $PORT"
        kill -9 $PID
    else
        echo $PORT
    fi
}

# # Function to clean up background processes
# cleanup() {
#   echo "Stopping background processes..."
#   kill $HOME_PID $PI_PID
# }
# trap cleanup EXIT

# Find available ports
HOME_PORT=$(make_port_available 5005)  # Capture the output of the function
PI_PORT=$(make_port_available 50051)  # Capture the output of the function

# Start pi-controller-websocket in the background
go run cmd/pi-controller-websocket/main.go --stage="local" --port=$HOME_PORT &
HOME_PID=$!
echo "Home server started with PID $HOME_PID on port $HOME_PORT"

# Start pi-controller in the background
go run cmd/pi-controller/main.go --stage="local" --port=$PI_PORT &
PI_PID=$!
echo "Pi server started with PID $PI_PID on port $PI_PORT"

# Start the queue-ingestor-service in the foreground
go run cmd/pi-controller-processor/main.go --stage="local" --server-address dev-desktop.vt:$PI_PORT
echo "Queue ingestor service started on dev-desktop.vt:$PI_PORT"

# Wait for the foreground process to complete
wait
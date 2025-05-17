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

PI_PORT=$(make_port_available 50051)  # Capture the output of the function

# Start pi-controller in the background
go run cmd/pi-controller/main.go --stage="local" --port=$PI_PORT &
PI_PID=$!
echo "Pi server started with PID $PI_PID on port $PI_PORT"


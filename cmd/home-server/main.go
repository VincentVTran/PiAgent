package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	amqp "github.com/rabbitmq/amqp091-go"
	pb "github.com/vincentvtran/homeserver/api/types"
)

var (
	port       = flag.Int("port", 5005, "The gRPC server port")
	stage      = flag.String("stage", "local", "Stage for RabbitMQ URL (e.g., local or production)")
	rabbitURL  string
	exchange   = "amq.direct"
	routingKey = "pi-queue-key"
	upgrader   = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
)

// Payload defines the message structure
type Payload struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}

func determineRabbitMQURL() {
	switch *stage {
	case "local":
		rabbitURL = "amqp://admin:admin-ui-password@192.168.2.4:5672/"
		log.Println("Using local RabbitMQ URL")
	case "prod":
		rabbitURL = "amqp://admin:admin-ui-password@rabbitmq.rabbitmq.service.cluster.local:5672/"
		log.Println("Using cluster RabbitMQ URL")
	default:
		rabbitURL = "amqp://admin:admin-ui-password@rabbitmq.rabbitmq.service.cluster.local:5672/"
		log.Printf("Using RabbitMQ URL for stage '%s'", *stage)
	}
}

func publishToExchange(url, exchange, routingKey string, payload Payload) error {
	// Connect to RabbitMQ
	conn, err := amqp.Dial(url)
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	// Open a channel
	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open a channel: %v", err)
	}
	defer ch.Close()

	// Declare an exchange
	err = ch.ExchangeDeclare(
		exchange, // name
		"direct", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare an exchange: %v", err)
	}

	// Serialize the payload
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %v", err)
	}

	// Publish the message to the exchange
	err = ch.Publish(
		exchange,   // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish a message: %v", err)
	}

	log.Printf("Message published to exchange %s with routing key %s: %v", exchange, routingKey, payload)
	return nil
}

// WebSocket handler
func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}
	defer conn.Close()

	log.Println("WebSocket connection established")

	for {
		// Read message from client
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			break
		}
		log.Printf("WebSocket received message: %s", message)
		publishToExchange(rabbitURL, exchange, routingKey, Payload{Message: string(message)})
		// Echo the message back to the client
		err = conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Printf("Error writing message: %v", err)
			break
		}
	}
}

// gRPC server implementation
type server struct {
	pb.UnimplementedHomeServiceServer
	version string
}

func (s *server) Invoke(ctx context.Context, in *pb.OperationRequest) (*pb.OperationResponse, error) {
	log.Println("Received gRPC connection from client")
	log.Println("Client request: ", in)

	// Create a payload and publish it to the exchange
	payload := Payload{
		ID:      "12345",
		Message: "Hello from home-server",
	}
	err := publishToExchange(rabbitURL, exchange, routingKey, payload)
	if err != nil {
		log.Printf("Error publishing to RabbitMQ exchange: %v", err)
		return nil, err
	}

	return &pb.OperationResponse{ApiVersion: s.version}, nil
}

func startWebSocketServer() {
	http.HandleFunc("/ws", handleWebSocket)
	log.Printf("WebSocket server listening on port %d", *port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
	if err != nil {
		log.Fatalf("Failed to start WebSocket server: %v", err)
	}
}

func main() {
	flag.Parse()

	// Determine RabbitMQ URL based on stage
	determineRabbitMQURL()

	// Start WebSocket server
	startWebSocketServer()
}

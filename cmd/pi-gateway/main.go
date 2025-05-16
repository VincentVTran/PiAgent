package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	pb "github.com/vincentvtran/pi-controller/api/types"
	model "github.com/vincentvtran/pi-controller/pkg/model"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	addr      = flag.String("server-address", "dev-desktop.vt:50051", "the address to connect to")
	stage     = flag.String("stage", "local", "Stage for RabbitMQ URL (e.g., local or production)")
	queueName = flag.String("queue", "pi-queue", "RabbitMQ queue name")
	rabbitURL string
	config    model.ApplicationConfig
)

func loadConfig() {
	file, err := os.Open("config/application-config.json")
	if err != nil {
		log.Fatalf("Failed to open config file: %v", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		log.Fatalf("Failed to decode config file: %v", err)
	}
}

func determineRabbitMQURL() {
	switch *stage {
	case "local":
		rabbitURL = config.Local.RabbitMQLink
		log.Println("Using local RabbitMQ URL")
	case "prod":
		rabbitURL = config.Prod.RabbitMQLink
		log.Println("Using cluster RabbitMQ URL")
	default:
		log.Fatalf("RabbitMQ URL for stage '%s' not found in config", *stage)
	}
}

func consumeFromRabbitMQ(url, queue string) error {
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

	// Consume messages from the queue
	msgs, err := ch.Consume(
		queue,        // queue
		"pi-gateway", // consumer
		true,         // auto-ack
		false,        // exclusive
		false,        // no-local (deprecated in amqp091-go)
		false,        // no-wait
		nil,          // args
	)
	if err != nil {
		return fmt.Errorf("failed to register a consumer: %v", err)
	}

	// Process messages
	forever := make(chan bool)
	go func() {
		for d := range msgs {
			var p model.PiOperation
			err := json.Unmarshal(d.Body, &p)
			if err != nil {
				log.Printf("Error unmarshalling message: %v", err)
				continue
			}
			log.Printf("[Queue Ingestor] Received message: %v", p)
		}
	}()

	log.Println("Waiting for messages. To exit press CTRL+C")
	<-forever
	return nil
}

func main() {
	flag.Parse()

	// Load configuration
	loadConfig()

	// Set up a connection to the server.
	conn, err := grpc.NewClient(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewPiAgentControllerClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	request := &pb.StreamRequest{
		Parameters: &pb.StreamParameter{
			Enable: true,
		},
	}
	r, err := c.ConfigureStream(ctx, request)
	if err != nil {
		log.Fatalf("Could not reach server: %v", err)
	}
	log.Printf("Server is currently set to version: %s", r.ApiVersion)

	// Determine RabbitMQ URL based on stage
	determineRabbitMQURL()

	log.Printf("Connecting to RabbitMQ at %s and consuming from queue %s", rabbitURL, *queueName)
	err = consumeFromRabbitMQ(rabbitURL, *queueName)
	if err != nil {
		log.Fatalf("Error consuming from RabbitMQ: %v", err)
	}
}

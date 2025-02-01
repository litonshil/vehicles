package conn

import (
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

// Global connection
var rc *amqp.Connection
var ch *amqp.Channel

func InitRabbitMQ() {
	connectString := "amqps://avxpoguo:Da6pggbTCzcN6BiyrTnva-7549c5dU89@fuji.lmq.cloudamqp.com/avxpoguo" //"amqp://guest:guest@localhost:15672/"
	fmt.Println("Connecting to RabbitMQ at", connectString)

	rc, err := amqp.Dial(connectString)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}

	ch, err = rc.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %s", err)
	}

	log.Println("Connected to RabbitMQ")
}

// CloseRabbitMQ Close RabbitMQ connection
func CloseRabbitMQ() {
	_ = ch.Close()
	_ = rc.Close()
	log.Println("RabbitMQ connection closed")
}

// PublishMessage Publish a message to RabbitMQ
func PublishMessage(exchange, routingKey string, message interface{}) error {
	err := ch.ExchangeDeclare(
		exchange, // Exchange name
		"direct", // Type
		true,     // Durable
		false,    // Auto-deleted
		false,    // Internal
		false,    // No-wait
		nil,      // Arguments
	)
	if err != nil {
		return err
	}

	// Declare a queue and bind it to the exchange
	queue, err := ch.QueueDeclare(
		"vehicle_queue", // Queue name
		true,            // Durable
		false,           // Auto-delete
		false,           // Exclusive
		false,           // No-wait
		nil,             // Arguments
	)
	if err != nil {
		return err
	}

	err = ch.QueueBind(
		queue.Name, // Queue name
		"",         // Routing key (ignored for fanout)
		exchange,   // Exchange name
		false,
		nil,
	)
	if err != nil {
		return err
	}

	// Convert the message to JSON
	body, err := json.Marshal(message)
	if err != nil {
		return err
	}

	// Publish the message
	err = ch.Publish(
		exchange, // Exchange
		"",       // Routing key (ignored for fanout)
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent, // Make message persistent
		},
	)
	if err != nil {
		log.Printf("Failed to publish message to %s: %v", routingKey, err)
		return err
	}

	log.Printf("Message published to exchange %s: %+v\n", exchange, message)
	return nil
}

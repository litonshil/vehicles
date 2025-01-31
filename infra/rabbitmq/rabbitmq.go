package rabbitmq

import (
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

type RabbitMQ struct {
	conn *amqp.Connection
}

var rabbitMQ *RabbitMQ

func InitRabbitMQ() *RabbitMQ {
	connectString := "amqp://guest:guest@localhost:5672/"
	fmt.Println("Connecting to RabbitMQ at", connectString)

	rc, err := amqp.Dial(connectString)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}

	log.Println("Connected to RabbitMQ")
	rabbitMQ = &RabbitMQ{rc}
	return rabbitMQ
}

func (r *RabbitMQ) Close() {
	_ = r.conn.Close()
}

func RMQ() *RabbitMQ {
	return rabbitMQ
}

func (r *RabbitMQ) Publish(queueName string, message interface{}) error {
	ch, err := r.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	// Declare the queue (this will create the queue if it doesn't exist)
	_, err = ch.QueueDeclare(
		queueName, // name of the queue
		true,      // durable (the queue will survive server restarts)
		false,     // delete when unused
		false,     // exclusive (used by only one connection)
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return err
	}

	body, err := json.Marshal(message)
	if err != nil {
		return err
	}

	err = ch.Publish(
		"",        // exchange
		queueName, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent, // Make message persistent
		},
	)
	if err != nil {
		return err
	}

	log.Printf("Published message to %s: %s", queueName, body)
	return nil
}

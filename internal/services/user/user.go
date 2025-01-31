package user

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"vehicles/infra/conn"
	"vehicles/internal/domain"
	"vehicles/types"
)

type UserService struct {
	repo        domain.UserRepo
	redisClient *conn.CacheClient
}

func NewUserService(userRepo domain.UserRepo, redisClient *conn.CacheClient) *UserService {
	return &UserService{
		repo:        userRepo,
		redisClient: redisClient,
	}
}

func (service *UserService) CreateUser(ctx context.Context, userReq types.UserReq) error {

	if userReq.Email == "" || userReq.FirstName == "" || userReq.LastName == "" {
		return errors.New("missing required fields")
	}

	user := domain.User{
		ID: primitive.NewObjectID(),
		Profile: domain.Profile{
			Email:     userReq.Email,
			FirstName: userReq.FirstName,
			LastName:  userReq.LastName,
		},
	}

	_ = PublishUserCreateEvent(user)

	// Call repository to save the user
	//err := service.repo.CreateUser(user)
	//if err != nil {
	//	return err
	//}
	return nil
}
func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func checkQueueExists(ch *amqp.Channel, queueName string) error {
	_, err := ch.QueueDeclarePassive(
		queueName, // queue name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	return err
}

func PublishUserCreateEvent(user domain.User) error {
	go ConsumeUserCreateEvent()
	con, err := amqp.Dial("amqps://avxpoguo:Da6pggbTCzcN6BiyrTnva-7549c5dU89@fuji.lmq.cloudamqp.com/avxpoguo")
	if err != nil {
		panic(err)
	}
	defer con.Close()

	ch, err := con.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// Check if the queue exists
	err = checkQueueExists(ch, "Liton")
	if err != nil {
		log.Printf("Queue does not exist: %s", err)
	} else {
		log.Println("Queue exists!")
	}

	err = ch.ExchangeDeclare(
		"logs",   // name
		"fanout", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	// Declare a queue
	q, err := ch.QueueDeclare(
		"Liton", // queue name
		true,    // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	// Bind the queue to the exchange
	err = ch.QueueBind(
		q.Name, // queue name
		"",     // routing key
		"logs", // exchange name
		false,  // no-wait
		nil,    // arguments
	)
	failOnError(err, "Failed to bind the queue")

	// Serialize user struct to JSON
	body, err := json.Marshal(user)
	if err != nil {
		log.Fatalf("Failed to serialize user: %s", err)
	}

	err = ch.Publish(
		"logs", // exchange
		"",     // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent, // make message persistent
			ContentType:  "application/json",
			Body:         body,
		})
	failOnError(err, "Failed to publish a message")

	log.Printf(" [x] Sent %s", body)

	return nil
}

func (service *UserService) GetUsers(ctx context.Context, filter domain.UserFilter) ([]domain.User, error) {

	if filter.ID != "" {
		key := fmt.Sprintf("user:%s", filter.ID)

		user := domain.User{}
		err := service.redisClient.GetStruct(key, user)
		if err != nil {
			return nil, err
		}

		return []domain.User{user}, nil
	}

	users, err := service.repo.GetUsers(filter)
	if err != nil {
		return nil, err
	}
	go func() {
		key := fmt.Sprintf("user:%s", filter.ID)
		_ = service.redisClient.SetStruct(key, users, 10)
	}()

	return users, nil
}

func ConsumeUserCreateEvent() {

	con, err := amqp.Dial("amqps://avxpoguo:Da6pggbTCzcN6BiyrTnva-7549c5dU89@fuji.lmq.cloudamqp.com/avxpoguo")
	if err != nil {
		panic(err)
	}
	defer con.Close()

	ch, err := con.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	msgs, err := ch.Consume(
		"Liton", // queue name
		"",      // consumer
		true,    // auto-ack
		false,   // exclusive
		false,   // no-local
		false,   // no-wait
		nil,     // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			var user domain.User
			if err := json.Unmarshal(d.Body, &user); err != nil {
				log.Printf("Error decoding JSON: %s", err)
				continue
			}
			log.Printf("Received a message: %v", user)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/mohammadMghi/notificationService/models"
	amqp "github.com/rabbitmq/amqp091-go"
)

type SagaError struct{
	Massage string
}
func main(){
	var sagaError SagaError
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()
	
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()
	saga, err := ch.QueueDeclare(
		"saga_notif", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	  )

	  ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	  defer cancel()
	  body, err := json.Marshal(&sagaError)



	if err != nil {
		sagaError.Massage = err.Error()
		err = ch.PublishWithContext(ctx,
			"",     // exchange
			saga.Name, // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing {
			  ContentType: "application/json",
			  Body:        []byte(body),
			})
	}
	
	q, err := ch.QueueDeclare(
	  "send_email_auth", // name
	  false,   // durable
	  false,   // delete when unused
	  false,   // exclusive
	  false,   // no-wait
	  nil,     // arguments
	)

	if err != nil {
		sagaError.Massage = err.Error()
		err = ch.PublishWithContext(ctx,
			"",     // exchange
			saga.Name, // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing {
			  ContentType: "application/json",
			  Body:        []byte(body),
			})
	}



	failOnError(err, "Failed to declare a queue")

	user, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	  )

	if err != nil {
		sagaError.Massage = err.Error()
		err = ch.PublishWithContext(ctx,
			"",     // exchange
			saga.Name, // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing {
			  ContentType: "application/json",
			  Body:        []byte(body),
			})
	}
	  failOnError(err, "Failed to register a consumer")
	  
	  var forever chan struct{}
	  
	  go func() {
		for u := range user {
			var user models.User
			err := json.Unmarshal(u.Body, &user)
			if err != nil {
				log.Println("Error:", err)
				continue
			}

			fmt.Println("Received message:", user)
		}
	  }()
	  
	  log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	  <-forever
}

func failOnError(err error, msg string) {
	if err != nil {
	  log.Panicf("%s: %s", msg, err)
	}
  }
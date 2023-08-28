package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/mohammadMghi/notificationService/models"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main(){
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()
	
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()
	
	q, err := ch.QueueDeclare(
	  "send_email_auth", // name
	  false,   // durable
	  false,   // delete when unused
	  false,   // exclusive
	  false,   // no-wait
	  nil,     // arguments
	)
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
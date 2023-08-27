package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/mohammadMghi/accountService/models"
	"github.com/mohammadMghi/accountService/repo"
	"github.com/mohammadMghi/accountService/usecase"
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
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	  )
	  failOnError(err, "Failed to declare a queue")
	  
	  ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	  defer cancel()
	  
	  body := "Hello World!"
	  err = ch.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing {
		  ContentType: "text/plain",
		  Body:        []byte(body),
		})
	  failOnError(err, "Failed to publish a message")
	  log.Printf(" [x] Sent %s\n", body)


	  http.HandleFunc("/json", myHandler)

	  http.ListenAndServe(":8080", nil)
}

func failOnError(err error, msg string) {
  if err != nil {
    log.Panicf("%s: %s", msg, err)
  }
}
func myHandler(w http.ResponseWriter, r *http.Request) {
	var  user models.User


	_  = json.NewDecoder(r.Body).Decode(&user)
 
	repo := repo.Mysql{}
	signinDomain := usecase.SigninUsecase{Mysql: repo};

	signinDomain.Signin(user)
}
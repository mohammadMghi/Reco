package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/mohammadMghi/accountService/models"
	 
	"github.com/mohammadMghi/accountService/usecase"
	amqp "github.com/rabbitmq/amqp091-go"
)


func main(){

	e := echo.New()
 

	

	  e.POST("/users", myHandler)

	  e.Logger.Fatal(e.Start(":1323"))
}

func failOnError(err error, msg string) {
  if err != nil {
    log.Panicf("%s: %s", msg, err)
  }
}
func myHandler(c echo.Context) error {
	var  user models.User


	_  = json.NewDecoder(c.Request().Body).Decode(&user)
  
	signinDomain := usecase.SigninUsecase{};

	err , u := signinDomain.Signin(user)

	if err != nil{
		return err
	}

	sendAuthRabbitMQ(u)

	return c.String(http.StatusOK , "user saved")
}


func sendAuthRabbitMQ(user  *models.User)  {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()
	body, err := json.Marshal(&user)

	if err != nil {
		failOnError(err, "Error")
	}

	q, err := ch.QueueDeclare(
		"send_email_auth", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	  )
	  failOnError(err, "Failed to declare a queue")
	  
	  ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	  defer cancel()
	  
	 
	  err = ch.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing {
		  ContentType: "application/json",
		  Body:        []byte(body),
		})
	  failOnError(err, "Failed to publish a message")
	  log.Printf(" [x] Sent %s\n", body)
	  

}

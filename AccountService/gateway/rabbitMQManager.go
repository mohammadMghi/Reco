package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/mohammadMghi/accountService/models"
	"github.com/mohammadMghi/accountService/repo"
	"github.com/rabbitmq/amqp091-go"
)

type SagaError struct{
	Message string `json:"message"`
}

type RabbitMQManager struct{
	Mysql repo.Mysql
	
}

 
func (r *RabbitMQManager)  Signin(user models.User) ( error ,*models.User){
	db, err :=  repo.ConnectToDB()
    if err != nil {
        fmt.Println("Error connecting to database:", err)
        return err , nil
    }
    // defer db.Close()

 
    result := db.Create(&user)
 
    if result.Error != nil {
        fmt.Println("Error inserting data:", result.Error)
        return err ,  nil
    }

    fmt.Println("Data inserted successfully!")

	return nil , &user
}




func (r *RabbitMQManager)RollBack(user models.User) (error , *models.User){
	conn, err := amqp091.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()
	ch, err := conn.Channel()
	db, err :=  repo.ConnectToDB()
    if err != nil {
        fmt.Println("Error connecting to database:", err)
        return err , nil
    }


	qq, err := ch.QueueDeclare(
		"saga_notif", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	  )
    // defer db.Close()
	usersRollBack, err := ch.Consume(
		qq.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	  )
 
	  go func() {
	
		for u := range usersRollBack {
			var user models.User
			err := json.Unmarshal(u.Body, &user)
			repo.RollBack(user)
			if err != nil {
				log.Println("Error:", err)
				continue
			}

			fmt.Println("Received message:", user)
		}
	  }()
 
    result := db.Delete(&user)
 
    if result.Error != nil {
        fmt.Println("Error inserting data:", result.Error)
        return err ,  nil
    }

    fmt.Println("Data inserted successfully!")

	return nil , &user
}





func (r *RabbitMQManager)SendAuthRabbitMQ(user  *models.User)  {
	conn, err := amqp091.Dial("amqp://guest:guest@localhost:5672/")
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
		amqp091.Publishing {
		  ContentType: "application/json",
		  Body:        []byte(body),
		})


		  
	  failOnError(err, "Failed to publish a message")
	  log.Printf(" [x] Sent %s\n", body)
	  

}


func failOnError(err error, msg string) {
	if err != nil {
	  log.Panicf("%s: %s", msg, err)
	}
}
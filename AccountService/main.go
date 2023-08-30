package main

import (
 
	"encoding/json"
 
 
	"net/http"
 
 
	"github.com/labstack/echo/v4"
	"github.com/mohammadMghi/accountService/gateway"
	"github.com/mohammadMghi/accountService/models"
	"github.com/mohammadMghi/accountService/repo"

	"github.com/mohammadMghi/accountService/usecase"
 
)


func main(){

	e := echo.New()
 

	

	  e.POST("/users", myHandler)

	  e.Logger.Fatal(e.Start(":1323"))
}

func myHandler(c echo.Context) error {
	var  user models.User


	_  = json.NewDecoder(c.Request().Body).Decode(&user)
	
	repo := repo.Mysql{}
  
	signinDomain := usecase.SigninUsecase{};



	mq := gateway.RabbitMQManager{Mysql:repo}

	mq.RollBack(user)

	err , u := signinDomain.Signin(user)

	if err != nil{
		return err
	}

	mq.SendAuthRabbitMQ(u)

	return c.String(http.StatusOK , "user saved")
}





type SagaError struct{
	Massage string
}
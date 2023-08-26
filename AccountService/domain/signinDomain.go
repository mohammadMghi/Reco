package domain

import (
	"github.com/mohammadMghi/accountService/models"
 
)

type SigninDomain interface{
	 Signin(user models.User) models.User 
}


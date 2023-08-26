package usecase

import (
	"github.com/mohammadMghi/accountService/domain"
	"github.com/mohammadMghi/accountService/models"
)

type SigninUsecase struct{
	domain.SigninDomain
}


func (s SigninUsecase)  Signin(user models.User) models.User{

}
package usecase

import (
	"github.com/mohammadMghi/accountService/domain"
	"github.com/mohammadMghi/accountService/models"
	"github.com/mohammadMghi/accountService/repo"
)

type SigninUsecase struct{
	SigninDomain domain.SigninDomain
	Mysql repo.Mysql
}


func (s *SigninUsecase)  Signin(user models.User)( e error , u *models.User){
	 e , u = s.SigninDomain.Signin(user)

	 

	 if e != nil{
		return e , nil
	 }

	 return nil , u
}
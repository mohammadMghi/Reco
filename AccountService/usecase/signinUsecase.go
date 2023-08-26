package usecase

import (
	"github.com/mohammadMghi/accountService/domain"
	"github.com/mohammadMghi/accountService/models"
)

type SigninUsecase struct{
	domain.SigninDomain
}


func (s *SigninUsecase)  Signin(user models.User)( e error , u *models.User){
	 e , u = s.Signin(user)

	 if e != nil{
		return e , nil
	 }

	 return nil , u
}
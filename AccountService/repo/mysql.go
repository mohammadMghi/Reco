package repo

import (
	"github.com/mohammadMghi/accountService/domain"
	"github.com/mohammadMghi/accountService/models"
)


type Mysql struct{
	domain.SigninDomain
}


func (m *Mysql)connection(){
	
} 


func (m *Mysql)  Signin(user models.User) models.User{

}
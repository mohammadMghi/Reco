package domain

import (
	"github.com/mohammadMghi/accountService/models"
 
)

type SingnupDomain interface{
	 Singup(user models.User) ( error,  models.User )
}


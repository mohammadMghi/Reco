package repo

import (
 
	"fmt"
	"gorm.io/driver/mysql"
 
	"github.com/mohammadMghi/accountService/domain"
	"github.com/mohammadMghi/accountService/models"
	"gorm.io/gorm"
)

var db *gorm.DB
type Mysql struct{
	domain.SigninDomain
}

func connectToDB() (*gorm.DB, error) {
	dsn := "user:pass@tcp(127.0.0.1:3306)/recoAccountService?charset=utf8mb4&parseTime=True&loc=Local"
	
	if db != nil{
		return db, nil
	}
	
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

 

	if err != nil{
		return nil , err
	}

	return db , err


}
 
func (m *Mysql)  Signin(user models.User) ( error ,*models.User){
	db, err := connectToDB()
    if err != nil {
        fmt.Println("Error connecting to database:", err)
        return err , nil
    }
    // defer db.Close()

    newUser := models.User{
        PhoneNumber  : user.PhoneNumber,
		Name: user.Name,
		Email: user.Email,
    }

    result := db.Create(&newUser)
 
    if result.Error != nil {
        fmt.Println("Error inserting data:", result.Error)
        return err ,  nil
    }

    fmt.Println("Data inserted successfully!")

	return nil , &user
}
package models

 

type User struct{
	ID string `gorm:"primaryKey"`
	PhoneNumber string  
	Name string 
	Email string 
}
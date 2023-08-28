package models

 

type User struct{
	ID *string `gorm:"primaryKey"`
	PhoneNumber string   `json:"phone_number"`
	Name string  `json:"name"`
	Email string  `json:"email"`
}
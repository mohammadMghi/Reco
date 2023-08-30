package models

 

type User struct{
	ID *int64 `gorm:"primaryKey"`
	PhoneNumber string   `json:"phone_number"`
	Name string  `json:"name"`
	Email string  `json:"email"`
}
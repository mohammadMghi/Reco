package models

 

type User struct{
	id string "json:'id'"
	phone_number string "json:'phone_number'"
	name string "json:'name'"
	email string "json:'email'"
}
package model

type User struct {
	ID        uint   `gorm:"primary_key"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	UserName  string `json:"user_name" gorm:"unique"`
	Email     string `json:"email" gorm:"unique"`
	Password  string `json:"-"`
}

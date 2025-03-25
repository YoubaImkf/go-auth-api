package model

type User struct {
	ID       uint   `gorm:"primary_key"`
	Name     string `json:"name" gorm:"unique"`
	Email    string `json:"email" gorm:"unique"`
	Password string `json:"-"`
}

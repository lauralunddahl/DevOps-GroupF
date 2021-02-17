package dto

import (
	database "db"
	"fmt"
)

// var db = database.DB

type User struct {
	//gorm.Model
	UserId   int
	Username string
	Email    string
	PwHash   string
}

// func initialMigration() {
// 	db.AutoMigrate(&User{})
// }

func GetUserID(username string) string {
	//initialMigration()
	user := User{}

	database.DB.Where("username = ?", username).First(&user)

	return user.Email
}

func (User) TableName() string {
	return "user"
}

func GetUsers() {
	var user []User
	database.DB.Find(&user)
	fmt.Println("{}", user)
}

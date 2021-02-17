package controller

import (
	database "db"

	"github.com/jinzhu/gorm"
)

var db *gorm.DB = database.OpenDB()

type User struct {
	//gorm.Model
	UserId   int
	Username string
	Email    string
	PwHash   string
}

func initialMigration() {
	db.AutoMigrate(&User{})
}

func GetUserID(username string) int {
	user := User{}

	db.Where("username = ?", username).First(&user)

	return user.UserId
}

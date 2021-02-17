package dto

import database "db"

type User struct {
	UserId   int
	Username string
	Email    string
	PwHash   string
}

func (User) TableName() string {
	return "user"
}

func GetUserID(username string) int {
	user := User{}
	database.DB.Where("username = ?", username).First(&user)
	return user.UserId
}

func GetUser(username string) User {
	user := User{}
	database.DB.Where("username = ?", username).First(&user)
	return user
}

func RegisterUser(username string, email string, password string) {
	user := User{Username: username, Email: email, PwHash: password}
	result := database.DB.Create(&user)
	if result.Error != nil {
		print(result.Error)
	}
}

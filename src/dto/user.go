package dto

import database "github.com/lauralunddahl/DevOps-GroupF/src/db"

type User struct {
	UserId   int
	Username string
	Email    string
	PwHash   string
	Image    string
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

func GetUserById(user_id int) User {
	user := User{}
	database.DB.Where("user_id = ?", user_id).First(&user)
	return user
}

func GetUsername(userid int) string {
	user := User{}
	database.DB.Where("user_id = ?", userid).First(&user)
	return user.Username
}

func RegisterUser(username string, email string, password string, image string) {
	user := User{Username: username, Email: email, PwHash: password, Image: image}
	result := database.DB.Create(&user)
	if result.Error != nil {
		print(result.Error)
	}
}

func GetTotalNumberOfUsers() int {
	var result int
	database.DB.Table("users").Count(&result)
	return result
}
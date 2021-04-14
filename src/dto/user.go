package dto

import (
	database "github.com/lauralunddahl/DevOps-GroupF/src/db"
	log "github.com/sirupsen/logrus"
)

type User struct {
	UserId   int
	Username string
	Email    string
	PwHash   string
	Image    string
}

func GetUserID(username string) int {
	user := User{}
	res := database.DB.Where("username = ?", username).First(&user)
	if res.Error != nil {
		log.Println("GetUserId")
		log.Error(res.Error)
		return 0
	}
	return user.UserId
}

func GetUser(username string) User {
	user := User{}
	res := database.DB.Where("username = ?", username).First(&user)
	if res.Error != nil {
		log.Println("GetUser")
		log.Error(res.Error)
	}
	return user
}

func GetUserById(user_id int) User {
	user := User{}
	res := database.DB.Where("user_id = ?", user_id).First(&user)
	if res.Error != nil {
		log.Println("GetUserById")
		log.Error(res.Error)
	}
	return user
}

func GetUsername(userid int) string {
	user := User{}
	res := database.DB.Where("user_id = ?", userid).First(&user)
	if res.Error != nil {
		log.Println("GetUsername")
		log.Error(res.Error)
		return ""
	}
	return user.Username
}

func RegisterUser(username string, email string, password string, image string) {
	user := User{Username: username, Email: email, PwHash: password, Image: image}
	result := database.DB.Create(&user)
	if result.Error != nil {
		log.Println("RegisterUser")
		log.Error(result.Error)
	}
}

func GetTotalNumberOfUsers() int {
	var result int
	database.DB.Table("users").Count(&result)
	return result
}
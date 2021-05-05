package dto

import (
	database "github.com/lauralunddahl/DevOps-GroupF/api/db"
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

func GetUserById(userId int) User {
	user := User{}
	res := database.DB.Where("user_id = ?", userId).First(&user)
	if res.Error != nil {
		log.Println("GetUserById")
		log.Error(res.Error)
	}
	return user
}

func GetUsername(userId int) string {
	user := User{}
	res := database.DB.Where("user_id = ?", userId).First(&user)
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

func GetTotalNumberOfUsers() int64 {
	var result int64
	database.DB.Table("users").Count(&result)
	return result
}

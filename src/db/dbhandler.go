package database

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

type Message struct {
	//gorm.Model
	MessageId int
	AuthorId  int
	Text      string
	PubDate   int
	Flagged   int
}

func OpenDB() *gorm.DB {
	db, err := gorm.Open("mysql", "groupf:f0psD3v1123@(localhost)/minitwit?charset=utf8&parseTime=True&loc=Local")
	defer db.Close()
	if err != nil {
		fmt.Println(err)
		panic("Failed to connect to the database!")
	}

	db.AutoMigrate(&Message{})

	// var user []User
	// db.Find(&user)
	// fmt.Println("{}", user)

	var message []Message
	db.Find(&message)
	fmt.Println("{}", message)

	return db
}

// func (User) TableName() string {
// 	return "user"
// }

func (Message) TableName() string {
	return "message"
}

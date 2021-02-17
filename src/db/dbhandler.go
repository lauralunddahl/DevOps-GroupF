package db

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var DB *gorm.DB

func init() {
	database, err := gorm.Open("mysql", "groupf:f0psD3v1123@(localhost)/minitwit?charset=utf8&parseTime=True&loc=Local")
	//defer db.Close()
	if err != nil {
		fmt.Println(err)
		panic("Failed to connect to the database!")
	}
	DB = database
}

func GetDB() *gorm.DB {
	return DB
}

// func initialMigration() {
// 	db.AutoMigrate(&User{})
// 	db.AutoMigrate(&Message{})
// 	db.AutoMigrate(&Follower{})
// }

package db

import (
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
)

var DB *gorm.DB

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Not running locally")
	}
	password := os.Getenv("DB_PASSWORD")

	database, err := gorm.Open("mysql", "fibonacci:"+password+"@(mydb.itu.dk)/minitwit?charset=utf8&parseTime=True&loc=Local")
	//database, err := gorm.Open("mysql", "**REMOVED**:**REMOVED**@(mydb.itu.dk)/minitwit_test?charset=utf8&parseTime=True&loc=Local")
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

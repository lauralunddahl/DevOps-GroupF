package db

import (
	"fmt"
	"os"
	"gorm.io/gorm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
)

var DB *gorm.DB

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Not running locally")
	}
	password := os.Getenv("DB_PASSWORD")
	database, err := gorm.Open(mysql.Open("fibonacci:" + password + "@(minitwit-mysql-db-do-user-8729061-0.b.db.ondigitalocean.com:25060)/defaultdb?charset=utf8&parseTime=True&loc=Local"))

	if err != nil {
		fmt.Println(err)
		panic("Failed to connect to the database!")
	}
	DB = database
}

func GetDB() *gorm.DB {
	return DB
}

package main

import (
	"database/sql"
	"fmt"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func connect_db() {
	db, err := sql.Open("sqlite3", "home/parallels/Desktop/DevOps/DevOps-GroupF/minitwit.db")
	checkErr(err)
	fmt.Printf(err.Error())
	db.Close()
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func helloWorld(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World")
}

func main() {
	connect_db()
	http.HandleFunc("/", helloWorld)
	http.ListenAndServe(":8080", nil)
}

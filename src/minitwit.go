package main

import (
	"database/sql"
	"fmt"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB
/*func connect_db() {
	db, err := sql.Open("sqlite3", "../minitwit.db")
	checkErr(err)
    }*/

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func helloWorld(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World")
}

func booksIndex(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "../minitwit.db")
	checkErr(err)
    rows, err := db.Query("Select * from user")
    checkErr(err)

    for rows.Next() {
		var username string
		var user_id string
		var email string
		var pw_hash string
        err := rows.Scan(&user_id,&username,&email,&pw_hash)
		checkErr(err)
		fmt.Fprintln(w, username)
	}
}     


func main() {
	http.HandleFunc("/", booksIndex)
	http.ListenAndServe(":8080", nil)
}

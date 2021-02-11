package main

import (
	"crypto/sha1"
	"database/sql"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const database string = "../minitwit.db"
const per_page int = 30
const debug bool = true
const secret_key string = "development key"

var DB *sql.DB = connect_db()

func connect_db() (DB *sql.DB) {
	db, err := sql.Open("sqlite3", database)
	checkErr(err)
	return db
}

func init_db() {
	//unsure if we are already doing this in connect_db()
}

func query_db(query string, args string, one bool) (*sql.Rows, error) {
	stmt, err := DB.Prepare(query)
	checkErr(err)
	defer stmt.Close()
	rows, err := stmt.Query(args)
	return rows, err
}

func format_datetime(timestamp string) string {
	const layout = "2016-03-28 @ 08:30"
	t, err := time.Parse(layout, timestamp)
	checkErr(err)
	return t.String()
}

func gravatar_url(email string) string {
	size := 80
	h := sha1.New()
	h.Write([]byte(strings.ToLower(strings.TrimSpace(email))))
	sha1_hash := hex.EncodeToString(h.Sum(nil))
	return fmt.Sprintf("http://www.gravatar.com/avatar/%s?d=identicon&s=%d", sha1_hash, size)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func booksIndex(w http.ResponseWriter, r *http.Request) {
	rows, err := query_db("Select * from user where username = ?", "a", false)
	checkErr(err)
	defer rows.Close()
	for rows.Next() {
		var username string
		var user_id string
		var email string
		var pw_hash string
		err := rows.Scan(&user_id, &username, &email, &pw_hash)
		checkErr(err)
		fmt.Fprintln(w, username)
	}
}

func main() {
	http.HandleFunc("/", booksIndex)
	http.ListenAndServe(":8080", nil)

	defer DB.Close()
}

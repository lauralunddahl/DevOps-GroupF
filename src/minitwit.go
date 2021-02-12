package main

import (
	"crypto/sha1"
	"database/sql"
	"encoding/hex"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

const database string = "../minitwit.db"
const per_page int = 30
const debug bool = true
const secret_key string = "development key"

var DB *sql.DB = connect_db()
var router = mux.NewRouter()

type User struct {
	Username string
	UserId   string
	Email    string
	PwHash   string
}

type Message struct {
	MessageId string
	AuthorId  string
	Text 	  string
	PubDate   string
	Flagged   int


}

type Timeline struct {
	Username string
	UserId   int
	Email    string
	PwHash   string
	
	MessageId int
	AuthorId  int
	Text 	  string
	PubDate   int
	Flagged   int
}


 

func connect_db() (DB *sql.DB) {
	db, err := sql.Open("sqlite3", database)
	checkErr(err)
	return db
}

func init_db() {
	//unsure if we are already doing this in connect_db()
}

func query_db(query string, args string, one bool) *sql.Rows {
	stmt, err := DB.Prepare(query)
	checkErr(err)
	defer stmt.Close()
	rows, err := stmt.Query(args)
	checkErr(err)
	return rows
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

func before_request(handler func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		//params :=
		//user_id := params["user_id"]
		rows := query_db("select * from user where user_id = ?", "1", true) //hardcoded user_id right now
		defer rows.Close()
		var user User
		for rows.Next() {
			err := rows.Scan(&user.UserId, &user.Username, &user.Email, &user.PwHash)
			checkErr(err)
		}
		handler(w, r)
	}
}

//Not sure if we need the after_request function

func timeline(w http.ResponseWriter, r *http.Request) {
	println(w, "We got a visitor from: "+r.RemoteAddr)
	//if (user.user_id < 0) //redirect to public_timeline
	/*user := User{
		Username: "Nanna",
		Email:    "nanm@itu.dk",
		UserId:   "42",
		PwHash:   "htjdejoi",
	}*/
	
	rows := query_db("select user.*, message.*  from message, user where message.flagged = 0 and message.author_id = user.user_id order by message.pub_date desc limit ?",string(per_page),false)
	defer rows.Close()
	var timelines []Timeline
	var timeline Timeline
	for rows.Next() {
		err := rows.Scan(&timeline.UserId, &timeline.Username, &timeline.Email, &timeline.PwHash,
		&timeline.MessageId, &timeline.AuthorId, &timeline.Text, &timeline.PubDate, &timeline.Flagged)
		checkErr(err)
		timelines = append(timelines,Timeline{UserId: timeline.UserId, Username: timeline.Username,
		Email: timeline.Email, PwHash: timeline.PwHash, MessageId: timeline.MessageId, AuthorId: timeline.AuthorId,
		Text: timeline.Text, PubDate: timeline.PubDate, Flagged: timeline.Flagged})
	}
	printSlice(timelines)
	templ, _ := template.ParseFiles("../templates/tmp.html")
	//pubTimeline, _ := template.ParseFiles("../templates/timeline.html")

	err := templ.Execute(w, timelines)
	if err != nil {
		fmt.Fprintln(w, err)
	}
}

func printSlice(s []Timeline) {
	fmt.Printf("len=%d cap=%d %v\n", len(s), cap(s), s)
}
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func booksIndex(w http.ResponseWriter, r *http.Request) {
	rows := query_db("Select * from user where username = ?", "a", false)
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
	

	router.HandleFunc("/", before_request(timeline))

	log.Fatal(http.ListenAndServe(":8080", router))

	defer DB.Close()
}

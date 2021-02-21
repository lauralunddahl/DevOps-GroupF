package main

import (
	"crypto/sha1"
	"database/sql"
	"encoding/hex"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
	"tawesoft.co.uk/go/dialog"
	api "api"
)

const database string = "../minitwit.db"
const per_page int = 30
const debug bool = true
const secret_key string = "development key"

var DB *sql.DB = connect_db()
var router = mux.NewRouter()

var (
	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	key   = []byte("super-secret-key")
	store = sessions.NewCookieStore(key)
)

type User struct {
	Username string
	UserId   string
	Email    string
	PwHash   string
}

type Message struct {
	MessageId string
	AuthorId  string
	Text      string
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
	Text      string
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

func query_db(query string, arg string) *sql.Rows {
	stmt, err := DB.Prepare(query)
	checkErr(err)
	defer stmt.Close()
	rows, err := stmt.Query(arg)
	checkErr(err)
	return rows
}

// TODO - fix so it can take a list of args instead
func query_db_multiple2(query string, arg1 string, arg2 string) *sql.Rows {
	stmt, err := DB.Prepare(query)
	checkErr(err)
	defer stmt.Close()
	rows, err := stmt.Query(arg1, arg2)
	checkErr(err)
	return rows
}

func query_db_multiple3(query string, arg1 string, arg2 string, arg3 string) *sql.Rows {
	stmt, err := DB.Prepare(query)
	checkErr(err)
	defer stmt.Close()
	rows, err := stmt.Query(arg1, arg2, arg3)
	checkErr(err)
	return rows
}

func get_user_id(username string) int {
	rows := query_db("SELECT user_id from user where username = ?", username)
	defer rows.Close()

	var uid int

	for rows.Next() {
		err := rows.Scan(&uid)
		checkErr(err)
	}
	return uid
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
		session, _ := store.Get(r, "session1")
		session.Values["authenticated"] = false
		session.Values["userid"] = 0
		err := session.Save(r, w)
		if err != nil {
			println(err.Error())
		}

		handler(w, r)
	}
}

//Not sure if we need the after_request function

func timeline(w http.ResponseWriter, r *http.Request) {
	println(w, "We got a visitor from: "+r.RemoteAddr)

	session, _ := store.Get(r, "session1")

	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Redirect(w, r, "/public", 302)
	} else {
		println("i dont want to be here")
		user_id := session.Values["userid"].(int)
		rows := query_db_multiple3("select user.*, message.* from message, user where message.flagged = 0 and message.author_id = user.user_id and (user.user_id = ? or user.user_id in (select whom_id from follower where who_id = ?)) order by message.pub_date desc limit ?", strconv.Itoa(user_id), strconv.Itoa(user_id), strconv.Itoa(per_page))
		defer rows.Close()
		var timelines []Timeline
		var timeline Timeline
		for rows.Next() {
			err := rows.Scan(&timeline.UserId, &timeline.Username, &timeline.Email, &timeline.PwHash,
				&timeline.MessageId, &timeline.AuthorId, &timeline.Text, &timeline.PubDate, &timeline.Flagged)
			checkErr(err)
			timelines = append(timelines, timeline)
		}

		templ := template.Must(template.ParseFiles("../templates/tmp.html", "../templates/layout.html"))

		err := templ.Execute(w, map[string]interface{}{
			"timeline": timelines,
			"public":   false,
		})
		if err != nil {
			fmt.Fprintln(w, err)
		}
	}
}

//marcus
func loginpage(w http.ResponseWriter, r *http.Request) {
	loginp, err := template.ParseFiles("../templates/login.html")
	if err != nil {
		println(err.Error())
	}
	err = loginp.Execute(w, nil)
	if err != nil {
		println(err.Error())
	}
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	println("handle login")
	username := r.FormValue("username")
	password := r.FormValue("password")
	println(username + " " + password)
	row := query_db("select * from user where username = ? ", username)
	var user User
	for row.Next() {
		err := row.Scan(&user.UserId, &user.Username, &user.Email, &user.PwHash)
		checkErr(err)
	}
	if user.Username == "" {
		fmt.Fprintln(w, "invalid username")
	}
	println(user.Username)
	//check password hash from database against input password from user
	byteHash := []byte(user.PwHash)
	bytePw := []byte(password)
	err := bcrypt.CompareHashAndPassword(byteHash, bytePw)
	if err != nil {
		println(err.Error())
		fmt.Fprintln(w, "invalid password")
	}

	session, _ := store.Get(r, "session1")
	session.Values["authenticated"] = true
	session.Values["userId"] = user.UserId
	session.Save(r, w)
	http.Redirect(w, r, "/", 302)
}

func public_timeline(w http.ResponseWriter, r *http.Request) {
	rows := query_db("select user.*, message.*  from message, user where message.flagged = 0 and message.author_id = user.user_id order by message.pub_date desc limit ?", strconv.Itoa(per_page))
	defer rows.Close()
	var timelines []Timeline
	var timeline Timeline
	for rows.Next() {
		err := rows.Scan(&timeline.UserId, &timeline.Username, &timeline.Email, &timeline.PwHash,
			&timeline.MessageId, &timeline.AuthorId, &timeline.Text, &timeline.PubDate, &timeline.Flagged)
		checkErr(err)
		timelines = append(timelines, timeline)
	}

	templ := template.Must(template.ParseFiles("../templates/tmp.html", "../templates/layout.html"))

	err := templ.Execute(w, map[string]interface{}{
		"timeline": timelines,
		"public":   true,
	})
	if err != nil {
		fmt.Fprintln(w, err)
	}
}

func user_timeline(w http.ResponseWriter, r *http.Request) {
	user_id := 0
	vars := mux.Vars(r)

	username := vars["username"]

	session, _ := store.Get(r, "session1")

	//session.Values["authenticated"] = false
	println(session.Values["authenticated"].(bool))
	if auth, _ := session.Values["authenticated"].(bool); auth {
		user_id = session.Values["userid"].(int)
	}
	profileuserRow := query_db("select * from user where username = ?", username)
	defer profileuserRow.Close()
	var profileuser User
	if profileuserRow.Next() {
		queryErr := profileuserRow.Scan(&profileuser.UserId, &profileuser.Username, &profileuser.Email, &profileuser.PwHash)
		checkErr(queryErr)

		followedRow := query_db_multiple2("select 1 from follower where follower.who_id = ? and follower.whom_id = ?", strconv.Itoa(user_id), profileuser.UserId)
		defer followedRow.Close()
		followed := false
		if followedRow.Next() {
			followed = true
		}

		rows := query_db_multiple2("select user.*, message.* from message, user where user.user_id = message.author_id and user.user_id = ? order by message.pub_date desc limit ?", profileuser.UserId, strconv.Itoa(per_page))

		defer rows.Close()
		var timelines []Timeline
		var timeline Timeline
		for rows.Next() {
			err := rows.Scan(&timeline.UserId, &timeline.Username, &timeline.Email, &timeline.PwHash,
				&timeline.MessageId, &timeline.AuthorId, &timeline.Text, &timeline.PubDate, &timeline.Flagged)
			checkErr(err)
			timelines = append(timelines, timeline)
		}

		templ := template.Must(template.ParseFiles("../templates/tmp.html", "../templates/layout.html"))

		err := templ.Execute(w, map[string]interface{}{
			"timeline":    timelines,
			"public":      false,
			"profileuser": profileuser,
			"followed":    followed,
		})
		if err != nil {
			fmt.Fprintln(w, err)
		}

	} else {
		http.NotFound(w, r)
	}
}

func follow_user(w http.ResponseWriter, r *http.Request) {
	user_id := 0
	vars := mux.Vars(r)

	username := vars["username"]

	session, _ := store.Get(r, "session1")
	if auth, _ := session.Values["authenticated"].(bool); auth {
		user_id = session.Values["userid"].(int)
	}
	if user_id == 0 {
		http.Error(w, "not authorized", 401)
	} else {
		whom_id := get_user_id(username)
		if whom_id == 0 {
			http.NotFound(w, r)
		}
		stmt, _ := DB.Prepare("insert into follower (who_id, whom_id) values (?,?)")
		_, err := stmt.Exec(user_id, whom_id)
		checkErr(err)
		dialog.Alert("You are now following %s", username)
		http.Redirect(w, r, "/{username}", 302)
	}
}

func unfollow_user(w http.ResponseWriter, r *http.Request) {
	user_id := 0
	vars := mux.Vars(r)

	username := vars["username"]

	session, _ := store.Get(r, "session1")
	if auth, _ := session.Values["authenticated"].(bool); auth {
		user_id = session.Values["userid"].(int)
	}
	if user_id == 0 {
		http.Error(w, "not authorized", 401)
	} else {
		whom_id := get_user_id(username)
		if whom_id == 0 {
			http.NotFound(w, r)
		}
		stmt, _ := DB.Prepare("delete from follower where who_id = ? and whom_id = ?")
		_, err := stmt.Exec(user_id, whom_id)
		checkErr(err)
		dialog.Alert("You are no longer following %s", username)
		http.Redirect(w, r, "/{username}", 302)
	}
}

func add_message(w http.ResponseWriter, r *http.Request) {
	user_id := 0
	session, _ := store.Get(r, "session1")
	if auth, _ := session.Values["authenticated"].(bool); auth {
		user_id = session.Values["userid"].(int)
	}
	if user_id == 0 {
		http.Error(w, "not authorized", 401)
	} else {
		text := r.FormValue("text")
		stmt, _ := DB.Prepare("insert into message (author_id, text, pub_date, flagged) values (?, ?, ?, 0)")
		_, err := stmt.Exec(user_id, text, time.Now())
		checkErr(err)
		dialog.Alert("Your message was recorded")
		http.Redirect(w, r, "/", 302)
	}
}

//Nanna
func register(w http.ResponseWriter, r *http.Request) {
	register, err := template.ParseFiles("../templates/register.html")
	if err != nil {
		println(err.Error())
	}
	session, _ := store.Get(r, "session1")
	if auth, _ := session.Values["authenticated"].(bool); auth {
		http.Redirect(w, r, "/", 302)
	} else {
		err = register.Execute(w, nil)
		if err != nil {
			println(err.Error())
		}
	}
}

func handleRegister(w http.ResponseWriter, r *http.Request) {
	println("handle register")
	err := ""
	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")
	password2 := r.FormValue("password2")
	if len(username) == 0 {
		err = "You have to enter a username\n"
	} else if len(email) == 0 || strings.Contains(email, "@") {
		err += "You have to enter a vlid email address\n"
	} else if len(password) == 0 {
		err += "You have to enter a password\n"
	} else if password != password2 {
		err += "The two passwords do not match\n"
	} else if get_user_id(username) > 0 { //this might have to be another check at some point
		err += "The username is already taken"
	} else {
		pw_hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
		if err != nil {
			println(err.Error())
		}
		query_db_multiple3("insert into user (username, email, pw_hash) values (?, ?, ?)", username, email, string(pw_hash))
		fmt.Println(w, "You were successfully registered and can login now")
		http.Redirect(w, r, "/login", 302)
	}
}

//Louise
func logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session1") //What should this be called?
	session.Values["authenticated"] = false
	session.Values["userId"] = ""
	session.Save(r, w)
	http.Redirect(w, r, "/", 302)
}

func printSlice(s []Timeline) {
	fmt.Printf("len=%d cap=%d %v\n", len(s), cap(s), s)
}
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {

	router.HandleFunc("/", before_request(timeline))
	router.HandleFunc("/register", before_request(register))
	router.HandleFunc("/registerfunc", handleRegister).Methods("POST")
	router.HandleFunc("/", before_request(timeline))
	router.HandleFunc("/login", before_request(loginpage))
	router.HandleFunc("/loginfunc", handleLogin).Methods("POST")
	router.HandleFunc("/public", public_timeline)
	api.HandleApiRequest(router)
	router.HandleFunc("/{username}", user_timeline)

	router.HandleFunc("/{username}/follow", follow_user)
	router.HandleFunc("/{username}/unfollow", unfollow_user)
	router.HandleFunc("/add_message", add_message).Methods("POST")
	router.HandleFunc("/logout", logout)

	log.Fatal(http.ListenAndServe(":8080", router))

	defer DB.Close()
}

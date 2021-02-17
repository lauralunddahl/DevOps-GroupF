package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

type Latest struct {
	La string `json:"latest"`
}

type Response struct {
	Status   int    `json:"status"`
	ErrorMsg string `json:"error_msg"`
}

type ApiMessage struct {
	Content string `json:"content"`
	PubDate string `json:"pub_date"`
	User    string `json:"user"`
}

type Register struct {
	Username  string `json:"username"`
	Email     string `json:"email" `
	Password  string `json:"password"`
	Password2 string `json:"password2"`
}

var latest = 0

func not_req_from_simulator(w http.ResponseWriter, r *http.Request) {
	fromSim := r.Header.Get("Authorization")
	w.Header().Set("Content-Type", "application/json")
	if fromSim != "Basic c2ltdWxhdG9yOnN1cGVyX3NhZmUh" {
		err := "You are not authorized to use this resource!"
		json.NewEncoder(w).Encode(err)
	}

	fmt.Println("Endpoint Hit: homePage")
}

func update_latest(w http.ResponseWriter, r *http.Request) int {
	late := r.URL.Query().Get("latest")
	if late == "" {
		late = "-1"
	}
	x, err := strconv.Atoi(late)
	if err != nil {
		println(err.Error())
	}
	if x != -1 {
		latest = x
		return latest
	} else {
		return latest
	}
}

//Get

func get_latest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var ls Latest
	ls.La = strconv.Itoa(latest)
	json.NewEncoder(w).Encode(ls)

}
func apiRegister(w http.ResponseWriter, r *http.Request) {
	update_latest(w, r)
	w.Header().Set("Content-Type", "application/json")
	err := ""
	var newReg Register
	json.NewDecoder(r.Body).Decode(&newReg)
	println(newReg.Username)
	println(newReg.Email)
	println(newReg.Password)
	println(newReg.Password2)

	if len(newReg.Username) == 0 {
		err = "You have to enter a username\n"
	} else if len(newReg.Email) == 0 { //todo: check for @ code: || strings.Contains(newReg.Email, "@")
		err += "You have to enter a valid email address\n"
	} else if len(newReg.Password) == 0 {
		err += "You have to enter a password\n"
	} else if newReg.Password != newReg.Password2 {
		err += "The two passwords do not match\n"
	} else if get_user_id(newReg.Username) > 0 { //this might have to be another check at some point
		err += "The username is already taken"
	} else {
		pw_hash, err := bcrypt.GenerateFromPassword([]byte(newReg.Password), bcrypt.MinCost)
		if err != nil {
			println(err.Error())
		}
		query_db_multiple3("insert into user (username, email, pw_hash) values (?, ?, ?)", newReg.Username, newReg.Email, string(pw_hash))
		fmt.Println(w, "You were successfully registered and can login now")
	}
	if err != "" {
		var res Response
		res.Status = 400
		res.ErrorMsg = err
		json.NewEncoder(w).Encode(res)
		return
	} else {
		var res Response
		res.Status = 204
		res.ErrorMsg = ""
		json.NewEncoder(w).Encode(res)
		return
	}
}

//GET
func messages(w http.ResponseWriter, r *http.Request) {
	update_latest(w, r)

	//not_req_from_simulator(w, r)

	noMsg := r.URL.Query().Get("no")
	if noMsg == "" {
		noMsg = "100"
	}
	rows := query_db("select user.*, message.*  from message, user where message.flagged = 0 and message.author_id = user.user_id order by message.pub_date desc limit ?", noMsg)
	var timelines []Timeline
	var timeline Timeline
	for rows.Next() {
		err := rows.Scan(&timeline.UserId, &timeline.Username, &timeline.Email, &timeline.PwHash,
			&timeline.MessageId, &timeline.AuthorId, &timeline.Text, &timeline.PubDate, &timeline.Flagged)
		checkErr(err)
		timelines = append(timelines, timeline)
	}
	var messages []ApiMessage
	for _, t := range timelines {
		var message ApiMessage
		message.Content = t.Text
		message.PubDate = strconv.Itoa(t.PubDate)
		message.User = t.Username
		messages = append(messages, message)
	}
	json.NewEncoder(w).Encode(messages)
}

//GET, POST
func messages_per_user(w http.ResponseWriter, r *http.Request) {
	update_latest(w, r)

	vars := mux.Vars(r)

	username := vars["username"]

	//not_req_from_simulator(w,r)

	noMsg := r.URL.Query().Get("no")
	if noMsg == "" {
		noMsg = "100"
	}

	switch r.Method {
	case "GET":
		userId := get_user_id(username)
		if userId == 0 {
			var res Response
			res.Status = 400
			res.ErrorMsg = "No user found"
			json.NewEncoder(w).Encode(res)
		} else {
			rows := query_db_multiple2("select user.*, message.* from message, user where user.user_id = message.author_id and user.user_id = ? order by message.pub_date desc limit ?", strconv.Itoa(userId), noMsg)
			var timelines []Timeline
			var timeline Timeline
			for rows.Next() {
				err := rows.Scan(&timeline.UserId, &timeline.Username, &timeline.Email, &timeline.PwHash,
					&timeline.MessageId, &timeline.AuthorId, &timeline.Text, &timeline.PubDate, &timeline.Flagged)
				checkErr(err)
				timelines = append(timelines, timeline)
			}
			var messages []ApiMessage
			for _, t := range timelines {
				var message ApiMessage
				message.Content = t.Text
				message.PubDate = strconv.Itoa(t.PubDate)
				message.User = t.Username
				messages = append(messages, message)
			}
			json.NewEncoder(w).Encode(messages)
		}
	case "POST":
		w.Header().Set("Content-Type", "application/json")
		var message Message
		json.NewDecoder(r.Body).Decode(&message)
		stmt, _ := DB.Prepare("insert into message (author_id, text, pub_date, flagged) values (?, ?, ?, 0)")
		_, err := stmt.Exec(get_user_id(username), message.Text, time.Now().Format("1612171952"))
		if err != nil {
			println(err.Error())
		}
		var res Response
		res.Status = 204
		res.ErrorMsg = ""
		json.NewEncoder(w).Encode(res)

	}

}

func handleApiRequest(router *mux.Router) {
	router.HandleFunc("/api/test", not_req_from_simulator)
	router.HandleFunc("/latest", get_latest).Methods("GET")
	router.HandleFunc("/api/register", apiRegister).Methods("POST")
	router.HandleFunc("/msgs", messages).Methods("GET")
	router.HandleFunc("/msgs/{username}", messages_per_user).Methods("GET", "POST")

}

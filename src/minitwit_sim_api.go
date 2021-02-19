package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
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

type Followers struct {
	Username string `json:"username"`
}
type FollowUser struct {
	Follow   string `json:"follow"`
	Unfollow string `json:"unfollow"`
}

var latest = 0

func not_req_from_simulator(w http.ResponseWriter, r *http.Request) {
	fromSim := r.Header.Get("Authorization")
	w.Header().Set("Content-Type", "application/json")
	if fromSim != "Basic c2ltdWxhdG9yOnN1cGVyX3NhZmUh" {
		err := "You are not authorized to use this resource!"
		json.NewEncoder(w).Encode(err)
	}
}

func update_latest(w http.ResponseWriter, r *http.Request) int {
	late, _ := strconv.Atoi(r.URL.Query().Get("latest"))
	if late != 0 {
		latest = late
	}
	return latest
}

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

	if len(newReg.Username) == 0 {
		err = "You have to enter a username\n"
	} else if len(newReg.Email) == 0 || !strings.Contains(newReg.Email, "@") {
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

func messages(w http.ResponseWriter, r *http.Request) {
	update_latest(w, r)

	//not_req_from_simulator(w, r)

	no_msg := r.URL.Query().Get("no")
	if no_msg == "" {
		no_msg = "100"
	}
	rows := query_db("select user.*, message.*  from message, user where message.flagged = 0 and message.author_id = user.user_id order by message.pub_date desc limit ?", no_msg)
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

func messages_per_user(w http.ResponseWriter, r *http.Request) {
	update_latest(w, r)

	vars := mux.Vars(r)
	username := vars["username"]

	//not_req_from_simulator(w, r)

	no_msg := r.URL.Query().Get("no")
	if no_msg == "" {
		no_msg = "100"
	}

	switch r.Method {
	case "GET":
		user_id := get_user_id(username)
		if user_id == 0 {
			var res Response
			res.Status = 500
			res.ErrorMsg = "No user found"
			json.NewEncoder(w).Encode(res)
		} else {
			rows := query_db_multiple2("select user.*, message.* from message, user where user.user_id = message.author_id and user.user_id = ? order by message.pub_date desc limit ?", strconv.Itoa(user_id), no_msg)
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

func follow(w http.ResponseWriter, r *http.Request) {
	update_latest(w, r)

	//not_req_from_simulator(w, r)
	vars := mux.Vars(r)
	username := vars["username"]

	no_followers := r.URL.Query().Get("no")
	if no_followers == "" {
		no_followers = "100"
	}

	user_id := get_user_id(username)
	if user_id == 0 {
		var res Response
		res.Status = 404
		res.ErrorMsg = "No user found"
		json.NewEncoder(w).Encode(res)
	}

	switch r.Method {

	case "POST":
		w.Header().Set("Content-Type", "application/json")
		var follows FollowUser
		json.NewDecoder(r.Body).Decode(&follows)
		if len(follows.Follow) > 0 {
			follows_username := follows.Follow
			follows_user_id := get_user_id(follows_username)
			if follows_user_id == 0 {
				var res Response
				res.Status = 500
				res.ErrorMsg = "No user found"
				json.NewEncoder(w).Encode(res)
			} else {
				stmt, _ := DB.Prepare("insert into follower (who_id, whom_id) values (?, ?)")
				_, err := stmt.Exec(user_id, follows_user_id)
				if err != nil {
					println(err.Error())
				}
				var res Response
				res.Status = 204
				res.ErrorMsg = ""
				json.NewEncoder(w).Encode(res)
			}
		} else if len(follows.Unfollow) > 0 {
			unfollows_username := follows.Unfollow
			unfollows_user_id := get_user_id(unfollows_username)
			if user_id == 0 {
				var res Response
				res.Status = 500
				res.ErrorMsg = "No user found"
				json.NewEncoder(w).Encode(res)
			} else {
				stmt, _ := DB.Prepare("delete from follower where who_id=? and whom_id=?")
				_, err := stmt.Exec(user_id, unfollows_user_id)
				if err != nil {
					println(err.Error())
				}
				var res Response
				res.Status = 204
				res.ErrorMsg = ""
				json.NewEncoder(w).Encode(res)
			}
		}
	case "GET":
		rows := query_db_multiple2("select user.username from user inner join follower on follower.whom_id=user.user_id where follower.who_id = ? limit ?", strconv.Itoa(user_id), no_followers)
		defer rows.Close()
		var followers []Followers
		var follower Followers
		for rows.Next() {
			err := rows.Scan(&follower.Username)
			if err != nil {
				println(err.Error())
			}
			followers = append(followers, follower)
		}
		json.NewEncoder(w).Encode(followers)
	}
}

func handleApiRequest(router *mux.Router) {
	router.HandleFunc("/api/test", not_req_from_simulator)
	router.HandleFunc("/latest", get_latest).Methods("GET")
	router.HandleFunc("/api/register", apiRegister).Methods("POST")
	router.HandleFunc("/msgs", messages).Methods("GET")
	router.HandleFunc("/fllws/{username}", follow).Methods("GET", "POST")
	router.HandleFunc("/msgs/{username}", messages_per_user).Methods("GET", "POST")
}

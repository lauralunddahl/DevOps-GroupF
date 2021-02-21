package api

import (
	"dto"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

type Timeline struct {
	Username string
	UserId   int
	Email    string
	PwHash   string

	MessageId int
	AuthorId  int
	Text      string
	PubDate   time.Time
	Flagged   int
}

type Message struct {
	MessageId string
	AuthorId  string
	Text      string
	PubDate   time.Time
	Flagged   int
}

type Latest struct {
	La string `json:"latest"`
}

type Response struct {
	Status   int    `json:"status"`
	ErrorMsg string `json:"error_msg"`
}

type ApiMessage struct {
	Content string    `json:"content"`
	PubDate time.Time `json:"pub_date"`
	User    string    `json:"user"`
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

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

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
	} else if dto.GetUserID(newReg.Username) > 0 { //this might have to be another check at some point
		err += "The username is already taken"
	} else {
		pw_hash, err := bcrypt.GenerateFromPassword([]byte(newReg.Password), bcrypt.MinCost)
		if err != nil {
			println(err.Error())
		}
		dto.RegisterUser(newReg.Username, newReg.Email, string(pw_hash))
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
	var timelines = dto.GetPublicTimeline() //update to no_msg
	var messages []ApiMessage
	for _, t := range timelines {
		var message ApiMessage
		message.Content = t.Text
		message.PubDate = t.PubDate
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
		user_id := dto.GetUserID(username)
		if user_id == 0 {
			var res Response
			res.Status = 500
			res.ErrorMsg = "No user found"
			json.NewEncoder(w).Encode(res)
		} else {
			var timelines = dto.GetUserTimeline(user_id) //update to no_msg
			var messages []ApiMessage
			for _, t := range timelines {
				var message ApiMessage
				message.Content = t.Text
				message.PubDate = t.PubDate
				message.User = t.Username
				messages = append(messages, message)
			}
			json.NewEncoder(w).Encode(messages)
		}
	case "POST":
		user_id := dto.GetUserID(username)
		w.Header().Set("Content-Type", "application/json")
		var message Message
		json.NewDecoder(r.Body).Decode(&message)
		dto.AddMessage(strconv.Itoa(user_id), message.Text, time.Now(), 0)
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

	user_id := dto.GetUserID(username)
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
			follows_user_id := dto.GetUserID(follows_username)
			if follows_user_id == 0 {
				var res Response
				res.Status = 500
				res.ErrorMsg = "No user found"
				json.NewEncoder(w).Encode(res)
			} else {
				dto.FollowUser(user_id, follows_user_id)
				var res Response
				res.Status = 204
				res.ErrorMsg = ""
				json.NewEncoder(w).Encode(res)
			}
		} else if len(follows.Unfollow) > 0 {
			unfollows_username := follows.Unfollow
			unfollows_user_id := dto.GetUserID(unfollows_username)
			if user_id == 0 {
				var res Response
				res.Status = 500
				res.ErrorMsg = "No user found"
				json.NewEncoder(w).Encode(res)
			} else {
				dto.UnfollowUser(user_id, unfollows_user_id)
				var res Response
				res.Status = 204
				res.ErrorMsg = ""
				json.NewEncoder(w).Encode(res)
			}
		}
	case "GET":
		numb, _ := strconv.Atoi(no_followers)
		var followers = dto.GetFollowers(user_id, numb)
		json.NewEncoder(w).Encode(followers)
	}
}

func HandleApiRequest(router *mux.Router) {
	router.HandleFunc("/api/test", not_req_from_simulator)
	router.HandleFunc("/latest", get_latest).Methods("GET")
	router.HandleFunc("/api/register", apiRegister).Methods("POST")
	router.HandleFunc("/msgs", messages).Methods("GET")
	router.HandleFunc("/fllws/{username}", follow).Methods("GET", "POST")
	router.HandleFunc("/msgs/{username}", messages_per_user).Methods("GET", "POST")
}

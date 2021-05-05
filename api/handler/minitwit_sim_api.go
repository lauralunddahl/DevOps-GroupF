package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	dto "github.com/lauralunddahl/DevOps-GroupF/api/dto"
	metrics "github.com/lauralunddahl/DevOps-GroupF/api/metrics"
	helper "github.com/lauralunddahl/DevOps-GroupF/api/helper"
	log "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

var latest = 0
var noUserFound = "No user found"

func updateLatest(w http.ResponseWriter, r *http.Request) int {
	late, _ := strconv.Atoi(r.URL.Query().Get("latest"))
	if late != 0 {
		latest = late
	}
	return latest
}

func GetLatest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var ls Latest
	ls.Latest = latest
	metrics.IncrementRequests()
	json.NewEncoder(w).Encode(ls)
}

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	updateLatest(w, r)
	metrics.IncrementRequests()
	err := ""
	var newReg Register
	json.NewDecoder(r.Body).Decode(&newReg)

	if len(newReg.Username) == 0 {
		err = "You have to enter a username\n"
	} else if len(newReg.Email) == 0 || !strings.Contains(newReg.Email, "@") {
		err += "You have to enter a valid email address\n"
	} else if len(newReg.Password) == 0 {
		err += "You have to enter a password\n"
	} else if dto.GetUserID(newReg.Username) > 0 { //this might have to be another check at some point
		err += "The username is already taken"
	}
	var res Response
	if err != "" {
		println(err)
		res.Status = 400
		res.ErrorMsg = err
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		log.Info(http.StatusText(http.StatusBadRequest))
	} else {
		pwHash, err := bcrypt.GenerateFromPassword([]byte(newReg.Password), bcrypt.MinCost)
		if err != nil {
			println(err.Error())
		} else {
			image := helper.GravatarUrl(newReg.Email)
			dto.RegisterUser(newReg.Username, newReg.Email, string(pwHash), image)
			res.Status = 204
			res.ErrorMsg = ""
			w.WriteHeader(http.StatusNoContent)
		}
	}
	duration := time.Since(start)
	route := r.URL.Path
	metrics.ObserveResponseTime(route, r.Method, duration.Seconds())
	json.NewEncoder(w).Encode(res)
}

func Messages(w http.ResponseWriter, r *http.Request) {
	updateLatest(w, r)
	metrics.IncrementRequests()

	noMsg := r.URL.Query().Get("no")
	if noMsg == "" {
		noMsg = "100"
	}
	var timelines = dto.GetPublicTimeline() //update to noMsg
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

func MessagesPerUser(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	updateLatest(w, r)
	metrics.IncrementRequests()
	vars := mux.Vars(r)
	username := vars["username"]

	noMsg := r.URL.Query().Get("no")
	if noMsg == "" {
		noMsg = "100"
	}

	switch r.Method {
	case "GET":
		userId := dto.GetUserID(username)
		if userId == 0 {
			var res Response
			res.Status = 404
			res.ErrorMsg = "No user found for " + username
			json.NewEncoder(w).Encode(res)
			log.Info("User id for user " + username + " was not found")
		} else {
			var timelines = dto.GetUserTimeline(userId) //update to noMsg
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
		userId := dto.GetUserID(username)
		var message ApiMessage
		json.NewDecoder(r.Body).Decode(&message)
		dto.AddMessage(strconv.Itoa(userId), message.Content, time.Now(), 0)
		var res Response
		res.Status = 204
		res.ErrorMsg = ""
		w.WriteHeader(http.StatusNoContent)
		json.NewEncoder(w).Encode(res)
	}
	duration := time.Since(start)
	route := r.URL.Path
	metrics.ObserveResponseTime(route, r.Method, duration.Seconds())
}

func Follow(w http.ResponseWriter, r *http.Request) {
	updateLatest(w, r)
	metrics.IncrementRequests()

	vars := mux.Vars(r)
	username := vars["username"]

	noFollowers := r.URL.Query().Get("no")
	if noFollowers == "" {
		noFollowers = "100"
	}

	userId := dto.GetUserID(username)
	if userId == 0 {
		var res Response
		res.Status = 404
		res.ErrorMsg = noUserFound
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		json.NewEncoder(w).Encode(res)
		log.Info("User id for user " + username + " was not found")
	}

	switch r.Method {
	case "POST":
		w.Header().Set("Content-Type", "application/json")
		var follows FollowUser
		json.NewDecoder(r.Body).Decode(&follows)
		if len(follows.Follow) > 0 {
			followsUsername := follows.Follow
			followsUserId := dto.GetUserID(followsUsername)
			if followsUserId == 0 {
				var res Response
				res.Status = 404
				res.ErrorMsg = noUserFound
				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
				json.NewEncoder(w).Encode(res)
				log.Info("User id for user to follow " + followsUsername + " not found")
			} else {
				dto.FollowUser(userId, followsUserId)
				metrics.IncrementFollows()
				var res Response
				res.Status = 204
				res.ErrorMsg = ""
				w.WriteHeader(http.StatusNoContent)
				json.NewEncoder(w).Encode(res)
			}
		} else if len(follows.Unfollow) > 0 {
			unfollowsUsername := follows.Unfollow
			unfollowsUserId := dto.GetUserID(unfollowsUsername)
			if userId == 0 {
				var res Response
				res.Status = 404
				res.ErrorMsg = noUserFound
				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
				json.NewEncoder(w).Encode(res)
				log.Info("User id for user to unfollow " + unfollowsUsername + " not found")
			} else {
				dto.UnfollowUser(userId, unfollowsUserId)
				metrics.IncrementUnfollows()
				var res Response
				res.Status = 204
				res.ErrorMsg = ""
				w.WriteHeader(http.StatusNoContent)
				json.NewEncoder(w).Encode(res)
			}
		}
	case "GET":
		numb, _ := strconv.Atoi(noFollowers)
		var followers = dto.GetFollowers(userId, numb)
		json.NewEncoder(w).Encode(followers)
	}
}

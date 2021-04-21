package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	minitwit "github.com/lauralunddahl/DevOps-GroupF/program/minitwit"
)

func main() {
	router := mux.NewRouter()

	router.PathPrefix("/css/").Handler(http.StripPrefix("/css/", http.FileServer(http.Dir("./css"))))

	router.HandleFunc("/", minitwit.Private_timeline).Methods("GET")
	router.HandleFunc("/register", minitwit.Before_request(minitwit.Register)).Methods("GET")
	router.HandleFunc("/registerfunc", minitwit.HandleRegister).Methods("POST")
	router.HandleFunc("/", minitwit.Before_request(minitwit.Private_timeline)).Methods("GET")
	router.HandleFunc("/login", minitwit.Before_request(minitwit.Loginpage))
	router.HandleFunc("/loginfunc", minitwit.HandleLogin).Methods("POST")
	router.HandleFunc("/public", minitwit.Public_timeline)
	router.HandleFunc("/add_message", minitwit.Add_message).Methods("POST")
	router.HandleFunc("/logout", minitwit.Logout)

	router.HandleFunc("/{username}", minitwit.User_timeline).Methods("GET")
	router.HandleFunc("/{username}/follow", minitwit.Follow_user)
	router.HandleFunc("/{username}/unfollow", minitwit.Unfollow_user)

	log.Fatal(http.ListenAndServe(":8080", router))
}

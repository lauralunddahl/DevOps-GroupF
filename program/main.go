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

	router.HandleFunc("/", minitwit.PrivateTimeline).Methods("GET")
	router.HandleFunc("/register", minitwit.BeforeRequest(minitwit.Register)).Methods("GET")
	router.HandleFunc("/registerfunc", minitwit.HandleRegister).Methods("POST")
	router.HandleFunc("/", minitwit.BeforeRequest(minitwit.PrivateTimeline)).Methods("GET")
	router.HandleFunc("/login", minitwit.BeforeRequest(minitwit.Loginpage))
	router.HandleFunc("/loginfunc", minitwit.HandleLogin).Methods("POST")
	router.HandleFunc("/public", minitwit.PublicTimeline)
	router.HandleFunc("/add_message", minitwit.AddMessage).Methods("POST")
	router.HandleFunc("/logout", minitwit.Logout)

	router.HandleFunc("/{username}", minitwit.UserTimeline).Methods("GET")
	router.HandleFunc("/{username}/follow", minitwit.FollowUser)
	router.HandleFunc("/{username}/unfollow", minitwit.UnfollowUser)

	log.Fatal(http.ListenAndServe(":8080", router))
}

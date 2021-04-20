package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	logging "github.com/lauralunddahl/DevOps-GroupF/src/program/logging"
	metrics "github.com/lauralunddahl/DevOps-GroupF/src/program/metrics"
	minitwit "github.com/lauralunddahl/DevOps-GroupF/src/program/minitwit"
	promhttp "github.com/prometheus/client_golang/prometheus/promhttp"
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

	metrics.RecordMetrics()
	router.Handle("/metrics", promhttp.Handler())

	router.HandleFunc("/{username}", minitwit.User_timeline).Methods("GET")
	router.HandleFunc("/{username}/follow", minitwit.Follow_user)
	router.HandleFunc("/{username}/unfollow", minitwit.Unfollow_user)
	logging.Logging()

	log.Fatal(http.ListenAndServe(":8080", router))
}

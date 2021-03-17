package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	api "github.com/lauralunddahl/DevOps-GroupF/src/api"
	metrics "github.com/lauralunddahl/DevOps-GroupF/src/metrics"
	minitwit "github.com/lauralunddahl/DevOps-GroupF/src/minitwit"
	promhttp "github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	router := mux.NewRouter()
	apirouter := mux.NewRouter()

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

	apirouter.HandleFunc("/latest", api.Get_latest).Methods("GET")
	apirouter.HandleFunc("/register", api.ApiRegister).Methods("POST")
	apirouter.HandleFunc("/msgs", api.Messages).Methods("GET")
	apirouter.HandleFunc("/fllws/{username}", api.Follow).Methods("GET", "POST")
	apirouter.HandleFunc("/msgs/{username}", api.Messages_per_user).Methods("GET", "POST")

	metrics.RecordMetrics()
	router.Handle("/metrics", promhttp.Handler())

	router.HandleFunc("/{username}", minitwit.User_timeline).Methods("GET")
	router.HandleFunc("/{username}/follow", minitwit.Follow_user)
	router.HandleFunc("/{username}/unfollow", minitwit.Unfollow_user)

	//go func() { log.Fatal(http.ListenAndServe(":9092", nil)) }()

	go func() { log.Fatal(http.ListenAndServe(":8081", apirouter)) }()

	log.Fatal(http.ListenAndServe(":8080", router))

}

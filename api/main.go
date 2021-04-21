package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	api "github.com/lauralunddahl/DevOps-GroupF/api/handler"
	logging "github.com/lauralunddahl/DevOps-GroupF/api/logging"
	metrics "github.com/lauralunddahl/DevOps-GroupF/api/metrics"
)

func main() {
	apirouter := mux.NewRouter()

	apirouter.HandleFunc("/latest", api.Get_latest).Methods("GET")
	apirouter.HandleFunc("/register", api.RegisterUser).Methods("POST")
	apirouter.HandleFunc("/msgs", api.Messages).Methods("GET")
	apirouter.HandleFunc("/fllws/{username}", api.Follow).Methods("GET", "POST")
	apirouter.HandleFunc("/msgs/{username}", api.Messages_per_user).Methods("GET", "POST")

	metrics.RecordMetrics()
	logging.Logging()

	log.Fatal(http.ListenAndServe(":8081", apirouter))
}

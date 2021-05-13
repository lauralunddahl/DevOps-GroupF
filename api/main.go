package main

import (
	"log"
	"net/http"
	"github.com/gorilla/mux"
	api "github.com/lauralunddahl/DevOps-GroupF/api/handler"
	metrics "github.com/lauralunddahl/DevOps-GroupF/api/metrics"
	promhttp "github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	apirouter := mux.NewRouter()

	apirouter.HandleFunc("/latest", api.GetLatest).Methods("GET")
	apirouter.HandleFunc("/register", api.RegisterUser).Methods("POST")
	apirouter.HandleFunc("/msgs", api.Messages).Methods("GET")
	apirouter.HandleFunc("/fllws/{username}", api.Follow).Methods("GET", "POST")
	apirouter.HandleFunc("/msgs/{username}", api.MessagesPerUser).Methods("GET", "POST")

	metrics.RecordMetrics()
	apirouter.Handle("/metrics", promhttp.Handler())

	log.Fatal(http.ListenAndServe(":8081", apirouter))
}

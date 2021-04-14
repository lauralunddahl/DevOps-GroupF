package metrics

import (
	//"fmt"
	"strings"
	"time"

	dto "github.com/lauralunddahl/DevOps-GroupF/src/dto"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	cpu "github.com/shirou/gopsutil/cpu"
	mem "github.com/shirou/gopsutil/mem"
)

const interval = 5

var bytes_to_gigabytes = float64(1073741824)
var number_of_https_requests = 0

//Prometheus objects needed as global variables:
var followed = promauto.NewCounter(getCounterOpts("followed_users", "Number of follows"))
var unfollowed = promauto.NewCounter(getCounterOpts("unfollowed_users", "Number of unfollows"))
var httpRequests = promauto.NewCounter(getCounterOpts("http_requests", "Number of http requests"))
var responseTimeRegister = prometheus.NewHistogramVec(getHistogramOpts("http_register_request_duration_seconds", "Histogram of response time for registering a user in seconds"), []string{"route", "method"})
var responseTimeSendMessage = prometheus.NewHistogramVec(getHistogramOpts("http_send_message_request_duration_seconds", "Histogram of response time for sending a message in seconds"), []string{"route", "method"})
var responseTimeRetrieveMessage = prometheus.NewHistogramVec(getHistogramOpts("http_retrieve_message_request_duration_seconds", "Histogram of response time for retrieving a message in seconds"), []string{"route", "method"})

func RecordMetrics() {
	databaseMetrics()
	cpuMetric()
	virtualMemoryMetrics()
	responseTimeMetrics()
}

func databaseMetrics() {
	users := prometheus.NewGauge(getGaugeOpts("users_registered", "Total number of users registered"))
	prometheus.MustRegister(users)

	averageFollowers := prometheus.NewGauge(getGaugeOpts("average_followers_per_user", "Number of average followers per user"))
	prometheus.MustRegister(averageFollowers)

	averagePosts := prometheus.NewGauge(getGaugeOpts("average_posts_per_user", "Number of average posts per user"))
	prometheus.MustRegister(averagePosts)

	go func() {
		for {
			var numberOfUsers = float64(dto.GetTotalNumberOfUsers())
			var numberOfFollowers = float64(dto.GetTotalNumberOfFollowerEntries())
			var numberOfPosts = float64(dto.GetTotalNumberOfMessages())

			users.Set(numberOfUsers)
			averageFollowers.Set(numberOfFollowers / numberOfUsers)
			averagePosts.Set(numberOfPosts / numberOfUsers)
			//fmt.Println("Number of users and average of followers and posts")
			time.Sleep(interval * time.Second)
		}
	}()
}

func cpuMetric() {
	cpuPercentage := prometheus.NewGauge(getGaugeOpts("cpu_total_percentage", "Displays the proportion of the cpu being used"))
	prometheus.MustRegister(cpuPercentage)
	go func() {
		for {
			c, _ := cpu.Percent(0, false)
			cpuPercentage.Set(c[0])
			//fmt.Println("CPU")
			time.Sleep(interval * time.Second)
		}
	}()
}

func virtualMemoryMetrics() {
	memoryPercentage := prometheus.NewGauge(getGaugeOpts("virtual_memory_percentage", "Displays the proportion of the virtual memory being used"))
	prometheus.MustRegister(memoryPercentage)

	memoryAvailable := prometheus.NewGauge(getGaugeOpts("virtual_memory_available", "Information about the RAM available for programs to allocate in gigabytes"))
	prometheus.MustRegister(memoryAvailable)
	go func() {
		for {
			v, _ := mem.VirtualMemory()
			memoryAvailable.Set(float64(v.Available) / bytes_to_gigabytes)
			memoryPercentage.Set(v.UsedPercent)
			//fmt.Println("Virtual Memory")
			time.Sleep(interval * time.Second)
		}
	}()
}

func responseTimeMetrics() {
	prometheus.MustRegister(responseTimeRegister)
	prometheus.MustRegister(responseTimeSendMessage)
	prometheus.MustRegister(responseTimeRetrieveMessage)
}

func IncrementFollows() {
	followed.Inc()
	//fmt.Println("Followed incremented")
}

func IncrementUnfollows() {
	unfollowed.Inc()
	//fmt.Println("Unfollowed incremented")
}

func IncrementRequests() {
	//httpRequests.Inc()
	//fmt.Println("HTTP requests incremented")
}

func ObserveResponseTime(route string, method string, duration float64) {
	switch {
	case route == "/register":
		responseTimeRegister.WithLabelValues(route, method).Observe(duration)
		//fmt.Println("Register")
	case strings.Contains(route, "/msgs"):
		switch method {
		case "GET":
			responseTimeRetrieveMessage.WithLabelValues(route, method).Observe(duration)
			//fmt.Println("Retrive messages")
		case "POST":
			responseTimeSendMessage.WithLabelValues(route, method).Observe(duration)
			//fmt.Println("Send message")
		}
	}
}

func getCounterOpts(name string, help string) prometheus.CounterOpts {
	return prometheus.CounterOpts{
		Name: name,
		Help: help,
	}
}

func getHistogramOpts(name string, help string) prometheus.HistogramOpts {
	return prometheus.HistogramOpts{
		Name: name,
		Help: help,
	}
}

func getGaugeOpts(name string, help string) prometheus.GaugeOpts {
	return prometheus.GaugeOpts{
		Name: name,
		Help: help,
	}
}

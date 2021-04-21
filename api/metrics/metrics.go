package metrics

import (
	"strings"
	"time"

	dto "github.com/lauralunddahl/DevOps-GroupF/api/dto"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	cpu "github.com/shirou/gopsutil/cpu"
	mem "github.com/shirou/gopsutil/mem"
)

const interval = 5

var bytes_to_gigabytes = float64(1073741824)

//Prometheus objects needed as global variables:
var followed = promauto.NewCounter(prometheus.CounterOpts{
	Name: "followed_users",
	Help: "Number of follows",
})
var unfollowed = promauto.NewCounter(prometheus.CounterOpts{
	Name: "unfollowed_users",
	Help: "Number of unfollows",
})
var httpRequests = promauto.NewCounter(prometheus.CounterOpts{
	Name: "http_requests",
	Help: "Number of http requests",
})
var users = prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "users_registered",
	Help: "Total number of users registered",
})
var averageFollowers = prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "average_followers_per_user",
	Help: "Number of average followers per user",
})
var averagePosts = prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "average_posts_per_user",
	Help: "Number of average posts per user",
})
var cpuPercentage = prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "cpu_total_percentage",
	Help: "Displays the proportion of the cpu being used",
})
var memoryPercentage = prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "virtual_memory_percentage",
	Help: "Displays the proportion of the virtual memory being used",
})
var memoryAvailable = prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "virtual_memory_available",
	Help: "Information about the RAM available for programs to allocate in gigabytes",
})
var responseTimeRegister = prometheus.NewHistogram(prometheus.HistogramOpts{
	Name: "http_register_request_duration_seconds",
	Help: "Histogram of response time for registering a user in seconds",
})
var responseTimeSendMessage = prometheus.NewHistogram(prometheus.HistogramOpts{
	Name: "http_send_message_request_duration_seconds",
	Help: "Histogram of response time for sending a message in seconds",
})
var responseTimeRetrieveMessage = prometheus.NewHistogram(prometheus.HistogramOpts{
	Name: "http_retrieve_message_request_duration_seconds",
	Help: "Histogram of response time for retrieving a message in seconds",
})

func RecordMetrics() {
	registerMetrics()
	databaseMetrics()
	cpuMetric()
	virtualMemoryMetrics()
}

func registerMetrics() {
	prometheus.MustRegister(users)
	prometheus.MustRegister(averageFollowers)
	prometheus.MustRegister(averagePosts)
	prometheus.MustRegister(cpuPercentage)
	prometheus.MustRegister(memoryPercentage)
	prometheus.MustRegister(memoryAvailable)
	prometheus.MustRegister(responseTimeRegister)
	prometheus.MustRegister(responseTimeSendMessage)
	prometheus.MustRegister(responseTimeRetrieveMessage)
}

func databaseMetrics() {
	go func() {
		for {
			var numberOfUsers = float64(dto.GetTotalNumberOfUsers())
			var numberOfFollowers = float64(dto.GetTotalNumberOfFollowerEntries())
			var numberOfPosts = float64(dto.GetTotalNumberOfMessages())

			users.Set(numberOfUsers)
			averageFollowers.Set(numberOfFollowers / numberOfUsers)
			averagePosts.Set(numberOfPosts / numberOfUsers)
			time.Sleep(interval * time.Second)
		}
	}()
}

func cpuMetric() {
	go func() {
		for {
			c, _ := cpu.Percent(0, false)
			cpuPercentage.Set(c[0])
			time.Sleep(interval * time.Second)
		}
	}()
}

func virtualMemoryMetrics() {
	go func() {
		for {
			v, _ := mem.VirtualMemory()
			memoryAvailable.Set(float64(v.Available) / bytes_to_gigabytes)
			memoryPercentage.Set(v.UsedPercent)
			time.Sleep(interval * time.Second)
		}
	}()
}

func IncrementFollows() {
	followed.Inc()
}

func IncrementUnfollows() {
	unfollowed.Inc()
}

func IncrementRequests() {
	httpRequests.Inc()
}

func ObserveResponseTime(route string, method string, duration float64) {
	switch {
	case route == "/register":
		responseTimeRegister.Observe(duration)
	case strings.Contains(route, "/msgs"):
		switch method {
		case "GET":
			responseTimeRetrieveMessage.Observe(duration)
		case "POST":
			responseTimeSendMessage.Observe(duration)
		}
	}
}

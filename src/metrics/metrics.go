package metrics

import (
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
	databaseMetrics()
	cpuMetric()
	virtualMemoryMetrics()
	responseTimeMetrics()
}

func databaseMetrics() {
	users := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "users_registered",
		Help: "Total number of users registered",
	})
	prometheus.MustRegister(users)

	averageFollowers := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "average_followers_per_user",
		Help: "Number of average followers per user",
	})
	prometheus.MustRegister(averageFollowers)

	averagePosts := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "average_posts_per_user",
		Help: "Number of average posts per user",
	})
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
	cpuPercentage := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "cpu_total_percentage",
		Help: "Displays the proportion of the cpu being used",
	})
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
	memoryPercentage := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "virtual_memory_percentage",
		Help: "Displays the proportion of the virtual memory being used",
	})
	prometheus.MustRegister(memoryPercentage)

	memoryAvailable := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "virtual_memory_available",
		Help: "Information about the RAM available for programs to allocate in gigabytes",
	})
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
	httpRequests.Inc()
	//fmt.Println("HTTP requests incremented")
}

func ObserveResponseTime(route string, method string, duration float64) {
	switch {
	case route == "/register":
		responseTimeRegister.Observe(duration)
		//fmt.Println("Register")
	case strings.Contains(route, "/msgs"):
		switch method {
		case "GET":
			responseTimeRetrieveMessage.Observe(duration)
			//fmt.Println("Retrive messages")
		case "POST":
			responseTimeSendMessage.Observe(duration)
			//fmt.Println("Send message")
		}
	}
}

// func getCounterOpts(name string, help string) prometheus.CounterOpts {
// 	return prometheus.CounterOpts{
// 		Name: name,
// 		Help: help,
// 	}
// }

// func getHistogramOpts(name string, help string) prometheus.HistogramOpts {
// 	return prometheus.HistogramOpts{
// 		Name: name,
// 		Help: help,
// 	}
// }

// func getGaugeOpts(name string, help string) prometheus.GaugeOpts {
// 	return prometheus.GaugeOpts{
// 		Name: name,
// 		Help: help,
// 	}
// }

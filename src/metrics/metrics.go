package metrics

import (
	"fmt"
	"time"
	"strings"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	dto "github.com/lauralunddahl/DevOps-GroupF/src/dto"
	mem "github.com/shirou/gopsutil/mem"
	cpu "github.com/shirou/gopsutil/cpu"
)

var bytes_to_gigabytes = float64(1073741824)

func RecordMetrics() {
	go func() {
		prometheus.MustRegister(users)
		prometheus.MustRegister(averageFollowers)
		prometheus.MustRegister(averagePosts)
		prometheus.MustRegister(memoryAvailable)
		prometheus.MustRegister(memoryPercentage)
		prometheus.MustRegister(cpuPercentage)
		prometheus.MustRegister(responseTimeRegister)
		prometheus.MustRegister(responseTimeSendMessage)
		prometheus.MustRegister(responseTimeRetrieveMessage)
		
		for {
			var numberOfUsers = float64(dto.GetTotalNumberOfUsers())
			var numberOfFollowers = float64(dto.GetTotalNumberOfFollowerEntries())
			var numberOfPosts = float64(dto.GetTotalNumberOfMessages())
			users.Set(numberOfUsers)
			averageFollowers.Set(numberOfFollowers/numberOfUsers)
			averagePosts.Set(numberOfPosts/numberOfUsers)
			v, _ := mem.VirtualMemory()
			c,_ := cpu.Percent(0,false)
			memoryAvailable.Set(float64(v.Available)/bytes_to_gigabytes)
			memoryPercentage.Set(v.UsedPercent)
			cpuPercentage.Set(c[0])
			time.Sleep(60*60*time.Second)
		}
	}()
}

func IncrementFollows() {
	go func() {
		followed.Inc()
	}()
}

func IncrementUnfollows() {
	go func() {
		unfollowed.Inc()
	}()
}

func ResponseTimeHistogram(route string, method string, duration float64){
	go func(){
		switch {
		case route == "/register":
			responseTimeRegister.WithLabelValues(route, method).Observe(duration)
			fmt.Println("Register")
		case strings.Contains(route, "/msgs"):
			switch method {
			case "GET":
				responseTimeRetrieveMessage.WithLabelValues(route, method).Observe(duration)
				fmt.Println("Retrive messages")
			case "POST":
				responseTimeSendMessage.WithLabelValues(route, method).Observe(duration)
				fmt.Println("Send message")
			}
		}
	}()
}

var (
	users = prometheus.NewGauge(prometheus.GaugeOpts{
		Name:       "users_registered",
		Help:       "Total number of users registered",
	})
	followed = promauto.NewCounter(prometheus.CounterOpts{
		Name:       "followed_users",
		Help:       "Number of follows",
	})
	unfollowed = promauto.NewCounter(prometheus.CounterOpts{
		Name:       "unfollowed_users",
		Help:       "Number of unfollows",
	})
	averageFollowers = prometheus.NewGauge(prometheus.GaugeOpts{
		Name:       "average_followers_per_user",
		Help:       "Number of average followers per user",
	})
	averagePosts = prometheus.NewGauge(prometheus.GaugeOpts{
		Name:       "average_posts_per_user",
		Help:       "Number of average posts per user",
	})
	memoryPercentage = prometheus.NewGauge(prometheus.GaugeOpts{
		Name:       "virtual_memory_percentage",
		Help:       "Displays the proportion of the virtual memory being used",
	})
	memoryAvailable = prometheus.NewGauge(prometheus.GaugeOpts{
		Name:       "virtual_memory_available",
		Help:       "Information about the RAM available for programs to allocate in gigabytes",
	})
	cpuPercentage = prometheus.NewGauge(prometheus.GaugeOpts{
		Name:       "cpu_total_percentage",
		Help:       "Displays the proportion of the cpu being used",
	})
	//At some point we could maybe expand this to also look at the status code
	responseTimeRegister = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:      "http_register_request_duration_seconds",
		Help:      "Histogram of response time for registering a user in seconds",
	}, []string{"route", "method"})
	responseTimeSendMessage = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:      "http_send_message_request_duration_seconds",
		Help:      "Histogram of response time for sending a message in seconds",
	}, []string{"route", "method"})
	responseTimeRetrieveMessage = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:      "http_retrieve_message_request_duration_seconds",
		Help:      "Histogram of response time for retrieving a message in seconds",
	}, []string{"route", "method"})
)

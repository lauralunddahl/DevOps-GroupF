package metrics

import (
	"fmt"
	"time"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	dto "github.com/lauralunddahl/DevOps-GroupF/src/dto"
	mem "github.com/shirou/gopsutil/mem"
	cpu "github.com/shirou/gopsutil/cpu"
)

func RecordMetrics() {
	go func() {
		prometheus.MustRegister(users)
		prometheus.MustRegister(averageFollowers)
		prometheus.MustRegister(averagePosts)
		prometheus.MustRegister(memoryFree)
		prometheus.MustRegister(memoryPercentage)
		prometheus.MustRegister(memoryActive)
		prometheus.MustRegister(cpuPercentage)
		
		for {
			var numberOfUsers = float64(dto.GetTotalNumberOfUsers())
			var numberOfFollowers = float64(dto.GetTotalNumberOfFollowerEntries())
			var numberOfPosts = float64(dto.GetTotalNumberOfMessages())
			users.Set(numberOfUsers)
			averageFollowers.Set(numberOfFollowers/numberOfUsers)
			averagePosts.Set(numberOfPosts/numberOfUsers)
			v, _ := mem.VirtualMemory()

			c,_ := cpu.Percent(0,false)

			memoryFree.Set(float64(v.Free))
			memoryPercentage.Set(v.UsedPercent)
			memoryActive.Set(float64(v.Active))
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
	memoryFree = prometheus.NewGauge(prometheus.GaugeOpts{
		Name:       "virtual_memory_free",
		Help:       "Information about the amount of free virtual memory",
	})
	memoryActive = prometheus.NewGauge(prometheus.GaugeOpts{
		Name:       "virtual_memory_active",
		Help:       "Information about the amount of active virtual memory",
	})
	cpuPercentage = prometheus.NewGauge(prometheus.GaugeOpts{
		Name:       "cpu_total_percentage",
		Help:       "Displays the proportion of the cpu being used",
	})
)

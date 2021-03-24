package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	dto "github.com/lauralunddahl/DevOps-GroupF/src/dto"
	mem "github.com/shirou/gopsutil/mem"
)

func RecordMetrics() {
	go func() {
		prometheus.MustRegister(users)
		prometheus.MustRegister(averageFollowers)
		prometheus.MustRegister(averagePosts)
		prometheus.MustRegister(memory)
		for {
			var numberOfUsers = float64(dto.GetTotalNumberOfUsers())
			var numberOfFollowers = float64(dto.GetTotalNumberOfFollowerEntries())
			var numberOfPosts = float64(dto.GetTotalNumberOfMessages())
			users.Set(numberOfUsers)
			averageFollowers.Set(numberOfFollowers/numberOfUsers)
			averagePosts.Set(numberOfPosts/numberOfUsers)
			v, _ := mem.VirtualMemory()
			memory.Set(float64(v.Free))
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
	memory = prometheus.NewGauge(prometheus.GaugeOpts{
		Name:       "virtual_memory",
		Help:       "Information about the virtual memory",
	})
)
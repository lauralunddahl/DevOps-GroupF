package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

/*func recordCPU() {
	go func() {
		for {
			avg by (job, instance, mode) (rate(node_cpu_seconds_total[5]))
			//count(process_cpu_seconds_total)
			time.Sleep(10 * time.Second)
		}
	}()
}*/

func RecordMetrics() {
	go func() {
		for {
			opsProcessed.Inc()
			time.Sleep(2 * time.Second)
		}
	}()
}

var (
	opsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "myapp_processed_ops_total",
		Help: "The total number of processed events",
	})
)

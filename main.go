package main

import (
	"github.com/kataras/iris/v12"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"push/controllers"
	"push/repository"
	"time"
)

var (
	processedTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "app_processed_total",
		Help: "Number of times ran",
	}, []string{"status"})
	responseStatus = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "response_status",
			Help: "Status of HTTP response",
		},
		[]string{"status"})
	latency = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace:  "push_notification",
			Name:       "latency_seconds",
			Help:       "Latency distributions.",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		},
		[]string{"method", "path"})
)

func RecordRequestLatency(c iris.Context) {
	start := time.Now()
	elapsed := time.Since(start).Seconds()

	latency.WithLabelValues(
		c.Method(),
		c.Path(),
	).Observe(elapsed)

	Counter.Inc()
	controllers.CreatePushNotificationHandler(c)
}

func doInit() {
	prometheus.MustRegister(processedTotal)
	prometheus.MustRegister(responseStatus)
	prometheus.MustRegister(latency)
	prometheus.MustRegister(Counter)
}

func main() {
	repository.Init()

	var errLoc error
	time.Local, errLoc = time.LoadLocation("America/Caracas")
	if errLoc != nil {
		log.Printf("error loading location %v\n", errLoc)
	}

	app := iris.New()

	v1 := app.Party("/v1")
	{
		pushNotification := v1.Party("/push_notification")
		{
			doInit()
			pushNotification.Post("/send", RecordRequestLatency)
			pushNotification.Get("/metrics", iris.FromStd(promhttp.HandlerFor(
				prometheus.DefaultGatherer,
				promhttp.HandlerOpts{
					EnableOpenMetrics: true,
				},
			)))
		}
		metrics := v1.Party("/metrics")
		{
			metrics.Get("/get", controllers.GetMetricsHandler)
		}
	}

	app.Listen(":8080")

}

var Counter = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "push_notification_count",
		Help: "No of request handled by Push Notification handler",
	},
)

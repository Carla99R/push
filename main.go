package main

import (
	prometheusMiddleware "github.com/iris-contrib/middleware/prometheus"
	"github.com/kataras/iris/v12"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"push/controllers"
	"push/repository"
	"time"
)

var (
	fileOnDisk     = prometheus.NewRegistry()
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
			Namespace:  "api",
			Name:       "latency_seconds",
			Help:       "Latency distributions.",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		},
		[]string{"method", "path"})
)

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
			m := prometheusMiddleware.New("push_notification", 0.3, 1.2, 5.0)
			doInit()
			pushNotification.Use(m.ServeHTTP)
			pushNotification.Post("/send", pushNotificationHandler)
			pushNotification.Get("/metrics", iris.FromStd(iris.FromStd(promhttp.InstrumentMetricHandler(fileOnDisk, promhttp.Handler()))))
		}
	}
	app.Listen(":8080")
}

func pushNotificationHandler(ctx iris.Context) {
	Counter.Inc()
	//totalRequests.WithLabelValues("/push_notification").Inc()
	controllers.CreatePushNotificationHandler(ctx)
}

var Counter = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "push_notification_count",
		Help: "No of request handled by Push Notification handler",
	},
)

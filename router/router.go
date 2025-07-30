package router

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	_ "github.com/swaggo/files"
	_ "github.com/yilei-pixocial/kubeapi/pkg/api"
	"github.com/yilei-pixocial/kubeapi/pkg/service"
	_ "github.com/yilei-pixocial/kubeapi/router/middleware"
	"log"
	"time"
)

func SetRoutes(app *iris.Application) {

	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path", "status"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	app.Use(func(ctx iris.Context) {
		start := time.Now()
		ctx.Next()
		duration := time.Since(start).Seconds()
		status := ctx.GetStatusCode()
		method := ctx.Method()
		path := ctx.Path()

		requestCounter.WithLabelValues(method, path, fmt.Sprintf("%d", status)).Inc()
		requestDuration.WithLabelValues(method, path, fmt.Sprintf("%d", status)).Observe(duration)
	})

	app.Get("/metrics", iris.FromStd(promhttp.Handler()))

	rootApi := app.Party("api/v1")

	k8s, err := service.NewK8sService()
	if err != nil {
		log.Fatalln(fmt.Errorf("new k8s service error: %v", err))
		return
	}

	rootApi.Get("/k8s/namespaces", k8s.GetNamespaces) // 获取namespace
	rootApi.Get("/k8s/services", k8s.GetServices)     // 获取service

}

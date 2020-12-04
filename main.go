package main

import (
	"flag"
	"net/http"
	"os"

	"github.com/eliecharra/hilink-exporter/collector"
	"github.com/eliecharra/hilink-exporter/logger"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

var (
	url         = flag.String("endpoint", "http://192.168.8.1/", "api endpoint")
	username    = flag.String("user", "admin", "username")
	password    = flag.String("password", "", "password")
	port        = flag.String("port", ":9770", "port number to listen on")
	metricsPath = flag.String("path", "/metrics", "path to answer requests on")
	logLevel    = flag.String("log-level", "info", "log level")

	withTraffic = flag.Bool("traffic", true, "fetch traffic related stats")
	withSignal  = flag.Bool("signal", true, "fetch signal related stats")
	withWan     = flag.Bool("wan", true, "fetch wan related stats")
)

func main() {

	flag.Parse()

	logger.Init(logLevel)

	err := start()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func start() error {

	h, err := createMetricsHandler()
	if err != nil {
		return err
	}
	http.Handle(*metricsPath, h)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`<html>
			<head><title>Hilink Exporter</title></head>
			<body>
			<h1>Hilink Exporter</h1>
			<p><a href="` + *metricsPath + `">Metrics</a></p>
			</body>
			</html>`))
	})

	log.Info("Listening on ", *port)
	log.Fatal(http.ListenAndServe(*port, nil))

	return nil
}

func createMetricsHandler() (http.Handler, error) {
	opts := collector.Options{
		Url:           url,
		Username:      username,
		Password:      password,
		EnableWan:     *withWan,
		EnableTraffic: *withTraffic,
		EnableSignal:  *withSignal,
	}
	c := collector.NewCollector(opts)
	registry := prometheus.NewRegistry()
	err := registry.Register(c)
	if err != nil {
		return nil, err
	}
	return promhttp.HandlerFor(registry,
		promhttp.HandlerOpts{
			ErrorLog:      log.New(),
			ErrorHandling: promhttp.ContinueOnError,
		}), nil
}

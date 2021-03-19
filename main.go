package main

import (
	"flag"
	"log"
	"net/http"
	"strconv"

	"github.com/deemon87/ip_exporter/collector"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	var port = flag.Int("port", 19001, "listen port")
	flag.Parse()

	networks := collector.NewCollector()
	prometheus.MustRegister(networks)

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(*port), nil))
}

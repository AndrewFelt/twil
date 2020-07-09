package main

import (
	"flag"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

//Account - account token for Twilio
var Account = flag.String("account", "", "The port metrics are exposed on")

//Token - access token for Twilio
var Token = flag.String("token", "", "The port metrics are exposed on")

//Port - Port metrics are exposed on, include the colon. E.G. :2112
var Port = flag.String("port", ":2112", "The port metrics are exposed on")

func main() {

	flag.Parse()

	usage := newUsageCollector()
	prometheus.MustRegister(usage)

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(*Port, nil)
}

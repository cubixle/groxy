package main

import (
	"log"
	"net/http"
	"os"

	"github.com/cubixle/groxy/internal/logging"
	"github.com/cubixle/groxy/internal/reverseproxy"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gopkg.in/yaml.v2"

	_ "github.com/cubixle/groxy/internal/metrics"
)

func main() {
	data, err := os.ReadFile("./groxy.yml")
	if err != nil {
		log.Fatal(err)
	}

	eps := &reverseproxy.File{}

	err = yaml.Unmarshal(data, eps)
	if err != nil {
		log.Fatal(err)
	}

	logger := logging.Logger()
	defer logger.Sync()

	handler := reverseproxy.NewReverseProxy(eps.Endpoints, logger)
	go func() {
		logger.Debug("starting metrics server")
		if err := http.ListenAndServe("localhost:9090", promhttp.Handler()); err != nil {
			log.Fatal(err)
		}
	}()

	logger.Debug("starting reverse proxy")
	if err := http.ListenAndServe("localhost:8080", handler); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

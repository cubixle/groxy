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
	cfgFile := getEnvVar("GROXY_CONFIG_FILE", "./groxy.yml")

	data, err := os.ReadFile(cfgFile)
	if err != nil {
		log.Fatal(err)
	}

	eps := &reverseproxy.File{}

	err = yaml.Unmarshal(data, eps)
	if err != nil {
		log.Fatal(err)
	}

	logger := logging.Logger(eps.Debug)
	defer logger.Sync()

	handler := reverseproxy.NewReverseProxy(eps.Endpoints, logger)
	go func() {
		logger.Info("starting metrics server")
		port := getEnvVar("GROXY_METRICS_PORT", "9090")
		if err := http.ListenAndServe(":"+port, promhttp.Handler()); err != nil {
			log.Fatal(err)
		}
	}()

	port := getEnvVar("GROXY_PORT", "8080")
	logger.Info("starting reverse proxy")
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func getEnvVar(name, defaultValue string) string {
	v := os.Getenv(name)
	if v == "" {
		return defaultValue
	}

	return v
}

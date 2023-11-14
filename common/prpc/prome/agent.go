package prome

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"im/logger"
	"net/http"
	"sync"
)

var once sync.Once

func StartAgent(host string, port int) {
	go func() {
		once.Do(func() {
			http.Handle("/", promhttp.Handler())
			addr := fmt.Sprintf("%s:%d", host, port)
			logger.Infof("Starting prometheus agent at %s", addr)
			if err := http.ListenAndServe(addr, nil); err != nil {
				logger.Error(err)
			}
		})
	}()
}

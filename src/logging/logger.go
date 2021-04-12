package logging

import (
	log "github.com/sirupsen/logrus"
)

func Logging() {
	var logger = log.New()
	addr := "http://164.90.254.78:8080"
	
	logger.WithField("addr", addr).Info("starting server")	
}
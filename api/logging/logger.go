package logging

import (
	log "github.com/sirupsen/logrus"
)

func Logging() {
	var logger = log.New()
	port := "8080"
	logger.WithField("port", port).Info("starting server")	
}
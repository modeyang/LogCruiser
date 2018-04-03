package main

import (
	log "github.com/sirupsen/logrus"
	"testing"
	"os"
)

func TestLogrus(T *testing.T) {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.WarnLevel)
	log.WithFields(log.Fields{
		"omg":    true,
		"number": 122,
	}).Warn("The group's number increased tremendously!")

	contextLogger := log.WithFields(log.Fields{
		"common": "this is a common field",
		"other": "I also should be logged always",
	})
	contextLogger.Warn("test")
}


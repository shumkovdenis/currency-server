package main

import (
	"flag"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/shumkovdenis/currency-server/load"
	"github.com/shumkovdenis/currency-server/rate"
	"github.com/shumkovdenis/currency-server/server"
)

func main() {
	var (
		port     int
		length   int
		timeout  time.Duration
		logLevel string
	)

	flag.IntVar(&port, "port", 3030, "Port.")
	flag.IntVar(&length, "length", 60, "Length rates.")
	flag.DurationVar(&timeout, "timeout", 1*time.Minute, "Fail timeout.")
	flag.StringVar(&logLevel, "log", "info", "Log level.")
	flag.Parse()

	level, err := log.ParseLevel(logLevel)
	if err != nil {
		log.Warning("Set default log level")
		level = log.InfoLevel
	}
	log.SetLevel(level)
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	log.WithFields(log.Fields{
		"port":    port,
		"length":  length,
		"timeout": timeout,
		"log":     logLevel,
	}).Info("Start")

	loadService := load.NewTruefx("http://webrates.truefx.com/rates/connect.html?f=csv")

	rateService := rate.NewSingle(length, timeout, loadService)
	rateService.Start()

	server.Start(port, rateService)
}

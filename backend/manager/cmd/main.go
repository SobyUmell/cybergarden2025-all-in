package main

import (
	"io"
	"log"
	"manager/internal/app"
	"manager/internal/config"
	"os"
	"os/signal"
	"syscall"

	"github.com/PrototypeSirius/ruglogger/logger"
	"github.com/sirupsen/logrus"
)

var level string

func main() {
	cfg := config.MustLoad()
	log, logCloser := logInit(cfg.Env)
	defer func() {
		if err := logCloser(); err != nil {
			os.Stderr.WriteString("Error closing log file: " + err.Error() + "\n")
		}
	}()
	log.Info("Config has been successfully loaded")
	application := app.New(log, *cfg)
	log.Info("Application has been successfully initialized")
	go func() {
		application.HTTPApp.MustRun()
	}()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop

	log.Info("Gracefully stopped")
}

func logInit(env string) (*logrus.Logger, func() error) {
	switch env {
	case "production":
		level = "warn"
	case "local":
		level = "info"
	default:
		level = "debug"
	}
	logFile, err := os.OpenFile("app.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Critical error: failed to open log file: %v\n", err)
	}
	multiWriter := io.MultiWriter(os.Stdout, logFile)
	logger.Init(
		logger.WithLevel(level),
		logger.WithOutput(multiWriter),
	)
	log := logger.Get()
	log.Info("The logging system has been successfully initialized.")
	return log, logFile.Close
}

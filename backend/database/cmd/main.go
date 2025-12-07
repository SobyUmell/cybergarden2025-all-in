package main

import (
	"database/internal/app"
	"database/internal/config"
	"io"
	"log"
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

	application, closeDB := app.New(log, cfg.Database, cfg.GRPC.Port)
	log.Info("Application has been successfully initialized.")
	go func() {
		application.GRPCServer.MustRun()
	}()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop
	application.GRPCServer.Stop()
	closeDB()
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
	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
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

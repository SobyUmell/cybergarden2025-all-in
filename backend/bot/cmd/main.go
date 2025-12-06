package main

import (
	"bot/internal/app"
	"bot/internal/config"
	"context"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/PrototypeSirius/ruglogger/apperror"
	"github.com/PrototypeSirius/ruglogger/logger"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
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
	application := app.New(log, cfg)
	log.Info("Application has been successfully initialized")
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	g, gCtx := errgroup.WithContext(ctx)
	g.Go(func() error {
		return application.BotServer.Run(gCtx)
	})
	g.Go(func() error {
		return application.GRPCServer.Run()
	})
	g.Go(func() error {
		<-gCtx.Done()
		logger.Info("Shutting down the application")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := application.BotServer.Stop(shutdownCtx); err != nil {
			appErr := apperror.SystemError(err, 1001, "error stopping bot server")
			logger.LogOnError(appErr, "stopping bot server")
		}
		application.GRPCServer.Stop()
		return nil
	})
	if err := g.Wait(); err != nil {
		appErr := apperror.SystemError(err, 1002, "error running application")
		logger.LogOnError(appErr, "running application")
	} else {
		log.Info("Application has been successfully stopped")
	}
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

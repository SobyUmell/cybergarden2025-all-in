package app

import (
	httpapp "manager/internal/app/http"
	"manager/internal/config"
	botrepo "manager/internal/repository/bot"
	databaserepo "manager/internal/repository/database"
	mlrepo "manager/internal/repository/ml"
	"manager/internal/router"
	httpservice "manager/internal/services"

	"github.com/PrototypeSirius/ruglogger/logger"
	"github.com/sirupsen/logrus"
)

type App struct {
	HTTPApp *httpapp.HTTPApp
}

func New(log *logrus.Logger, cfg config.Config) *App {
	botClient, err := botrepo.New(log, cfg.Client.Bot.Host, cfg.Client.Bot.Port)
	if err != nil {
		logger.FatalOnError(err, "error init bot client")
	}

	databaseClient, err := databaserepo.New(log, cfg.Client.Database.Host, cfg.Client.Database.Port)
	if err != nil {
		logger.FatalOnError(err, "error init database client")
	}

	mlClient := mlrepo.New(log, cfg.Client.ML.Host, cfg.Client.ML.Port)

	managerService := httpservice.New(log, botClient, databaseClient, mlClient)

	httpHandler := router.New(log, managerService)

	httpServer, engine := httpapp.New(log, cfg.HttpServer.Port)

	httpHandler.RouterRegister(engine)

	return &App{HTTPApp: httpServer}
}

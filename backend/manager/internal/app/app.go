package app

import (
	httpapp "manager/internal/app/http"
	"manager/internal/config"
	botrepo "manager/internal/repository/bot"
	databaserepo "manager/internal/repository/database"
	"manager/internal/router"
	httpservice "manager/internal/services"

	"github.com/PrototypeSirius/ruglogger/logger"
	"github.com/sirupsen/logrus"
)

type App struct {
	HTTPApp *httpapp.HTTPApp
}

func New(log *logrus.Logger, cfg config.Config) *App {
	// 1. Initialize Repositories (gRPC Clients)
	botClient, err := botrepo.New(log, cfg.Client.Bot.Host, cfg.Client.Bot.Port)
	if err != nil {
		logger.FatalOnError(err, "error init bot client")
	}

	databaseClient, err := databaserepo.New(log, cfg.Client.Database.Host, cfg.Client.Database.Port)
	if err != nil {
		logger.FatalOnError(err, "error init database client")
	}

	// 2. Initialize Service (Business Logic)
	managerService := httpservice.New(log, botClient, databaseClient)

	// 3. Initialize Handler (Router)
	httpHandler := router.New(log, managerService)

	// 4. Initialize HTTP Server (Gin)
	httpServer, engine := httpapp.New(log, cfg.HttpServer.Port)

	// 5. Register Routes
	httpHandler.RouterRegister(engine)

	return &App{HTTPApp: httpServer}
}

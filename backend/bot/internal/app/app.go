package app

import (
	botapp "bot/internal/app/bot"
	grpcapp "bot/internal/app/grpc"
	"bot/internal/config"
	botrepo "bot/internal/repository/bot"
	botservice "bot/internal/services"

	"github.com/sirupsen/logrus"
)

type App struct {
	GRPCServer *grpcapp.App
	BotServer  *botapp.App
}

func New(log *logrus.Logger, cfg *config.Config) *App {
	botServer, bot := botapp.New(log, cfg.Bot)
	botRepo := botrepo.New(bot, cfg.Bot.Token)
	botService := botservice.New(log, botRepo)
	grpcServer := grpcapp.New(log, botService, cfg.GRPC.Port)
	return &App{GRPCServer: grpcServer, BotServer: botServer}
}

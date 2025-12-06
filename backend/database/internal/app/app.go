package app

import (
	grpcapp "database/internal/app/grpc"
	"database/internal/config"
	dbrepo "database/internal/repository/database"
	dbservice "database/internal/services"

	"github.com/sirupsen/logrus"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(log *logrus.Logger, cfg config.DatabaseConfig, port int) (*App, func() error) {
	dbrepo, closeDB := dbrepo.New(log, cfg)
	dbservice := dbservice.New(log, dbrepo)
	grpcapp := grpcapp.New(log, dbservice, port)
	return &App{GRPCServer: grpcapp}, closeDB
}

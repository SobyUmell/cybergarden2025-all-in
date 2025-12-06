package httpservice

import (
	"context"
	"fmt"
	model "manager/internal/models"

	"github.com/PrototypeSirius/ruglogger/logger"
	"github.com/sirupsen/logrus"
)

type BotRepository interface {
	Auth(ctx context.Context, initData string) (int64, error)
}

type DatabaseRepository interface {
	AddUser(ctx context.Context, uid int64) error
	AddTransaction(ctx context.Context, uid int64, t model.Transaction) error
	DeleteTransaction(ctx context.Context, uid, tid int64) error
	EditTransaction(ctx context.Context, uid int64, t model.Transaction) error
	RequestUserTransactions(ctx context.Context, uid int64) ([]model.Transaction, error)
}

type ManagerService struct {
	log *logrus.Logger
	bot BotRepository
	db  DatabaseRepository
}

func New(log *logrus.Logger, bot BotRepository, db DatabaseRepository) *ManagerService {
	return &ManagerService{
		log: log,
		bot: bot,
		db:  db,
	}
}

func (s *ManagerService) AuthUser(ctx context.Context, initData string) (int64, error) {
	uid, err := s.bot.Auth(ctx, initData)
	if err != nil {
		return 0, err 
	}
	if err := s.db.AddUser(ctx, uid); err != nil {
		logger.LogOnError(err, fmt.Sprintf("Failed to ensure user %d exists in DB", uid))
		return 0, err
	}

	return uid, nil
}

func (s *ManagerService) AddTransaction(ctx context.Context, uid int64, t model.Transaction) error {
	return s.db.AddTransaction(ctx, uid, t)
}

func (s *ManagerService) DeleteTransaction(ctx context.Context, uid, tid int64) error {
	return s.db.DeleteTransaction(ctx, uid, tid)
}

func (s *ManagerService) EditTransaction(ctx context.Context, uid int64, t model.Transaction) error {
	return s.db.EditTransaction(ctx, uid, t)
}

func (s *ManagerService) GetHistory(ctx context.Context, uid int64) ([]model.Transaction, error) {
	return s.db.RequestUserTransactions(ctx, uid)
}
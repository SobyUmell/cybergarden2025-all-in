package botservice

import (
	"bot/internal/models"
	botrepo "bot/internal/repository/bot"
	"context"

	"github.com/sirupsen/logrus"
)

type BotService struct {
	log  *logrus.Logger
	auth Auther
	send SendMessager
}

func New(log *logrus.Logger, b *botrepo.Bot) *BotService {
	return &BotService{log: log, auth: b, send: b}
}

type Auther interface {
	Auth(ctx context.Context, authData string) (models.UserResponse, error)
}

type SendMessager interface {
	SendMessage(ctx context.Context, uid int64, text string) (string, error)
}

func (b *BotService) Auth(ctx context.Context, authData string) (models.UserResponse, error) {
	return b.auth.Auth(ctx, authData)
}

func (b *BotService) SendMessage(ctx context.Context, uid int64, text string) (string, error) {
	return b.send.SendMessage(ctx, uid, text)
}

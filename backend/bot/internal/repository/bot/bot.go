package botrepo

import (
	"bot/internal/models"
	"context"
	"time"

	"github.com/PrototypeSirius/ruglogger/apperror"
	"github.com/go-telegram/bot"
	initdata "github.com/telegram-mini-apps/init-data-golang"
)

type Bot struct {
	Bot   *bot.Bot
	Token string
}

func New(bot *bot.Bot, token string) *Bot {
	return &Bot{Bot: bot, Token: token}
}

func (b *Bot) Auth(ctx context.Context, authData string) (models.UserResponse, error) {
	expIn := 24 * time.Hour
	if err := initdata.Validate(authData, b.Token, expIn); err != nil {
		return models.UserResponse{Authorized: false, UserID: 0, Error: err.Error()}, apperror.BadRequestError(err, 1081, "Invalid or expired authorization data.")
	}
	initData, err := initdata.Parse(authData)
	if err != nil {
		return models.UserResponse{Authorized: false, UserID: 0, Error: err.Error()}, apperror.SystemError(err, 1082, "Failed to parse valid init data.")
	}
	return models.UserResponse{Authorized: true, UserID: initData.User.ID, Error: ""}, nil
}

func (b *Bot) SendMessage(ctx context.Context, uid int64, text string) (string, error) {
	_, err := b.Bot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: uid,
		Text:   text,
	})
	if err != nil {
		return "faled", apperror.SystemError(err, 1083, "Failed to send message.")
	}
	return "success", nil
}

package bothandler

import (
	"bot/internal/config"
	"context"
	"errors"
	"fmt"

	"github.com/PrototypeSirius/ruglogger/apperror"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/sirupsen/logrus"
)

func Start(ctx context.Context, log *logrus.Logger, b *bot.Bot) {
	log.Info("Starting bot")
	go b.StartWebhook(ctx)
}

func SetWebhook(ctx context.Context, log *logrus.Logger, b *bot.Bot, cfg config.BotConfig) error {
	webHookURL := fmt.Sprintf("%s%s", cfg.WebURL, "bot/webhook")
	log.Info("Setting webhook", logrus.Fields{"url": webHookURL})
	ok, err := b.SetWebhook(ctx, &bot.SetWebhookParams{
		URL:         webHookURL,
		SecretToken: cfg.WebhookToken,
	})
	if err != nil {
		return apperror.SystemError(err, 1051, "error set webhook")
	}
	if !ok {
		return apperror.SystemError(errors.New("webhook was not set (api returned false)"), 1052, "error set webhook")
	}
	return nil
}

func HandleStart(cfg config.BotConfig) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		kb := &models.InlineKeyboardMarkup{InlineKeyboard: [][]models.InlineKeyboardButton{{{Text: "OPEN", WebApp: &models.WebAppInfo{URL: cfg.WebURL}}}}}
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      update.Message.Chat.ID,
			Text:        "CentKeeper MiniApp",
			ReplyMarkup: kb,
		})
	}
}

func HandleHelp(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Доступные команды:\n/start - Начало работы\n/help - Помощь",
	})
}

func HandleDefault(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message != nil && update.Message.Text != "" {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Введите команду. /help - список всех команд",
		})
	}
}

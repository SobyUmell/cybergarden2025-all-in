package botapp

import (
	"bot/internal/config"
	bothandler "bot/internal/handlers/bot"
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/PrototypeSirius/ruglogger/apperror"
	"github.com/PrototypeSirius/ruglogger/logger"
	"github.com/go-telegram/bot"
	"github.com/sirupsen/logrus"
)

type App struct {
	log        *logrus.Logger
	bot        *bot.Bot
	httpServer *http.Server
	cfg        config.BotConfig
}

func New(log *logrus.Logger, cfg config.BotConfig) (*App, *bot.Bot) {
	opts := []bot.Option{
		bot.WithDefaultHandler(bothandler.HandleDefault),
		bot.WithWebhookSecretToken(cfg.WebhookToken),
		// bot.WithDebug(),
	}
	b, err := bot.New(cfg.Token, opts...)
	if err != nil {
		appErr := apperror.SystemError(err, 1011, "error creating bot client")
		logger.FatalOnError(appErr, "creating bot handler")
	}
	b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, bothandler.HandleStart(cfg))
	b.RegisterHandler(bot.HandlerTypeMessageText, "/help", bot.MatchTypeExact, bothandler.HandleHelp)
	return &App{
		log: log,
		bot: b,
		cfg: cfg,
	}, b
}

func (a *App) Run(ctx context.Context) error {
	if err := bothandler.SetWebhook(ctx, a.log, a.bot, a.cfg); err != nil {
		return err
	}
	bothandler.Start(ctx, a.log, a.bot)
	mux := http.NewServeMux()
	a.httpServer = &http.Server{
		Addr:         fmt.Sprintf(":%d", a.cfg.Port),
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	mux.Handle("POST /bot/webhook", a.bot.WebhookHandler())
	a.log.Info("Starting Bot Webhook Server on port", logrus.Fields{"port": a.cfg.Port})
	err := a.httpServer.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return apperror.SystemError(err, 1012, "error starting webhook server")
	}
	return nil
}

func (a *App) Stop(ctx context.Context) error {
	a.log.Info("Stopping Bot Webhook Server.")
	if a.httpServer != nil {
		return a.httpServer.Shutdown(ctx)
	}
	return nil
}

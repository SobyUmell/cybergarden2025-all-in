package botrepo

import (
	"context"
	"errors"
	"fmt"

	cyberbott "github.com/PrototypeSirius/protos_service/gen/cybergarden/bot"
	"github.com/PrototypeSirius/ruglogger/apperror"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type BotClient struct {
	bot cyberbott.BotClient
	log *logrus.Logger
}

func New(log *logrus.Logger, host string, port int) (*BotClient, error) {
	botaddr := fmt.Sprintf("%s:%d", host, port)
	bcc, err := grpc.NewClient(botaddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, apperror.SystemError(err, 1021, "error load bot service")
	}
	return &BotClient{bot: cyberbott.NewBotClient(bcc), log: log}, nil
}

func (c *BotClient) Auth(ctx context.Context, initData string) (int64, error) {
	resp, err := c.bot.Auth(ctx, &cyberbott.AuthRequest{
		AuthData: initData,
	})
	if err != nil {
		return 0, apperror.SystemError(err, 1022, "grpc auth request failed")
	}
	if !resp.GetAuthorized() {
		return 0, apperror.BadRequestError(errors.New(resp.GetError()), 1023, "user not authorized by bot service")
	}
	return resp.GetUserID(), nil
}

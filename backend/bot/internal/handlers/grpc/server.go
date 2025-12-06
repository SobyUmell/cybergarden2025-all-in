package grpchandler

import (
	"bot/internal/models"
	"context"
	"errors"

	cyberbott "github.com/PrototypeSirius/protos_service/gen/cybergarden/bot"
	"github.com/PrototypeSirius/ruglogger/apperror"
	"github.com/PrototypeSirius/ruglogger/logger"
	"google.golang.org/grpc"
)

type BotService interface {
	Auth(ctx context.Context, AuthData string) (models.UserResponse, error)
	SendMessage(ctx context.Context, uid int64, text string) (string, error)
}

type serverAPI struct {
	cyberbott.UnimplementedBotServer
	bot BotService
}

func Register(s *grpc.Server, bot BotService) {
	cyberbott.RegisterBotServer(s, &serverAPI{bot: bot})
}

func (s *serverAPI) Auth(ctx context.Context, req *cyberbott.AuthRequest) (*cyberbott.AuthResponse, error) {
	if req.GetAuthData() == "" {
		appErr := apperror.BadRequestError(errors.New("empty"), 1061, "empty auth data")
		logger.LogOnError(appErr, "error in auth")
		return nil, appErr
	}
	userData, err := s.bot.Auth(ctx, req.GetAuthData())
	if err != nil {
		logger.LogOnError(err, "error in auth")
		return nil, err
	}
	return &cyberbott.AuthResponse{Authorized: userData.Authorized, UserID: userData.UserID, Error: userData.Error}, nil
}

func (s *serverAPI) SendMessage(ctx context.Context, req *cyberbott.SendMessageRequest) (*cyberbott.SendMessageResponse, error) {
	if req.GetUserID() == 0 {
		appErr := apperror.BadRequestError(errors.New("empty"), 1062, "empty uid")
		logger.LogOnError(appErr, "error in send message")
		return nil, appErr
	}
	if req.GetMessage() == "" {
		appErr := apperror.BadRequestError(errors.New("empty"), 1063, "empty text")
		logger.LogOnError(appErr, "error in send message")
		return nil, appErr
	}
	text, err := s.bot.SendMessage(ctx, req.GetUserID(), req.GetMessage())
	if err != nil {
		logger.LogOnError(err, "error in send message")
		return nil, err
	}
	return &cyberbott.SendMessageResponse{Error: text}, nil
}

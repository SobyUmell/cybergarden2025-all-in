package grpcapp

import (
	grpchandler "bot/internal/handlers/grpc"
	"context"
	"fmt"
	"net"
	"time"

	"github.com/PrototypeSirius/ruglogger/apperror"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type App struct {
	log        *logrus.Logger
	grpcServer *grpc.Server
	port       int
}

func New(log *logrus.Logger, serverAPI grpchandler.BotService, port int) *App {
	opts := []grpc.ServerOption{grpc.UnaryInterceptor(loggingInterceptor(log))}
	gRPCServer := grpc.NewServer(opts...)
	grpchandler.Register(gRPCServer, serverAPI)
	return &App{log: log, grpcServer: gRPCServer, port: port}
}

func (a *App) Run() error {
	a.log.Info("Starting the gRPC server")
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return apperror.SystemError(err, 1021, "error starting listener for gRPC server")
	}
	a.log.Info("grpc server is running", l.Addr().String())
	if err := a.grpcServer.Serve(l); err != nil {
		return apperror.SystemError(err, 1022, "error starting gRPC server")
	}
	return nil
}

func (a *App) Stop() error {
	a.log.Info("Stopping the gRPC server...")
	a.grpcServer.GracefulStop()
	return nil
}

func loggingInterceptor(log *logrus.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()
		resp, err := handler(ctx, req)
		duration := time.Since(start)
		fields := logrus.Fields{
			"method":   info.FullMethod,
			"duration": duration.String(),
			"req":      req,
		}
		if err != nil {
			fields["error"] = err.Error()
			if st, ok := status.FromError(err); ok {
				fields["code"] = st.Code().String()
			} else {
				fields["code"] = codes.Unknown.String()
			}
			log.WithFields(fields).Error("gRPC request failed")
		} else {
			fields["code"] = codes.OK.String()
			log.WithFields(fields).Info("gRPC request success")
		}
		return resp, err
	}
}

package httpapp

import (
	"fmt"

	"github.com/PrototypeSirius/ruglogger/apperror"
	"github.com/PrototypeSirius/ruglogger/logger"
	"github.com/PrototypeSirius/ruglogger/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type HTTPApp struct {
	log       *logrus.Logger
	ginServer *gin.Engine
	port      int
}

func New(log *logrus.Logger, port int) (*HTTPApp, *gin.Engine) {
	gin.ForceConsoleColor()
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.StructuredLogHandler())
	r.Use(middleware.ErrorHandler())
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://127.0.0.1:8080"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length", "Authorization"},
		AllowCredentials: true,
	}))
	return &HTTPApp{
		log:       log,
		ginServer: r,
		port:      port,
	}, r
}

func (a *HTTPApp) MustRun() {
	if err := a.Run(); err != nil {
		logger.FatalOnError(err, "Error starting REST server")
	}
}

func (a *HTTPApp) Run() error {
	addr := fmt.Sprintf(":%d", a.port)
	a.log.Info(fmt.Sprintf("Starting REST server on %s via Gin", addr))
	if err := a.ginServer.Run(addr); err != nil {
		return apperror.SystemError(err, 1031, "error starting REST server")
	}
	return nil
}

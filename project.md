### DIRECTORY . FOLDER STRUCTURE ###
FILE .env
FILE .gitignore
DIR backend/
    DIR bot/
        FILE app.log
        DIR cmd/
            FILE main.go
        DIR config/
            FILE config.yaml
        FILE Dockerfile
        FILE go.mod
        FILE go.sum
        DIR internal/
            DIR app/
                FILE app.go
                DIR bot/
                    FILE app.go
                DIR grpc/
                    FILE app.go
            DIR config/
                FILE config.go
            DIR handlers/
                DIR bot/
                    FILE server.go
                DIR grpc/
                    FILE server.go
            DIR models/
                FILE model.go
            DIR repository/
                DIR bot/
                    FILE bot.go
            DIR services/
                FILE bot.go
    DIR database/
        DIR cmd/
            FILE main.go
        DIR config/
            FILE config.yaml
        FILE Dockerfile
        FILE go.mod
        FILE go.sum
        DIR internal/
            DIR app/
                FILE app.go
                DIR grpc/
                    FILE app.go
            DIR config/
                FILE config.go
            DIR handlers/
                DIR grpc/
                    FILE server.go
            DIR migrations/
                FILE 1_init_schema.down.sql
                FILE 1_init_schema.up.sql
            DIR migrator/
                FILE migrator.go
            DIR models/
                FILE models.go
            DIR repository/
                DIR database/
                    FILE db.go
            DIR services/
                FILE db.go
    DIR database_save/
    DIR manager/
        DIR cmd/
            FILE main.go
        DIR config/
            FILE config.yaml
        FILE Dockerfile
        FILE go.mod
        FILE go.sum
        DIR internal/
            DIR app/
                FILE app.go
                DIR http/
                    FILE app.go
            DIR config/
                FILE config.go
            DIR models/
                FILE models.go
            DIR repository/
                DIR bot/
                    FILE client.go
                DIR database/
                    FILE client.go
            DIR router/
                FILE router.go
            DIR services/
                FILE http.go
    FILE README.md
FILE docker-compose.yml
DIR frontend/
    FILE .gitignore
    FILE components.json
    FILE next.config.ts
    FILE package.json
    FILE pnpm-lock.yaml
    FILE postcss.config.mjs
    DIR public/
        FILE file.svg
        FILE globe.svg
        FILE next.svg
        FILE vercel.svg
        FILE window.svg
    FILE README.md
    DIR src/
        DIR app/
            FILE favicon.ico
            FILE globals.css
            FILE layout.tsx
            FILE page.tsx
        DIR components/
            DIR ui/
                FILE button.tsx
                FILE form.tsx
                FILE label.tsx
                FILE table.tsx
        DIR lib/
            FILE utils.ts
    FILE tsconfig.json
DIR ml/
    FILE README.md
FILE project.md
FILE README.md
### DIRECTORY . FOLDER STRUCTURE ###

### DIRECTORY . FLATTENED CONTENT ###
### .\backend\bot\app.log BEGIN ###
{"level":"info","msg":"The logging system has been successfully initialized.","time":"2025-12-06T02:09:18.241+03:00"}
{"level":"info","msg":"Config has been successfully loaded","time":"2025-12-06T02:09:18.241+03:00"}
{"level":"info","msg":"Application has been successfully initialized","time":"2025-12-06T02:09:22.154+03:00"}
{"level":"info","msg":"Starting the gRPC server...","time":"2025-12-06T02:09:22.155+03:00"}
{"level":"info","msg":"Setting webhookmap[url:https://all-in-stardust.ru/bot/webhook]","time":"2025-12-06T02:09:22.155+03:00"}
{"level":"info","msg":"grpc server is running[::]:2011","time":"2025-12-06T02:09:22.155+03:00"}
{"level":"info","msg":"Starting bot","time":"2025-12-06T02:09:22.318+03:00"}
{"level":"info","msg":"Starting Bot Webhook Server on portmap[port:8081]","time":"2025-12-06T02:09:22.319+03:00"}
{"level":"info","msg":"Shutting down the application","time":"2025-12-06T02:10:20.946+03:00"}
{"level":"info","msg":"Stopping Bot Webhook Server.","time":"2025-12-06T02:10:20.946+03:00"}
{"level":"info","msg":"Stopping the gRPC server...","time":"2025-12-06T02:10:20.946+03:00"}
{"level":"info","msg":"Application has been successfully stopped","time":"2025-12-06T02:10:20.946+03:00"}

### .\backend\bot\app.log END ###

### .\backend\bot\cmd\main.go BEGIN ###
package main

import (
	"bot/internal/app"
	"bot/internal/config"
	"context"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/PrototypeSirius/ruglogger/apperror"
	"github.com/PrototypeSirius/ruglogger/logger"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

var level string

func main() {
	cfg := config.MustLoad()
	log, logCloser := logInit(cfg.Env)
	defer func() {
		if err := logCloser(); err != nil {
			os.Stderr.WriteString("Error closing log file: " + err.Error() + "\n")
		}
	}()
	log.Info("Config has been successfully loaded")
	application := app.New(log, cfg)
	log.Info("Application has been successfully initialized")
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	g, gCtx := errgroup.WithContext(ctx)
	g.Go(func() error {
		return application.BotServer.Run(gCtx)
	})
	g.Go(func() error {
		return application.GRPCServer.Run()
	})
	g.Go(func() error {
		<-gCtx.Done()
		logger.Info("Shutting down the application")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := application.BotServer.Stop(shutdownCtx); err != nil {
			appErr := apperror.SystemError(err, 1001, "error stopping bot server")
			logger.LogOnError(appErr, "stopping bot server")
		}
		application.GRPCServer.Stop()
		return nil
	})
	if err := g.Wait(); err != nil {
		appErr := apperror.SystemError(err, 1002, "error running application")
		logger.LogOnError(appErr, "running application")
	} else {
		log.Info("Application has been successfully stopped")
	}
}

func logInit(env string) (*logrus.Logger, func() error) {
	switch env {
	case "production":
		level = "warn"
	case "local":
		level = "info"
	default:
		level = "debug"
	}
	logFile, err := os.OpenFile("app.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Critical error: failed to open log file: %v\n", err)
	}
	multiWriter := io.MultiWriter(os.Stdout, logFile)
	logger.Init(
		logger.WithLevel(level),
		logger.WithOutput(multiWriter),
	)
	log := logger.Get()
	log.Info("The logging system has been successfully initialized.")
	return log, logFile.Close
}

### .\backend\bot\cmd\main.go END ###

### .\backend\bot\config\config.yaml BEGIN ###
env: "local" # local, production
grpc:
  port: 2011
  timeout: 20s
bot: 
  token: "8011424482:AAEQd6xZoJH9gCf2oJfmh6z1-2qWZnWp0_s"
  tradeToken: "CtTeKffFRJhhGUttyqaStT"
  webhookToken: "lqxvHTdJXgthPEuYUtSJGFXoGOXKABsibGlvNyBcGHp"
  webURL: "https://all-in-stardust.ru/"
  port: 8081
### .\backend\bot\config\config.yaml END ###

### .\backend\bot\Dockerfile BEGIN ###
FROM golang:alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o bot-app cmd/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

RUN addgroup -S appgroup && adduser -S appuser -G appgroup


COPY --from=builder /app/bot-app .

COPY --from=builder /app/config/config.yaml ./config/config.yaml

RUN touch app.log && chown appuser:appgroup app.log

USER appuser

EXPOSE 8081
EXPOSE 2011

CMD ["./bot-app", "--config=./config/config.yaml"]
### .\backend\bot\Dockerfile END ###

### .\backend\bot\go.mod BEGIN ###
module bot

go 1.24.2

require (
	github.com/PrototypeSirius/protos_service v0.0.0-20251206133529-ad042aab1d47
	github.com/PrototypeSirius/ruglogger v0.0.0-20251116235218-0aa354b7deb2
	github.com/go-telegram/bot v1.17.0
	github.com/ilyakaznacheev/cleanenv v1.5.0
	github.com/sirupsen/logrus v1.9.3
	github.com/telegram-mini-apps/init-data-golang v1.5.0
	golang.org/x/sync v0.17.0
	google.golang.org/grpc v1.77.0
)

require (
	github.com/BurntSushi/toml v1.2.1 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	golang.org/x/net v0.46.1-0.20251013234738-63d1a5100f82 // indirect
	golang.org/x/sys v0.37.0 // indirect
	golang.org/x/text v0.30.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251022142026-3a174f9686a8 // indirect
	google.golang.org/protobuf v1.36.10 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	olympos.io/encoding/edn v0.0.0-20201019073823-d3554ca0b0a3 // indirect
)

### .\backend\bot\go.mod END ###

### .\backend\bot\go.sum BEGIN ###
github.com/BurntSushi/toml v1.2.1 h1:9F2/+DoOYIOksmaJFPw1tGFy1eDnIJXg+UHjuD8lTak=
github.com/BurntSushi/toml v1.2.1/go.mod h1:CxXYINrC8qIiEnFrOxCa7Jy5BFHlXnUU2pbicEuybxQ=
github.com/PrototypeSirius/protos_service v0.0.0-20251205214121-8599937f9392 h1:/aA/TShz3rDAF821FJyYUZ0n9wjXFHacYZ+mjLpefRg=
github.com/PrototypeSirius/protos_service v0.0.0-20251205214121-8599937f9392/go.mod h1:hOJuOHkvh867yXdbopMLGrc0gRiSI+ENZWRBqSZSiNw=
github.com/PrototypeSirius/protos_service v0.0.0-20251206133529-ad042aab1d47 h1:slbJBbdCL1XnGLqam6WMGELX1VmilIE26xWRCXB1BzE=
github.com/PrototypeSirius/protos_service v0.0.0-20251206133529-ad042aab1d47/go.mod h1:hOJuOHkvh867yXdbopMLGrc0gRiSI+ENZWRBqSZSiNw=
github.com/PrototypeSirius/ruglogger v0.0.0-20251116235218-0aa354b7deb2 h1:67qBAnVo8NJxag8sBFZNsobQEWHhr0np3QhXfRAl7R8=
github.com/PrototypeSirius/ruglogger v0.0.0-20251116235218-0aa354b7deb2/go.mod h1:3KeQKrOLvDy6xVl90WmKDW3VyQwtWvvvQb3wWkR059s=
github.com/davecgh/go-spew v1.1.0/go.mod h1:J7Y8YcW2NihsgmVo/mv3lAwl/skON4iLHjSsI+c5H38=
github.com/davecgh/go-spew v1.1.1 h1:vj9j/u1bqnvCEfJOwUhtlOARqs3+rkHYY13jYWTU97c=
github.com/davecgh/go-spew v1.1.1/go.mod h1:J7Y8YcW2NihsgmVo/mv3lAwl/skON4iLHjSsI+c5H38=
github.com/go-logr/logr v1.4.3 h1:CjnDlHq8ikf6E492q6eKboGOC0T8CDaOvkHCIg8idEI=
github.com/go-logr/logr v1.4.3/go.mod h1:9T104GzyrTigFIr8wt5mBrctHMim0Nb2HLGrmQ40KvY=
github.com/go-logr/stdr v1.2.2 h1:hSWxHoqTgW2S2qGc0LTAI563KZ5YKYRhT3MFKZMbjag=
github.com/go-logr/stdr v1.2.2/go.mod h1:mMo/vtBO5dYbehREoey6XUKy/eSumjCCveDpRre4VKE=
github.com/go-telegram/bot v1.17.0 h1:Hs0kGxSj97QFqOQP0zxduY/4tSx8QDzvNI9uVRS+zmY=
github.com/go-telegram/bot v1.17.0/go.mod h1:i2TRs7fXWIeaceF3z7KzsMt/he0TwkVC680mvdTFYeM=
github.com/golang/protobuf v1.5.4 h1:i7eJL8qZTpSEXOPTxNKhASYpMn+8e5Q6AdndVa1dWek=
github.com/golang/protobuf v1.5.4/go.mod h1:lnTiLA8Wa4RWRcIUkrtSVa5nRhsEGBg48fD6rSs7xps=
github.com/google/go-cmp v0.7.0 h1:wk8382ETsv4JYUZwIsn6YpYiWiBsYLSJiTsyBybVuN8=
github.com/google/go-cmp v0.7.0/go.mod h1:pXiqmnSA92OHEEa9HXL2W4E7lf9JzCmGVUdgjX3N/iU=
github.com/google/uuid v1.6.0 h1:NIvaJDMOsjHA8n1jAhLSgzrAzy1Hgr+hNrb57e+94F0=
github.com/google/uuid v1.6.0/go.mod h1:TIyPZe4MgqvfeYDBFedMoGGpEw/LqOeaOT+nhxU+yHo=
github.com/ilyakaznacheev/cleanenv v1.5.0 h1:0VNZXggJE2OYdXE87bfSSwGxeiGt9moSR2lOrsHHvr4=
github.com/ilyakaznacheev/cleanenv v1.5.0/go.mod h1:a5aDzaJrLCQZsazHol1w8InnDcOX0OColm64SlIi6gk=
github.com/joho/godotenv v1.5.1 h1:7eLL/+HRGLY0ldzfGMeQkb7vMd0as4CfYvUVzLqw0N0=
github.com/joho/godotenv v1.5.1/go.mod h1:f4LDr5Voq0i2e/R5DDNOoa2zzDfwtkZa6DnEwAbqwq4=
github.com/pmezard/go-difflib v1.0.0 h1:4DBwDE0NGyQoBHbLQYPwSUPoCMWR5BEzIk/f1lZbAQM=
github.com/pmezard/go-difflib v1.0.0/go.mod h1:iKH77koFhYxTK1pcRnkKkqfTogsbg7gZNVY4sRDYZ/4=
github.com/sirupsen/logrus v1.9.3 h1:dueUQJ1C2q9oE3F7wvmSGAaVtTmUizReu6fjN8uqzbQ=
github.com/sirupsen/logrus v1.9.3/go.mod h1:naHLuLoDiP4jHNo9R0sCBMtWGeIprob74mVsIT4qYEQ=
github.com/stretchr/objx v0.1.0/go.mod h1:HFkY916IF+rwdDfMAkV7OtwuqBVzrE8GR6GFx+wExME=
github.com/stretchr/testify v1.7.0/go.mod h1:6Fq8oRcR53rry900zMqJjRRixrwX3KX962/h/Wwjteg=
github.com/stretchr/testify v1.11.1 h1:7s2iGBzp5EwR7/aIZr8ao5+dra3wiQyKjjFuvgVKu7U=
github.com/stretchr/testify v1.11.1/go.mod h1:wZwfW3scLgRK+23gO65QZefKpKQRnfz6sD981Nm4B6U=
github.com/telegram-mini-apps/init-data-golang v1.5.0 h1:rtpsmQ/nihkicPvnrdRXmHHtTnPvG1FmxMRZJwMKPz0=
github.com/telegram-mini-apps/init-data-golang v1.5.0/go.mod h1:GG4HnRx9ocjD4MjjzOw7gf9Ptm0NvFbDr5xqnfFOYuY=
go.opentelemetry.io/auto/sdk v1.2.1 h1:jXsnJ4Lmnqd11kwkBV2LgLoFMZKizbCi5fNZ/ipaZ64=
go.opentelemetry.io/auto/sdk v1.2.1/go.mod h1:KRTj+aOaElaLi+wW1kO/DZRXwkF4C5xPbEe3ZiIhN7Y=
go.opentelemetry.io/otel v1.38.0 h1:RkfdswUDRimDg0m2Az18RKOsnI8UDzppJAtj01/Ymk8=
go.opentelemetry.io/otel v1.38.0/go.mod h1:zcmtmQ1+YmQM9wrNsTGV/q/uyusom3P8RxwExxkZhjM=
go.opentelemetry.io/otel/metric v1.38.0 h1:Kl6lzIYGAh5M159u9NgiRkmoMKjvbsKtYRwgfrA6WpA=
go.opentelemetry.io/otel/metric v1.38.0/go.mod h1:kB5n/QoRM8YwmUahxvI3bO34eVtQf2i4utNVLr9gEmI=
go.opentelemetry.io/otel/sdk v1.38.0 h1:l48sr5YbNf2hpCUj/FoGhW9yDkl+Ma+LrVl8qaM5b+E=
go.opentelemetry.io/otel/sdk v1.38.0/go.mod h1:ghmNdGlVemJI3+ZB5iDEuk4bWA3GkTpW+DOoZMYBVVg=
go.opentelemetry.io/otel/sdk/metric v1.38.0 h1:aSH66iL0aZqo//xXzQLYozmWrXxyFkBJ6qT5wthqPoM=
go.opentelemetry.io/otel/sdk/metric v1.38.0/go.mod h1:dg9PBnW9XdQ1Hd6ZnRz689CbtrUp0wMMs9iPcgT9EZA=
go.opentelemetry.io/otel/trace v1.38.0 h1:Fxk5bKrDZJUH+AMyyIXGcFAPah0oRcT+LuNtJrmcNLE=
go.opentelemetry.io/otel/trace v1.38.0/go.mod h1:j1P9ivuFsTceSWe1oY+EeW3sc+Pp42sO++GHkg4wwhs=
golang.org/x/net v0.46.1-0.20251013234738-63d1a5100f82 h1:6/3JGEh1C88g7m+qzzTbl3A0FtsLguXieqofVLU/JAo=
golang.org/x/net v0.46.1-0.20251013234738-63d1a5100f82/go.mod h1:Q9BGdFy1y4nkUwiLvT5qtyhAnEHgnQ/zd8PfU6nc210=
golang.org/x/sync v0.17.0 h1:l60nONMj9l5drqw6jlhIELNv9I0A4OFgRsG9k2oT9Ug=
golang.org/x/sync v0.17.0/go.mod h1:9KTHXmSnoGruLpwFjVSX0lNNA75CykiMECbovNTZqGI=
golang.org/x/sys v0.0.0-20220715151400-c0bba94af5f8/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
golang.org/x/sys v0.37.0 h1:fdNQudmxPjkdUTPnLn5mdQv7Zwvbvpaxqs831goi9kQ=
golang.org/x/sys v0.37.0/go.mod h1:OgkHotnGiDImocRcuBABYBEXf8A9a87e/uXjp9XT3ks=
golang.org/x/text v0.30.0 h1:yznKA/E9zq54KzlzBEAWn1NXSQ8DIp/NYMy88xJjl4k=
golang.org/x/text v0.30.0/go.mod h1:yDdHFIX9t+tORqspjENWgzaCVXgk0yYnYuSZ8UzzBVM=
gonum.org/v1/gonum v0.16.0 h1:5+ul4Swaf3ESvrOnidPp4GZbzf0mxVQpDCYUQE7OJfk=
gonum.org/v1/gonum v0.16.0/go.mod h1:fef3am4MQ93R2HHpKnLk4/Tbh/s0+wqD5nfa6Pnwy4E=
google.golang.org/genproto/googleapis/rpc v0.0.0-20251022142026-3a174f9686a8 h1:M1rk8KBnUsBDg1oPGHNCxG4vc1f49epmTO7xscSajMk=
google.golang.org/genproto/googleapis/rpc v0.0.0-20251022142026-3a174f9686a8/go.mod h1:7i2o+ce6H/6BluujYR+kqX3GKH+dChPTQU19wjRPiGk=
google.golang.org/grpc v1.77.0 h1:wVVY6/8cGA6vvffn+wWK5ToddbgdU3d8MNENr4evgXM=
google.golang.org/grpc v1.77.0/go.mod h1:z0BY1iVj0q8E1uSQCjL9cppRj+gnZjzDnzV0dHhrNig=
google.golang.org/protobuf v1.36.10 h1:AYd7cD/uASjIL6Q9LiTjz8JLcrh/88q5UObnmY3aOOE=
google.golang.org/protobuf v1.36.10/go.mod h1:HTf+CrKn2C3g5S8VImy6tdcUvCska2kB7j23XfzDpco=
gopkg.in/check.v1 v0.0.0-20161208181325-20d25e280405 h1:yhCVgyC4o1eVCa2tZl7eS0r+SDo693bJlVdllGtEeKM=
gopkg.in/check.v1 v0.0.0-20161208181325-20d25e280405/go.mod h1:Co6ibVJAznAaIkqp8huTwlJQCZ016jof/cbN4VW5Yz0=
gopkg.in/yaml.v3 v3.0.0-20200313102051-9f266ea9e77c/go.mod h1:K4uyk7z7BCEPqu6E+C64Yfv1cQ7kz7rIZviUmN+EgEM=
gopkg.in/yaml.v3 v3.0.1 h1:fxVm/GzAzEWqLHuvctI91KS9hhNmmWOoWu0XTYJS7CA=
gopkg.in/yaml.v3 v3.0.1/go.mod h1:K4uyk7z7BCEPqu6E+C64Yfv1cQ7kz7rIZviUmN+EgEM=
olympos.io/encoding/edn v0.0.0-20201019073823-d3554ca0b0a3 h1:slmdOY3vp8a7KQbHkL+FLbvbkgMqmXojpFUO/jENuqQ=
olympos.io/encoding/edn v0.0.0-20201019073823-d3554ca0b0a3/go.mod h1:oVgVk4OWVDi43qWBEyGhXgYxt7+ED4iYNpTngSLX2Iw=

### .\backend\bot\go.sum END ###

### .\backend\bot\internal\app\app.go BEGIN ###
package app

import (
	botapp "bot/internal/app/bot"
	grpcapp "bot/internal/app/grpc"
	"bot/internal/config"
	botrepo "bot/internal/repository/bot"
	botservice "bot/internal/services"

	"github.com/sirupsen/logrus"
)

type App struct {
	GRPCServer *grpcapp.App
	BotServer  *botapp.App
}

func New(log *logrus.Logger, cfg *config.Config) *App {
	botServer, bot := botapp.New(log, cfg.Bot)
	botRepo := botrepo.New(bot, cfg.Bot.Token)
	botService := botservice.New(log, botRepo)
	grpcServer := grpcapp.New(log, botService, cfg.GRPC.Port)
	return &App{GRPCServer: grpcServer, BotServer: botServer}
}

### .\backend\bot\internal\app\app.go END ###

### .\backend\bot\internal\app\bot\app.go BEGIN ###
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

### .\backend\bot\internal\app\bot\app.go END ###

### .\backend\bot\internal\app\grpc\app.go BEGIN ###
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

### .\backend\bot\internal\app\grpc\app.go END ###

### .\backend\bot\internal\config\config.go BEGIN ###
package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env  string     `yaml:"env" env-default:"local"` // local, production
	GRPC GRPCConfig `yaml:"grpc"`                    // gRPC config
	Bot  BotConfig  `yaml:"bot" env-required:"true"` // bot config
}

type GRPCConfig struct {
	Port    int           `yaml:"port"`    // gRPC port
	Timeout time.Duration `yaml:"timeout"` // gRPC timeout
}

type BotConfig struct {
	Token        string `yaml:"token" env-required:"true"`        // bot token
	WebhookToken string `yaml:"webhookToken" env-required:"true"` // webhook token
	WebURL       string `yaml:"webURL" env-required:"true"`       // webhook url
	TradeToken   string `yaml:"tradeToken" env-required:"true"`   // trade token
	Port         int    `yaml:"port" env-required:"true"`         // bot port
}

func MustLoad() *Config {
	path := fechPathConfig()
	if path == "" {
		panic("config path is empty")
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file not found: " + path)
	}
	var cfg Config
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("failed to load config: " + err.Error())
	}
	return &cfg
}

func fechPathConfig() string {
	var res string
	//--config="path/to/config.yaml"
	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()
	return res
}

### .\backend\bot\internal\config\config.go END ###

### .\backend\bot\internal\handlers\bot\server.go BEGIN ###
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
		kb := &models.InlineKeyboardMarkup{InlineKeyboard: [][]models.InlineKeyboardButton{{{Text: "MiniApp", WebApp: &models.WebAppInfo{URL: cfg.WebURL}}}}}
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      update.Message.Chat.ID,
			Text:        "GAAAMBLING!",
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
			Text:   "What the actual fuck?",
		})
	}
}

### .\backend\bot\internal\handlers\bot\server.go END ###

### .\backend\bot\internal\handlers\grpc\server.go BEGIN ###
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

### .\backend\bot\internal\handlers\grpc\server.go END ###

### .\backend\bot\internal\models\model.go BEGIN ###
package models

type UserResponse struct {
	Authorized bool
	UserID     int64
	Error      string
}

### .\backend\bot\internal\models\model.go END ###

### .\backend\bot\internal\repository\bot\bot.go BEGIN ###
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

### .\backend\bot\internal\repository\bot\bot.go END ###

### .\backend\bot\internal\services\bot.go BEGIN ###
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

### .\backend\bot\internal\services\bot.go END ###

### .\backend\database\cmd\main.go BEGIN ###
package main

import (
	"database/internal/app"
	"database/internal/config"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/PrototypeSirius/ruglogger/logger"
	"github.com/sirupsen/logrus"
)

var level string

func main() {
	cfg := config.MustLoad()
	log, logCloser := logInit(cfg.Env)
	defer func() {
		if err := logCloser(); err != nil {
			os.Stderr.WriteString("Error closing log file: " + err.Error() + "\n")
		}
	}()
	log.Info("Config has been successfully loaded")

	application, closeDB := app.New(log, cfg.Database, cfg.GRPC.Port)
	log.Info("Application has been successfully initialized.")
	go func() {
		application.GRPCServer.MustRun()
	}()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop
	application.GRPCServer.Stop()
	closeDB()
	log.Info("Gracefully stopped")
}

func logInit(env string) (*logrus.Logger, func() error) {
	switch env {
	case "production":
		level = "warn"
	case "local":
		level = "info"
	default:
		level = "debug"
	}
	logFile, err := os.OpenFile("app.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Critical error: failed to open log file: %v\n", err)
	}
	multiWriter := io.MultiWriter(os.Stdout, logFile)
	logger.Init(
		logger.WithLevel(level),
		logger.WithOutput(multiWriter),
	)
	log := logger.Get()
	log.Info("The logging system has been successfully initialized.")
	return log, logFile.Close
}

### .\backend\database\cmd\main.go END ###

### .\backend\database\config\config.yaml BEGIN ###
env: "local" # local, production
grpc:
  port: 2012
  timeout: 20s

postgres:
  host: "db"
  port: 5432
  user: "cybergardenadmin"
  password: "cybergardeninstall"
  database: "cybergardendata"
  migrations_path: "./migrations"
### .\backend\database\config\config.yaml END ###

### .\backend\database\Dockerfile BEGIN ###
FROM golang:alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o database-app cmd/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

RUN addgroup -S appgroup && adduser -S appuser -G appgroup


COPY --from=builder /app/database-app .

COPY --from=builder /app/config/config.yaml ./config/config.yaml

RUN touch app.log && chown appuser:appgroup app.log

USER appuser

EXPOSE 2012

CMD ["./database-app", "--config=./config/config.yaml"]
### .\backend\database\Dockerfile END ###

### .\backend\database\go.mod BEGIN ###
module database

go 1.24.2

require (
	github.com/PrototypeSirius/protos_service v0.0.0-20251206133529-ad042aab1d47
	github.com/PrototypeSirius/ruglogger v0.0.0-20251116235218-0aa354b7deb2
	github.com/golang-migrate/migrate/v4 v4.19.1
	github.com/ilyakaznacheev/cleanenv v1.5.0
	github.com/lib/pq v1.10.9
	github.com/sirupsen/logrus v1.9.3
	google.golang.org/grpc v1.77.0
)

require (
	github.com/BurntSushi/toml v1.2.1 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	golang.org/x/net v0.47.0 // indirect
	golang.org/x/sys v0.38.0 // indirect
	golang.org/x/text v0.31.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251022142026-3a174f9686a8 // indirect
	google.golang.org/protobuf v1.36.10 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	olympos.io/encoding/edn v0.0.0-20201019073823-d3554ca0b0a3 // indirect
)

### .\backend\database\go.mod END ###

### .\backend\database\go.sum BEGIN ###
github.com/Azure/go-ansiterm v0.0.0-20230124172434-306776ec8161 h1:L/gRVlceqvL25UVaW/CKtUDjefjrs0SPonmDGUVOYP0=
github.com/Azure/go-ansiterm v0.0.0-20230124172434-306776ec8161/go.mod h1:xomTg63KZ2rFqZQzSB4Vz2SUXa1BpHTVz9L5PTmPC4E=
github.com/BurntSushi/toml v1.2.1 h1:9F2/+DoOYIOksmaJFPw1tGFy1eDnIJXg+UHjuD8lTak=
github.com/BurntSushi/toml v1.2.1/go.mod h1:CxXYINrC8qIiEnFrOxCa7Jy5BFHlXnUU2pbicEuybxQ=
github.com/Microsoft/go-winio v0.6.2 h1:F2VQgta7ecxGYO8k3ZZz3RS8fVIXVxONVUPlNERoyfY=
github.com/Microsoft/go-winio v0.6.2/go.mod h1:yd8OoFMLzJbo9gZq8j5qaps8bJ9aShtEA8Ipt1oGCvU=
github.com/PrototypeSirius/protos_service v0.0.0-20251205214121-8599937f9392 h1:/aA/TShz3rDAF821FJyYUZ0n9wjXFHacYZ+mjLpefRg=
github.com/PrototypeSirius/protos_service v0.0.0-20251205214121-8599937f9392/go.mod h1:hOJuOHkvh867yXdbopMLGrc0gRiSI+ENZWRBqSZSiNw=
github.com/PrototypeSirius/protos_service v0.0.0-20251206132036-afab1f88818a h1:Zho5MURnYDCqxHd/y6iAFSr6sgNbW1lhkVAvi28dXqw=
github.com/PrototypeSirius/protos_service v0.0.0-20251206132036-afab1f88818a/go.mod h1:hOJuOHkvh867yXdbopMLGrc0gRiSI+ENZWRBqSZSiNw=
github.com/PrototypeSirius/protos_service v0.0.0-20251206133211-19ab78708444 h1:+M57f4ruXLEoCTqQm9xEg2jYRfoZ9smzXT6mRvFKpOE=
github.com/PrototypeSirius/protos_service v0.0.0-20251206133211-19ab78708444/go.mod h1:hOJuOHkvh867yXdbopMLGrc0gRiSI+ENZWRBqSZSiNw=
github.com/PrototypeSirius/protos_service v0.0.0-20251206133529-ad042aab1d47 h1:slbJBbdCL1XnGLqam6WMGELX1VmilIE26xWRCXB1BzE=
github.com/PrototypeSirius/protos_service v0.0.0-20251206133529-ad042aab1d47/go.mod h1:hOJuOHkvh867yXdbopMLGrc0gRiSI+ENZWRBqSZSiNw=
github.com/PrototypeSirius/ruglogger v0.0.0-20251116235218-0aa354b7deb2 h1:67qBAnVo8NJxag8sBFZNsobQEWHhr0np3QhXfRAl7R8=
github.com/PrototypeSirius/ruglogger v0.0.0-20251116235218-0aa354b7deb2/go.mod h1:3KeQKrOLvDy6xVl90WmKDW3VyQwtWvvvQb3wWkR059s=
github.com/containerd/errdefs v1.0.0 h1:tg5yIfIlQIrxYtu9ajqY42W3lpS19XqdxRQeEwYG8PI=
github.com/containerd/errdefs v1.0.0/go.mod h1:+YBYIdtsnF4Iw6nWZhJcqGSg/dwvV7tyJ/kCkyJ2k+M=
github.com/containerd/errdefs/pkg v0.3.0 h1:9IKJ06FvyNlexW690DXuQNx2KA2cUJXx151Xdx3ZPPE=
github.com/containerd/errdefs/pkg v0.3.0/go.mod h1:NJw6s9HwNuRhnjJhM7pylWwMyAkmCQvQ4GpJHEqRLVk=
github.com/davecgh/go-spew v1.1.0/go.mod h1:J7Y8YcW2NihsgmVo/mv3lAwl/skON4iLHjSsI+c5H38=
github.com/davecgh/go-spew v1.1.1/go.mod h1:J7Y8YcW2NihsgmVo/mv3lAwl/skON4iLHjSsI+c5H38=
github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc h1:U9qPSI2PIWSS1VwoXQT9A3Wy9MM3WgvqSxFWenqJduM=
github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc/go.mod h1:J7Y8YcW2NihsgmVo/mv3lAwl/skON4iLHjSsI+c5H38=
github.com/dhui/dktest v0.4.6 h1:+DPKyScKSEp3VLtbMDHcUq6V5Lm5zfZZVb0Sk7Ahom4=
github.com/dhui/dktest v0.4.6/go.mod h1:JHTSYDtKkvFNFHJKqCzVzqXecyv+tKt8EzceOmQOgbU=
github.com/distribution/reference v0.6.0 h1:0IXCQ5g4/QMHHkarYzh5l+u8T3t73zM5QvfrDyIgxBk=
github.com/distribution/reference v0.6.0/go.mod h1:BbU0aIcezP1/5jX/8MP0YiH4SdvB5Y4f/wlDRiLyi3E=
github.com/docker/docker v28.3.3+incompatible h1:Dypm25kh4rmk49v1eiVbsAtpAsYURjYkaKubwuBdxEI=
github.com/docker/docker v28.3.3+incompatible/go.mod h1:eEKB0N0r5NX/I1kEveEz05bcu8tLC/8azJZsviup8Sk=
github.com/docker/go-connections v0.5.0 h1:USnMq7hx7gwdVZq1L49hLXaFtUdTADjXGp+uj1Br63c=
github.com/docker/go-connections v0.5.0/go.mod h1:ov60Kzw0kKElRwhNs9UlUHAE/F9Fe6GLaXnqyDdmEXc=
github.com/docker/go-units v0.5.0 h1:69rxXcBk27SvSaaxTtLh/8llcHD8vYHT7WSdRZ/jvr4=
github.com/docker/go-units v0.5.0/go.mod h1:fgPhTUdO+D/Jk86RDLlptpiXQzgHJF7gydDDbaIK4Dk=
github.com/felixge/httpsnoop v1.0.4 h1:NFTV2Zj1bL4mc9sqWACXbQFVBBg2W3GPvqp8/ESS2Wg=
github.com/felixge/httpsnoop v1.0.4/go.mod h1:m8KPJKqk1gH5J9DgRY2ASl2lWCfGKXixSwevea8zH2U=
github.com/go-logr/logr v1.4.3 h1:CjnDlHq8ikf6E492q6eKboGOC0T8CDaOvkHCIg8idEI=
github.com/go-logr/logr v1.4.3/go.mod h1:9T104GzyrTigFIr8wt5mBrctHMim0Nb2HLGrmQ40KvY=
github.com/go-logr/stdr v1.2.2 h1:hSWxHoqTgW2S2qGc0LTAI563KZ5YKYRhT3MFKZMbjag=
github.com/go-logr/stdr v1.2.2/go.mod h1:mMo/vtBO5dYbehREoey6XUKy/eSumjCCveDpRre4VKE=
github.com/gogo/protobuf v1.3.2 h1:Ov1cvc58UF3b5XjBnZv7+opcTcQFZebYjWzi34vdm4Q=
github.com/gogo/protobuf v1.3.2/go.mod h1:P1XiOD3dCwIKUDQYPy72D8LYyHL2YPYrpS2s69NZV8Q=
github.com/golang-migrate/migrate/v4 v4.19.1 h1:OCyb44lFuQfYXYLx1SCxPZQGU7mcaZ7gH9yH4jSFbBA=
github.com/golang-migrate/migrate/v4 v4.19.1/go.mod h1:CTcgfjxhaUtsLipnLoQRWCrjYXycRz/g5+RWDuYgPrE=
github.com/golang/protobuf v1.5.4 h1:i7eJL8qZTpSEXOPTxNKhASYpMn+8e5Q6AdndVa1dWek=
github.com/golang/protobuf v1.5.4/go.mod h1:lnTiLA8Wa4RWRcIUkrtSVa5nRhsEGBg48fD6rSs7xps=
github.com/google/go-cmp v0.7.0 h1:wk8382ETsv4JYUZwIsn6YpYiWiBsYLSJiTsyBybVuN8=
github.com/google/go-cmp v0.7.0/go.mod h1:pXiqmnSA92OHEEa9HXL2W4E7lf9JzCmGVUdgjX3N/iU=
github.com/google/uuid v1.6.0 h1:NIvaJDMOsjHA8n1jAhLSgzrAzy1Hgr+hNrb57e+94F0=
github.com/google/uuid v1.6.0/go.mod h1:TIyPZe4MgqvfeYDBFedMoGGpEw/LqOeaOT+nhxU+yHo=
github.com/ilyakaznacheev/cleanenv v1.5.0 h1:0VNZXggJE2OYdXE87bfSSwGxeiGt9moSR2lOrsHHvr4=
github.com/ilyakaznacheev/cleanenv v1.5.0/go.mod h1:a5aDzaJrLCQZsazHol1w8InnDcOX0OColm64SlIi6gk=
github.com/joho/godotenv v1.5.1 h1:7eLL/+HRGLY0ldzfGMeQkb7vMd0as4CfYvUVzLqw0N0=
github.com/joho/godotenv v1.5.1/go.mod h1:f4LDr5Voq0i2e/R5DDNOoa2zzDfwtkZa6DnEwAbqwq4=
github.com/lib/pq v1.10.9 h1:YXG7RB+JIjhP29X+OtkiDnYaXQwpS4JEWq7dtCCRUEw=
github.com/lib/pq v1.10.9/go.mod h1:AlVN5x4E4T544tWzH6hKfbfQvm3HdbOxrmggDNAPY9o=
github.com/moby/docker-image-spec v1.3.1 h1:jMKff3w6PgbfSa69GfNg+zN/XLhfXJGnEx3Nl2EsFP0=
github.com/moby/docker-image-spec v1.3.1/go.mod h1:eKmb5VW8vQEh/BAr2yvVNvuiJuY6UIocYsFu/DxxRpo=
github.com/moby/term v0.5.0 h1:xt8Q1nalod/v7BqbG21f8mQPqH+xAaC9C3N3wfWbVP0=
github.com/moby/term v0.5.0/go.mod h1:8FzsFHVUBGZdbDsJw/ot+X+d5HLUbvklYLJ9uGfcI3Y=
github.com/morikuni/aec v1.0.0 h1:nP9CBfwrvYnBRgY6qfDQkygYDmYwOilePFkwzv4dU8A=
github.com/morikuni/aec v1.0.0/go.mod h1:BbKIizmSmc5MMPqRYbxO4ZU0S0+P200+tUnFx7PXmsc=
github.com/opencontainers/go-digest v1.0.0 h1:apOUWs51W5PlhuyGyz9FCeeBIOUDA/6nW8Oi/yOhh5U=
github.com/opencontainers/go-digest v1.0.0/go.mod h1:0JzlMkj0TRzQZfJkVvzbP0HBR3IKzErnv2BNG4W4MAM=
github.com/opencontainers/image-spec v1.1.0 h1:8SG7/vwALn54lVB/0yZ/MMwhFrPYtpEHQb2IpWsCzug=
github.com/opencontainers/image-spec v1.1.0/go.mod h1:W4s4sFTMaBeK1BQLXbG4AdM2szdn85PY75RI83NrTrM=
github.com/pkg/errors v0.9.1 h1:FEBLx1zS214owpjy7qsBeixbURkuhQAwrK5UwLGTwt4=
github.com/pkg/errors v0.9.1/go.mod h1:bwawxfHBFNV+L2hUp1rHADufV3IMtnDRdf1r5NINEl0=
github.com/pmezard/go-difflib v1.0.0/go.mod h1:iKH77koFhYxTK1pcRnkKkqfTogsbg7gZNVY4sRDYZ/4=
github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 h1:Jamvg5psRIccs7FGNTlIRMkT8wgtp5eCXdBlqhYGL6U=
github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2/go.mod h1:iKH77koFhYxTK1pcRnkKkqfTogsbg7gZNVY4sRDYZ/4=
github.com/sirupsen/logrus v1.9.3 h1:dueUQJ1C2q9oE3F7wvmSGAaVtTmUizReu6fjN8uqzbQ=
github.com/sirupsen/logrus v1.9.3/go.mod h1:naHLuLoDiP4jHNo9R0sCBMtWGeIprob74mVsIT4qYEQ=
github.com/stretchr/objx v0.1.0/go.mod h1:HFkY916IF+rwdDfMAkV7OtwuqBVzrE8GR6GFx+wExME=
github.com/stretchr/testify v1.7.0/go.mod h1:6Fq8oRcR53rry900zMqJjRRixrwX3KX962/h/Wwjteg=
github.com/stretchr/testify v1.11.1 h1:7s2iGBzp5EwR7/aIZr8ao5+dra3wiQyKjjFuvgVKu7U=
github.com/stretchr/testify v1.11.1/go.mod h1:wZwfW3scLgRK+23gO65QZefKpKQRnfz6sD981Nm4B6U=
go.opentelemetry.io/auto/sdk v1.2.1 h1:jXsnJ4Lmnqd11kwkBV2LgLoFMZKizbCi5fNZ/ipaZ64=
go.opentelemetry.io/auto/sdk v1.2.1/go.mod h1:KRTj+aOaElaLi+wW1kO/DZRXwkF4C5xPbEe3ZiIhN7Y=
go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.61.0 h1:F7Jx+6hwnZ41NSFTO5q4LYDtJRXBf2PD0rNBkeB/lus=
go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.61.0/go.mod h1:UHB22Z8QsdRDrnAtX4PntOl36ajSxcdUMt1sF7Y6E7Q=
go.opentelemetry.io/otel v1.38.0 h1:RkfdswUDRimDg0m2Az18RKOsnI8UDzppJAtj01/Ymk8=
go.opentelemetry.io/otel v1.38.0/go.mod h1:zcmtmQ1+YmQM9wrNsTGV/q/uyusom3P8RxwExxkZhjM=
go.opentelemetry.io/otel/metric v1.38.0 h1:Kl6lzIYGAh5M159u9NgiRkmoMKjvbsKtYRwgfrA6WpA=
go.opentelemetry.io/otel/metric v1.38.0/go.mod h1:kB5n/QoRM8YwmUahxvI3bO34eVtQf2i4utNVLr9gEmI=
go.opentelemetry.io/otel/sdk v1.38.0 h1:l48sr5YbNf2hpCUj/FoGhW9yDkl+Ma+LrVl8qaM5b+E=
go.opentelemetry.io/otel/sdk v1.38.0/go.mod h1:ghmNdGlVemJI3+ZB5iDEuk4bWA3GkTpW+DOoZMYBVVg=
go.opentelemetry.io/otel/sdk/metric v1.38.0 h1:aSH66iL0aZqo//xXzQLYozmWrXxyFkBJ6qT5wthqPoM=
go.opentelemetry.io/otel/sdk/metric v1.38.0/go.mod h1:dg9PBnW9XdQ1Hd6ZnRz689CbtrUp0wMMs9iPcgT9EZA=
go.opentelemetry.io/otel/trace v1.38.0 h1:Fxk5bKrDZJUH+AMyyIXGcFAPah0oRcT+LuNtJrmcNLE=
go.opentelemetry.io/otel/trace v1.38.0/go.mod h1:j1P9ivuFsTceSWe1oY+EeW3sc+Pp42sO++GHkg4wwhs=
golang.org/x/net v0.47.0 h1:Mx+4dIFzqraBXUugkia1OOvlD6LemFo1ALMHjrXDOhY=
golang.org/x/net v0.47.0/go.mod h1:/jNxtkgq5yWUGYkaZGqo27cfGZ1c5Nen03aYrrKpVRU=
golang.org/x/sys v0.0.0-20220715151400-c0bba94af5f8/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
golang.org/x/sys v0.38.0 h1:3yZWxaJjBmCWXqhN1qh02AkOnCQ1poK6oF+a7xWL6Gc=
golang.org/x/sys v0.38.0/go.mod h1:OgkHotnGiDImocRcuBABYBEXf8A9a87e/uXjp9XT3ks=
golang.org/x/text v0.31.0 h1:aC8ghyu4JhP8VojJ2lEHBnochRno1sgL6nEi9WGFGMM=
golang.org/x/text v0.31.0/go.mod h1:tKRAlv61yKIjGGHX/4tP1LTbc13YSec1pxVEWXzfoeM=
gonum.org/v1/gonum v0.16.0 h1:5+ul4Swaf3ESvrOnidPp4GZbzf0mxVQpDCYUQE7OJfk=
gonum.org/v1/gonum v0.16.0/go.mod h1:fef3am4MQ93R2HHpKnLk4/Tbh/s0+wqD5nfa6Pnwy4E=
google.golang.org/genproto/googleapis/rpc v0.0.0-20251022142026-3a174f9686a8 h1:M1rk8KBnUsBDg1oPGHNCxG4vc1f49epmTO7xscSajMk=
google.golang.org/genproto/googleapis/rpc v0.0.0-20251022142026-3a174f9686a8/go.mod h1:7i2o+ce6H/6BluujYR+kqX3GKH+dChPTQU19wjRPiGk=
google.golang.org/grpc v1.77.0 h1:wVVY6/8cGA6vvffn+wWK5ToddbgdU3d8MNENr4evgXM=
google.golang.org/grpc v1.77.0/go.mod h1:z0BY1iVj0q8E1uSQCjL9cppRj+gnZjzDnzV0dHhrNig=
google.golang.org/protobuf v1.36.10 h1:AYd7cD/uASjIL6Q9LiTjz8JLcrh/88q5UObnmY3aOOE=
google.golang.org/protobuf v1.36.10/go.mod h1:HTf+CrKn2C3g5S8VImy6tdcUvCska2kB7j23XfzDpco=
gopkg.in/check.v1 v0.0.0-20161208181325-20d25e280405 h1:yhCVgyC4o1eVCa2tZl7eS0r+SDo693bJlVdllGtEeKM=
gopkg.in/check.v1 v0.0.0-20161208181325-20d25e280405/go.mod h1:Co6ibVJAznAaIkqp8huTwlJQCZ016jof/cbN4VW5Yz0=
gopkg.in/yaml.v3 v3.0.0-20200313102051-9f266ea9e77c/go.mod h1:K4uyk7z7BCEPqu6E+C64Yfv1cQ7kz7rIZviUmN+EgEM=
gopkg.in/yaml.v3 v3.0.1 h1:fxVm/GzAzEWqLHuvctI91KS9hhNmmWOoWu0XTYJS7CA=
gopkg.in/yaml.v3 v3.0.1/go.mod h1:K4uyk7z7BCEPqu6E+C64Yfv1cQ7kz7rIZviUmN+EgEM=
olympos.io/encoding/edn v0.0.0-20201019073823-d3554ca0b0a3 h1:slmdOY3vp8a7KQbHkL+FLbvbkgMqmXojpFUO/jENuqQ=
olympos.io/encoding/edn v0.0.0-20201019073823-d3554ca0b0a3/go.mod h1:oVgVk4OWVDi43qWBEyGhXgYxt7+ED4iYNpTngSLX2Iw=

### .\backend\database\go.sum END ###

### .\backend\database\internal\app\app.go BEGIN ###
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

### .\backend\database\internal\app\app.go END ###

### .\backend\database\internal\app\grpc\app.go BEGIN ###
package grpcapp

import (
	"context"
	grpchandler "database/internal/handlers/grpc"
	"fmt"
	"net"
	"time"

	"github.com/PrototypeSirius/ruglogger/apperror"
	"github.com/PrototypeSirius/ruglogger/logger"
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

func New(log *logrus.Logger, serverAPI grpchandler.DBService, port int) *App {
	opts := []grpc.ServerOption{grpc.UnaryInterceptor(loggingInterceptor(log))}
	gRPCServer := grpc.NewServer(opts...)
	grpchandler.Register(gRPCServer, serverAPI)
	return &App{log: log, grpcServer: gRPCServer, port: port}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		logger.FatalOnError(err, "Error starting gRPC server")
	}
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
	a.log.Info("Stopping the gRPC server")
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
			"duration": duration,
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
			log.WithFields(fields).Info("gRPC request succeeded")
		}
		return resp, err
	}
}

### .\backend\database\internal\app\grpc\app.go END ###

### .\backend\database\internal\config\config.go BEGIN ###
package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env      string         `yaml:"env" env-default:"local"`      // local, production
	GRPC     GRPCConfig     `yaml:"grpc"`                         // gRPC config
	Database DatabaseConfig `yaml:"postgres" env-required:"true"` // database config
}

type GRPCConfig struct {
	Port    int           `yaml:"port"`    // gRPC port
	Timeout time.Duration `yaml:"timeout"` // gRPC timeout
}

type DatabaseConfig struct {
	Host           string `yaml:"host" env-required:"true"`            // database host
	Port           int    `yaml:"port" env-required:"true"`            // database port
	User           string `yaml:"user" env-required:"true"`            // database user
	Password       string `yaml:"password" env-required:"true"`        // database password
	Database       string `yaml:"database" env-required:"true"`        // database name
	MigrationsPath string `yaml:"migrations_path" env-required:"true"` // migrations path
}

func MustLoad() *Config {
	path := fechPathConfig()
	if path == "" {
		panic("config path is empty")
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file not found: " + path)
	}
	var cfg Config
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("failed to load config: " + err.Error())
	}
	return &cfg
}

func fechPathConfig() string {
	var res string
	//--config="path/to/config.yaml"
	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()
	return res
}

### .\backend\database\internal\config\config.go END ###

### .\backend\database\internal\handlers\grpc\server.go BEGIN ###
package grpchandler

import (
	"context"
	model "database/internal/models"
	"errors"

	database "github.com/PrototypeSirius/protos_service/gen/cybergarden/database"
	"github.com/PrototypeSirius/ruglogger/apperror"
	"github.com/PrototypeSirius/ruglogger/logger"
	"google.golang.org/grpc"
)

type DBService interface {
	AddTransaction(ctx context.Context, uid int64, t model.Transaction) (string, error)
	AddUser(ctx context.Context, uid int64) (string, error)
	DeleteTransaction(ctx context.Context, uid, tid int64) (string, error)
	DeleteUser(ctx context.Context, uid int64) (string, error)
	EditTransaction(ctx context.Context, uid int64, t model.Transaction) (string, error)
	RequestUserTransactions(ctx context.Context, uid int64) ([]model.Transaction, string, error)
}

type serverAPI struct {
	database.UnimplementedDatabaseServer
	db DBService
}

func Register(s *grpc.Server, db DBService) {
	database.RegisterDatabaseServer(s, &serverAPI{db: db})
}

func (s *serverAPI) AddTransaction(ctx context.Context, req *database.AddTransactionRequest) (*database.AddTransactionResponse, error) {
	if req.GetUserID() == 0 {
		appErr := apperror.BadRequestError(errors.New("empty user id"), 1050, "Error while adding transaction")
		logger.LogOnError(appErr, "Error in adding transaction")
		return nil, appErr
	}
	if req.GetTransaction() == nil {
		appErr := apperror.BadRequestError(errors.New("empty transaction"), 1051, "Error while adding transaction")
		logger.LogOnError(appErr, "Error in adding transaction")
		return nil, appErr
	}
	mes, err := s.db.AddTransaction(ctx, req.GetUserID(), interpretatorTransactionAdd(req))
	if err != nil {
		logger.LogOnError(err, "Error in adding transaction")
		return &database.AddTransactionResponse{ErrorMes: mes}, err
	}
	return &database.AddTransactionResponse{ErrorMes: mes}, nil
}

func (s *serverAPI) AddUser(ctx context.Context, req *database.AddUserRequest) (*database.AddUserResponse, error) {
	if req.GetUserID() == 0 {
		appErr := apperror.BadRequestError(errors.New("empty user id"), 1052, "Error while adding user")
		logger.LogOnError(appErr, "Error in adding user")
		return nil, appErr
	}
	mes, err := s.db.AddUser(ctx, req.GetUserID())
	if err != nil {
		logger.LogOnError(err, "Error in adding user")
		return &database.AddUserResponse{ErrorMes: mes}, err
	}
	return &database.AddUserResponse{ErrorMes: mes}, nil
}

func (s *serverAPI) DeleteTransaction(ctx context.Context, req *database.DeleteTransactionRequest) (*database.DeleteTransactionResponse, error) {
	if req.GetUserID() == 0 {
		appErr := apperror.BadRequestError(errors.New("empty user id"), 1053, "Error while deleting transaction")
		logger.LogOnError(appErr, "Error in deleting transaction")
		return nil, appErr
	}
	if req.GetID() == 0 {
		appErr := apperror.BadRequestError(errors.New("empty transaction id"), 1054, "Error while deleting transaction")
		logger.LogOnError(appErr, "Error in deleting transaction")
		return nil, appErr
	}
	mes, err := s.db.DeleteTransaction(ctx, req.GetUserID(), req.GetID())
	if err != nil {
		logger.LogOnError(err, "Error in deleting transaction")
		return &database.DeleteTransactionResponse{ErrorMes: mes}, err
	}
	return &database.DeleteTransactionResponse{ErrorMes: mes}, nil
}

func (s *serverAPI) DeleteUser(ctx context.Context, req *database.DeleteUserRequest) (*database.DeleteUserResponse, error) {
	if req.GetUserID() == 0 {
		appErr := apperror.BadRequestError(errors.New("empty user id"), 1055, "Error while deleting user")
		logger.LogOnError(appErr, "Error in deleting user")
		return nil, appErr
	}
	mes, err := s.db.DeleteUser(ctx, req.GetUserID())
	if err != nil {
		logger.LogOnError(err, "Error in deleting user")
		return &database.DeleteUserResponse{ErrorMes: mes}, err
	}
	return &database.DeleteUserResponse{ErrorMes: mes}, nil
}

func (s *serverAPI) EditTransaction(ctx context.Context, req *database.EditTransactionRequest) (*database.EditTransactionResponse, error) {
	if req.GetUserID() == 0 {
		appErr := apperror.BadRequestError(errors.New("empty user id"), 1056, "Error while editing transaction")
		logger.LogOnError(appErr, "Error in editing transaction")
		return nil, appErr
	}
	if req.GetTransaction() == nil {
		appErr := apperror.BadRequestError(errors.New("empty transaction"), 1057, "Error while editing transaction")
		logger.LogOnError(appErr, "Error in editing transaction")
		return nil, appErr
	}
	mes, err := s.db.EditTransaction(ctx, req.GetUserID(), interpretatorTransactionEdit(req))
	if err != nil {
		logger.LogOnError(err, "Error in editing transaction")
		return &database.EditTransactionResponse{ErrorMes: mes}, err
	}
	return &database.EditTransactionResponse{}, nil
}

func (s *serverAPI) RequestUserTransactions(ctx context.Context, req *database.RequestUserTransactionsRequest) (*database.RequestUserTransactionsResponse, error) {
	if req.GetUserID() == 0 {
		appErr := apperror.BadRequestError(errors.New("empty user id"), 1058, "Error while requesting user transactions")
		logger.LogOnError(appErr, "Error in requesting user transactions")
		return nil, appErr
	}
	transactions, mes, err := s.db.RequestUserTransactions(ctx, req.GetUserID())
	if err != nil {
		logger.LogOnError(err, "Error in requesting user transactions")
		return &database.RequestUserTransactionsResponse{ErrorMes: mes}, err
	}
	return interpretatorTransactionResponse(transactions), nil
}

func interpretatorTransactionAdd(req *database.AddTransactionRequest) model.Transaction {
	t := req.GetTransaction()
	return model.Transaction{
		ID:          t.ID,
		Date:        t.Date,
		Kategoria:   t.Kategoria,
		Type:        t.Type,
		Amount:      t.Amount,
		Description: t.Description,
	}
}

func interpretatorTransactionEdit(req *database.EditTransactionRequest) model.Transaction {
	t := req.GetTransaction()
	return model.Transaction{
		ID:          t.ID,
		Date:        t.Date,
		Kategoria:   t.Kategoria,
		Type:        t.Type,
		Amount:      t.Amount,
		Description: t.Description,
	}
}

func interpretatorTransactionResponse(transactions []model.Transaction) *database.RequestUserTransactionsResponse {
	var protoTransactions []*database.Transaction
	for _, t := range transactions {
		protoTransactions = append(protoTransactions, &database.Transaction{
			ID:          t.ID,
			Date:        t.Date,
			Kategoria:   t.Kategoria,
			Type:        t.Type,
			Amount:      t.Amount,
			Description: t.Description,
		})
	}
	return &database.RequestUserTransactionsResponse{
		Transactions: protoTransactions,
	}
}

### .\backend\database\internal\handlers\grpc\server.go END ###

### .\backend\database\internal\migrations\1_init_schema.down.sql BEGIN ###
DROP TABLE IF EXISTS transactions;
DROP TABLE IF EXISTS users;
### .\backend\database\internal\migrations\1_init_schema.down.sql END ###

### .\backend\database\internal\migrations\1_init_schema.up.sql BEGIN ###
CREATE TABLE IF NOT EXISTS users (
    id BIGINT PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS transactions (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    date BIGINT NOT NULL,
    category TEXT NOT NULL DEFAULT '',
    type TEXT NOT NULL DEFAULT '',
    amount BIGINT NOT NULL DEFAULT 0,
    description TEXT NOT NULL DEFAULT '',
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_transactions_user_id ON transactions(user_id);
### .\backend\database\internal\migrations\1_init_schema.up.sql END ###

### .\backend\database\internal\migrator\migrator.go BEGIN ###
package migrator

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/PrototypeSirius/ruglogger/apperror"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

func Run(log *logrus.Logger, db *sql.DB, migrationsPath string) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return apperror.SystemError(err, 1501, "could not create database driver")
	}
	m, err := migrate.NewWithDatabaseInstance(
		migrationsPath,
		"postgres",
		driver,
	)
	if err != nil {
		return apperror.SystemError(err, 1502, "could not create migrate instance")
	}
	err = m.Up()
	if err == nil || errors.Is(err, migrate.ErrNoChange) {
		return nil
	}
	version, dirty, verErr := m.Version()
	if verErr != nil {
		return apperror.SystemError(verErr, 1503, "could not get migrate version")
	}
	if dirty {
		log.Info(fmt.Sprintf("Dirty migration detected at version %d", version))
		prevVersion := int(version) - 1
		if prevVersion < 0 {
			prevVersion = 0
		}
		if forceErr := m.Force(prevVersion); forceErr != nil {
			return apperror.SystemError(forceErr, 1504, fmt.Sprintf("failed to force rollback to %d", prevVersion))
		}
		log.Info(fmt.Sprintf("Successfully forced version to %d", prevVersion))
		return apperror.SystemError(err, 1505, fmt.Sprintf("dirty migration at version %d rolled back", version))
	}
	return apperror.SystemError(err, 1506, "migration failed")
}

### .\backend\database\internal\migrator\migrator.go END ###

### .\backend\database\internal\models\models.go BEGIN ###
package model

type Transaction struct {
	ID          int64
	Date        int64
	Kategoria   string
	Type        string
	Amount      int64
	Description string
}

### .\backend\database\internal\models\models.go END ###

### .\backend\database\internal\repository\database\db.go BEGIN ###
package dbrepo

import (
	"context"
	"database/internal/config"
	"database/internal/migrator"
	model "database/internal/models"
	"database/sql"
	"fmt"

	"github.com/PrototypeSirius/ruglogger/apperror"
	"github.com/PrototypeSirius/ruglogger/logger"
	"github.com/sirupsen/logrus"
)

type DatabaseRepo struct {
	db *sql.DB
}

func New(log *logrus.Logger, cfg config.DatabaseConfig) (*DatabaseRepo, func() error) {
	connstr := fmt.Sprintf(
		"host=%s port=%d dbname=%s user=%s password=%s sslmode=disable",
		cfg.Host,
		cfg.Port,
		cfg.Database,
		cfg.User,
		cfg.Password,
	)
	db, err := sql.Open("postgres", connstr)
	if err != nil {
		appErr := apperror.SystemError(err, 1201, "Connecting to the database")
		logger.FatalOnError(appErr, "Failed to connect to the database", logrus.Fields{
			"host":     cfg.Host,
			"port":     cfg.Port,
			"database": cfg.Database,
		})
	}
	log.Info("The database is connected")
	err = migrator.Run(log, db, cfg.MigrationsPath)
	if err != nil {
		logger.FatalOnError(err, "Failed to migrate the database", logrus.Fields{
			"migrations_path": cfg.MigrationsPath,
		})
	}
	log.Info("Database migration completed")
	return &DatabaseRepo{db: db}, db.Close
}

func (d *DatabaseRepo) AddTransaction(ctx context.Context, uid int64, t model.Transaction) (string, error) {
	query := `INSERT INTO transactions (user_id, date, category, type, amount, description) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := d.db.ExecContext(ctx, query, uid, t.Date, t.Kategoria, t.Type, t.Amount, t.Description)
	if err != nil {
		return "failed to add transaction", apperror.BadRequestError(err, 1081, "Failed to add transaction")
	}
	return "", nil
}

func (d *DatabaseRepo) AddUser(ctx context.Context, uid int64) (string, error) {
	query := `INSERT INTO users (id) VALUES ($1) ON CONFLICT (id) DO NOTHING`
	_, err := d.db.ExecContext(ctx, query, uid)
	if err != nil {
		return "failed to add user", apperror.BadRequestError(err, 1082, "Failed to add user")
	}
	return "success", nil
}

func (d *DatabaseRepo) DeleteTransaction(ctx context.Context, uid, tid int64) (string, error) {
	query := `DELETE FROM transactions WHERE id = $1`
	res, err := d.db.ExecContext(ctx, query, tid)
	if err != nil {
		return "failed to delete transaction", apperror.BadRequestError(err, 1083, "Failed to delete transaction")
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return "failed to delete transaction", apperror.BadRequestError(err, 1084, "Error getting rows affected")
	}
	if rows == 0 {
		return "Transaction not found", apperror.BadRequestError(err, 1085, "Transaction not found")
	}
	return "success", nil
}

func (d *DatabaseRepo) DeleteUser(ctx context.Context, uid int64) (string, error) {
	query := `DELETE FROM users WHERE id = $1`
	res, err := d.db.ExecContext(ctx, query, uid)
	if err != nil {
		return "failed to delete user", apperror.BadRequestError(err, 1086, "Failed to delete user")
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return "failed to delete user", apperror.BadRequestError(err, 1087, "Error getting rows affected")
	}
	if rows == 0 {
		return "User not found", apperror.BadRequestError(err, 1088, "User not found")
	}
	return "success", nil
}

func (d *DatabaseRepo) EditTransaction(ctx context.Context, uid int64, t model.Transaction) (string, error) {
	query := `UPDATE transactions SET date = $1, category = $2, type = $3, amount = $4, description = $5 WHERE id = $6 AND user_id = $7`
	res, err := d.db.ExecContext(ctx, query, t.Date, t.Kategoria, t.Type, t.Amount, t.Description, t.ID, uid)
	if err != nil {
		return "failed to edit transaction", apperror.BadRequestError(err, 1089, "Failed to edit transaction")
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return "failed to edit transaction", apperror.BadRequestError(err, 1090, "Error getting rows affected")
	}
	if rows == 0 {
		return "Transaction not found", apperror.BadRequestError(err, 1091, "Transaction not found")
	}
	return "success", nil
}

func (d *DatabaseRepo) RequestUserTransactions(ctx context.Context, uid int64) ([]model.Transaction, string, error) {
	query := `SELECT id, date, category, type, amount, description FROM transactions WHERE user_id = $1 ORDER BY date DESC`
	rows, err := d.db.QueryContext(ctx, query, uid)
	if err != nil {
		return []model.Transaction{}, "failed to get transactions", apperror.BadRequestError(err, 1092, "Failed to get transactions")
	}
	var transactions []model.Transaction
	for rows.Next() {
		var t model.Transaction
		err := rows.Scan(&t.ID, &t.Date, &t.Kategoria, &t.Type, &t.Amount, &t.Description)
		if err != nil {
			return []model.Transaction{}, "failed to get transactions", apperror.BadRequestError(err, 1093, "Error scanning transaction")
		}
		transactions = append(transactions, t)
	}
	if err := rows.Err(); err != nil {
		return []model.Transaction{}, "failed to get transactions", apperror.BadRequestError(err, 1094, "Error getting transactions")
	}
	if len(transactions) == 0 {
		return []model.Transaction{}, "No transactions found", nil
	}
	return transactions, "success", nil
}

### .\backend\database\internal\repository\database\db.go END ###

### .\backend\database\internal\services\db.go BEGIN ###
package dbservice

import (
	"context"
	model "database/internal/models"
	dbrepo "database/internal/repository/database"

	"github.com/sirupsen/logrus"
)

type Database struct {
	log                     *logrus.Logger
	addTransaction          AdderTransaction
	addUser                 AdderUser
	deleteTransaction       DeleterTransaction
	deleteUser              DeleterUser
	editTransaction         EditorTransaction
	requestUserTransactions RequesterUserTransactions
}

func New(log *logrus.Logger, db *dbrepo.DatabaseRepo) *Database {
	return &Database{
		log:                     log,
		addTransaction:          db,
		addUser:                 db,
		deleteTransaction:       db,
		deleteUser:              db,
		editTransaction:         db,
		requestUserTransactions: db,
	}
}

type AdderTransaction interface {
	AddTransaction(ctx context.Context, uid int64, t model.Transaction) (string, error)
}

type AdderUser interface {
	AddUser(ctx context.Context, uid int64) (string, error)
}

type DeleterTransaction interface {
	DeleteTransaction(ctx context.Context, uid, tid int64) (string, error)
}

type DeleterUser interface {
	DeleteUser(ctx context.Context, uid int64) (string, error)
}

type EditorTransaction interface {
	EditTransaction(ctx context.Context, uid int64, t model.Transaction) (string, error)
}

type RequesterUserTransactions interface {
	RequestUserTransactions(ctx context.Context, uid int64) ([]model.Transaction, string, error)
}

func (d *Database) AddTransaction(ctx context.Context, uid int64, t model.Transaction) (string, error) {
	return d.addTransaction.AddTransaction(ctx, uid, t)
}

func (d *Database) AddUser(ctx context.Context, uid int64) (string, error) {
	return d.addUser.AddUser(ctx, uid)
}

func (d *Database) DeleteTransaction(ctx context.Context, uid, tid int64) (string, error) {
	return d.deleteTransaction.DeleteTransaction(ctx, uid, tid)
}

func (d *Database) DeleteUser(ctx context.Context, uid int64) (string, error) {
	return d.deleteUser.DeleteUser(ctx, uid)
}

func (d *Database) EditTransaction(ctx context.Context, uid int64, t model.Transaction) (string, error) {
	return d.editTransaction.EditTransaction(ctx, uid, t)
}

func (d *Database) RequestUserTransactions(ctx context.Context, uid int64) ([]model.Transaction, string, error) {
	return d.requestUserTransactions.RequestUserTransactions(ctx, uid)
}

### .\backend\database\internal\services\db.go END ###

### .\backend\manager\cmd\main.go BEGIN ###
package main

import (
	"io"
	"log"
	"manager/internal/app"
	"manager/internal/config"
	"os"
	"os/signal"
	"syscall"

	"github.com/PrototypeSirius/ruglogger/logger"
	"github.com/sirupsen/logrus"
)

var level string

func main() {
	cfg := config.MustLoad()
	log, logCloser := logInit(cfg.Env)
	defer func() {
		if err := logCloser(); err != nil {
			os.Stderr.WriteString("Error closing log file: " + err.Error() + "\n")
		}
	}()
	log.Info("Config has been successfully loaded")
	application := app.New(log, *cfg)
	log.Info("Application has been successfully initialized")
	go func() {
		application.HTTPApp.MustRun()
	}()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop

	log.Info("Gracefully stopped")
}

func logInit(env string) (*logrus.Logger, func() error) {
	switch env {
	case "production":
		level = "warn"
	case "local":
		level = "info"
	default:
		level = "debug"
	}
	logFile, err := os.OpenFile("app.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Critical error: failed to open log file: %v\n", err)
	}
	multiWriter := io.MultiWriter(os.Stdout, logFile)
	logger.Init(
		logger.WithLevel(level),
		logger.WithOutput(multiWriter),
	)
	log := logger.Get()
	log.Info("The logging system has been successfully initialized.")
	return log, logFile.Close
}

### .\backend\manager\cmd\main.go END ###

### .\backend\manager\config\config.yaml BEGIN ###
env: "local" # local, production

httpserver:
  port: 8080
  host: localhost

client:
  bot:
    port: 2011
    host: bot

  database:
    host: database-service
    port: 2012
  
  ml:
    host: ml-service
    port: 2014
  

### .\backend\manager\config\config.yaml END ###

### .\backend\manager\Dockerfile BEGIN ###
FROM golang:alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o manager-app cmd/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

RUN addgroup -S appgroup && adduser -S appuser -G appgroup


COPY --from=builder /app/manager-app .

COPY --from=builder /app/config/config.yaml ./config/config.yaml

RUN touch app.log && chown appuser:appgroup app.log

USER appuser

EXPOSE 8080

CMD ["./manager-app", "--config=./config/config.yaml"]
### .\backend\manager\Dockerfile END ###

### .\backend\manager\go.mod BEGIN ###
module manager

go 1.24.2

require (
	github.com/PrototypeSirius/protos_service v0.0.0-20251206133529-ad042aab1d47
	github.com/PrototypeSirius/ruglogger v0.0.0-20251116235218-0aa354b7deb2
	github.com/gin-contrib/cors v1.7.6
	github.com/gin-gonic/gin v1.11.0
	github.com/ilyakaznacheev/cleanenv v1.5.0
	github.com/sirupsen/logrus v1.9.3
	google.golang.org/grpc v1.77.0
)

require (
	github.com/BurntSushi/toml v1.2.1 // indirect
	github.com/bytedance/sonic v1.14.0 // indirect
	github.com/bytedance/sonic/loader v0.3.0 // indirect
	github.com/cloudwego/base64x v0.1.6 // indirect
	github.com/gabriel-vasile/mimetype v1.4.9 // indirect
	github.com/gin-contrib/sse v1.1.0 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.27.0 // indirect
	github.com/goccy/go-json v0.10.5 // indirect
	github.com/goccy/go-yaml v1.18.0 // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/cpuid/v2 v2.3.0 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pelletier/go-toml/v2 v2.2.4 // indirect
	github.com/quic-go/qpack v0.5.1 // indirect
	github.com/quic-go/quic-go v0.54.0 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/ugorji/go/codec v1.3.0 // indirect
	go.uber.org/mock v0.5.0 // indirect
	golang.org/x/arch v0.20.0 // indirect
	golang.org/x/crypto v0.43.0 // indirect
	golang.org/x/mod v0.28.0 // indirect
	golang.org/x/net v0.46.1-0.20251013234738-63d1a5100f82 // indirect
	golang.org/x/sync v0.17.0 // indirect
	golang.org/x/sys v0.37.0 // indirect
	golang.org/x/text v0.30.0 // indirect
	golang.org/x/tools v0.37.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251022142026-3a174f9686a8 // indirect
	google.golang.org/protobuf v1.36.10 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	olympos.io/encoding/edn v0.0.0-20201019073823-d3554ca0b0a3 // indirect
)

### .\backend\manager\go.mod END ###

### .\backend\manager\go.sum BEGIN ###
github.com/BurntSushi/toml v1.2.1 h1:9F2/+DoOYIOksmaJFPw1tGFy1eDnIJXg+UHjuD8lTak=
github.com/BurntSushi/toml v1.2.1/go.mod h1:CxXYINrC8qIiEnFrOxCa7Jy5BFHlXnUU2pbicEuybxQ=
github.com/PrototypeSirius/protos_service v0.0.0-20251206133529-ad042aab1d47 h1:slbJBbdCL1XnGLqam6WMGELX1VmilIE26xWRCXB1BzE=
github.com/PrototypeSirius/protos_service v0.0.0-20251206133529-ad042aab1d47/go.mod h1:hOJuOHkvh867yXdbopMLGrc0gRiSI+ENZWRBqSZSiNw=
github.com/PrototypeSirius/ruglogger v0.0.0-20251116235218-0aa354b7deb2 h1:67qBAnVo8NJxag8sBFZNsobQEWHhr0np3QhXfRAl7R8=
github.com/PrototypeSirius/ruglogger v0.0.0-20251116235218-0aa354b7deb2/go.mod h1:3KeQKrOLvDy6xVl90WmKDW3VyQwtWvvvQb3wWkR059s=
github.com/bytedance/sonic v1.14.0 h1:/OfKt8HFw0kh2rj8N0F6C/qPGRESq0BbaNZgcNXXzQQ=
github.com/bytedance/sonic v1.14.0/go.mod h1:WoEbx8WTcFJfzCe0hbmyTGrfjt8PzNEBdxlNUO24NhA=
github.com/bytedance/sonic/loader v0.3.0 h1:dskwH8edlzNMctoruo8FPTJDF3vLtDT0sXZwvZJyqeA=
github.com/bytedance/sonic/loader v0.3.0/go.mod h1:N8A3vUdtUebEY2/VQC0MyhYeKUFosQU6FxH2JmUe6VI=
github.com/cloudwego/base64x v0.1.6 h1:t11wG9AECkCDk5fMSoxmufanudBtJ+/HemLstXDLI2M=
github.com/cloudwego/base64x v0.1.6/go.mod h1:OFcloc187FXDaYHvrNIjxSe8ncn0OOM8gEHfghB2IPU=
github.com/creack/pty v1.1.9/go.mod h1:oKZEueFk5CKHvIhNR5MUki03XCEU+Q6VDXinZuGJ33E=
github.com/davecgh/go-spew v1.1.0/go.mod h1:J7Y8YcW2NihsgmVo/mv3lAwl/skON4iLHjSsI+c5H38=
github.com/davecgh/go-spew v1.1.1 h1:vj9j/u1bqnvCEfJOwUhtlOARqs3+rkHYY13jYWTU97c=
github.com/davecgh/go-spew v1.1.1/go.mod h1:J7Y8YcW2NihsgmVo/mv3lAwl/skON4iLHjSsI+c5H38=
github.com/gabriel-vasile/mimetype v1.4.9 h1:5k+WDwEsD9eTLL8Tz3L0VnmVh9QxGjRmjBvAG7U/oYY=
github.com/gabriel-vasile/mimetype v1.4.9/go.mod h1:WnSQhFKJuBlRyLiKohA/2DtIlPFAbguNaG7QCHcyGok=
github.com/gin-contrib/cors v1.7.6 h1:3gQ8GMzs1Ylpf70y8bMw4fVpycXIeX1ZemuSQIsnQQY=
github.com/gin-contrib/cors v1.7.6/go.mod h1:Ulcl+xN4jel9t1Ry8vqph23a60FwH9xVLd+3ykmTjOk=
github.com/gin-contrib/sse v1.1.0 h1:n0w2GMuUpWDVp7qSpvze6fAu9iRxJY4Hmj6AmBOU05w=
github.com/gin-contrib/sse v1.1.0/go.mod h1:hxRZ5gVpWMT7Z0B0gSNYqqsSCNIJMjzvm6fqCz9vjwM=
github.com/gin-gonic/gin v1.11.0 h1:OW/6PLjyusp2PPXtyxKHU0RbX6I/l28FTdDlae5ueWk=
github.com/gin-gonic/gin v1.11.0/go.mod h1:+iq/FyxlGzII0KHiBGjuNn4UNENUlKbGlNmc+W50Dls=
github.com/go-logr/logr v1.4.3 h1:CjnDlHq8ikf6E492q6eKboGOC0T8CDaOvkHCIg8idEI=
github.com/go-logr/logr v1.4.3/go.mod h1:9T104GzyrTigFIr8wt5mBrctHMim0Nb2HLGrmQ40KvY=
github.com/go-logr/stdr v1.2.2 h1:hSWxHoqTgW2S2qGc0LTAI563KZ5YKYRhT3MFKZMbjag=
github.com/go-logr/stdr v1.2.2/go.mod h1:mMo/vtBO5dYbehREoey6XUKy/eSumjCCveDpRre4VKE=
github.com/go-playground/assert/v2 v2.2.0 h1:JvknZsQTYeFEAhQwI4qEt9cyV5ONwRHC+lYKSsYSR8s=
github.com/go-playground/assert/v2 v2.2.0/go.mod h1:VDjEfimB/XKnb+ZQfWdccd7VUvScMdVu0Titje2rxJ4=
github.com/go-playground/locales v0.14.1 h1:EWaQ/wswjilfKLTECiXz7Rh+3BjFhfDFKv/oXslEjJA=
github.com/go-playground/locales v0.14.1/go.mod h1:hxrqLVvrK65+Rwrd5Fc6F2O76J/NuW9t0sjnWqG1slY=
github.com/go-playground/universal-translator v0.18.1 h1:Bcnm0ZwsGyWbCzImXv+pAJnYK9S473LQFuzCbDbfSFY=
github.com/go-playground/universal-translator v0.18.1/go.mod h1:xekY+UJKNuX9WP91TpwSH2VMlDf28Uj24BCp08ZFTUY=
github.com/go-playground/validator/v10 v10.27.0 h1:w8+XrWVMhGkxOaaowyKH35gFydVHOvC0/uWoy2Fzwn4=
github.com/go-playground/validator/v10 v10.27.0/go.mod h1:I5QpIEbmr8On7W0TktmJAumgzX4CA1XNl4ZmDuVHKKo=
github.com/goccy/go-json v0.10.5 h1:Fq85nIqj+gXn/S5ahsiTlK3TmC85qgirsdTP/+DeaC4=
github.com/goccy/go-json v0.10.5/go.mod h1:oq7eo15ShAhp70Anwd5lgX2pLfOS3QCiwU/PULtXL6M=
github.com/goccy/go-yaml v1.18.0 h1:8W7wMFS12Pcas7KU+VVkaiCng+kG8QiFeFwzFb+rwuw=
github.com/goccy/go-yaml v1.18.0/go.mod h1:XBurs7gK8ATbW4ZPGKgcbrY1Br56PdM69F7LkFRi1kA=
github.com/golang/protobuf v1.5.4 h1:i7eJL8qZTpSEXOPTxNKhASYpMn+8e5Q6AdndVa1dWek=
github.com/golang/protobuf v1.5.4/go.mod h1:lnTiLA8Wa4RWRcIUkrtSVa5nRhsEGBg48fD6rSs7xps=
github.com/google/go-cmp v0.7.0 h1:wk8382ETsv4JYUZwIsn6YpYiWiBsYLSJiTsyBybVuN8=
github.com/google/go-cmp v0.7.0/go.mod h1:pXiqmnSA92OHEEa9HXL2W4E7lf9JzCmGVUdgjX3N/iU=
github.com/google/gofuzz v1.0.0/go.mod h1:dBl0BpW6vV/+mYPU4Po3pmUjxk6FQPldtuIdl/M65Eg=
github.com/google/uuid v1.6.0 h1:NIvaJDMOsjHA8n1jAhLSgzrAzy1Hgr+hNrb57e+94F0=
github.com/google/uuid v1.6.0/go.mod h1:TIyPZe4MgqvfeYDBFedMoGGpEw/LqOeaOT+nhxU+yHo=
github.com/gorilla/websocket v1.5.3 h1:saDtZ6Pbx/0u+bgYQ3q96pZgCzfhKXGPqt7kZ72aNNg=
github.com/gorilla/websocket v1.5.3/go.mod h1:YR8l580nyteQvAITg2hZ9XVh4b55+EU/adAjf1fMHhE=
github.com/ilyakaznacheev/cleanenv v1.5.0 h1:0VNZXggJE2OYdXE87bfSSwGxeiGt9moSR2lOrsHHvr4=
github.com/ilyakaznacheev/cleanenv v1.5.0/go.mod h1:a5aDzaJrLCQZsazHol1w8InnDcOX0OColm64SlIi6gk=
github.com/joho/godotenv v1.5.1 h1:7eLL/+HRGLY0ldzfGMeQkb7vMd0as4CfYvUVzLqw0N0=
github.com/joho/godotenv v1.5.1/go.mod h1:f4LDr5Voq0i2e/R5DDNOoa2zzDfwtkZa6DnEwAbqwq4=
github.com/json-iterator/go v1.1.12 h1:PV8peI4a0ysnczrg+LtxykD8LfKY9ML6u2jnxaEnrnM=
github.com/json-iterator/go v1.1.12/go.mod h1:e30LSqwooZae/UwlEbR2852Gd8hjQvJoHmT4TnhNGBo=
github.com/klauspost/cpuid/v2 v2.3.0 h1:S4CRMLnYUhGeDFDqkGriYKdfoFlDnMtqTiI/sFzhA9Y=
github.com/klauspost/cpuid/v2 v2.3.0/go.mod h1:hqwkgyIinND0mEev00jJYCxPNVRVXFQeu1XKlok6oO0=
github.com/kr/pretty v0.3.0 h1:WgNl7dwNpEZ6jJ9k1snq4pZsg7DOEN8hP9Xw0Tsjwk0=
github.com/kr/pretty v0.3.0/go.mod h1:640gp4NfQd8pI5XOwp5fnNeVWj67G7CFk/SaSQn7NBk=
github.com/kr/text v0.2.0 h1:5Nx0Ya0ZqY2ygV366QzturHI13Jq95ApcVaJBhpS+AY=
github.com/kr/text v0.2.0/go.mod h1:eLer722TekiGuMkidMxC/pM04lWEeraHUUmBw8l2grE=
github.com/leodido/go-urn v1.4.0 h1:WT9HwE9SGECu3lg4d/dIA+jxlljEa1/ffXKmRjqdmIQ=
github.com/leodido/go-urn v1.4.0/go.mod h1:bvxc+MVxLKB4z00jd1z+Dvzr47oO32F/QSNjSBOlFxI=
github.com/mattn/go-isatty v0.0.20 h1:xfD0iDuEKnDkl03q4limB+vH+GxLEtL/jb4xVJSWWEY=
github.com/mattn/go-isatty v0.0.20/go.mod h1:W+V8PltTTMOvKvAeJH7IuucS94S2C6jfK/D7dTCTo3Y=
github.com/modern-go/concurrent v0.0.0-20180228061459-e0a39a4cb421/go.mod h1:6dJC0mAP4ikYIbvyc7fijjWJddQyLn8Ig3JB5CqoB9Q=
github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd h1:TRLaZ9cD/w8PVh93nsPXa1VrQ6jlwL5oN8l14QlcNfg=
github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd/go.mod h1:6dJC0mAP4ikYIbvyc7fijjWJddQyLn8Ig3JB5CqoB9Q=
github.com/modern-go/reflect2 v1.0.2 h1:xBagoLtFs94CBntxluKeaWgTMpvLxC4ur3nMaC9Gz0M=
github.com/modern-go/reflect2 v1.0.2/go.mod h1:yWuevngMOJpCy52FWWMvUC8ws7m/LJsjYzDa0/r8luk=
github.com/pelletier/go-toml/v2 v2.2.4 h1:mye9XuhQ6gvn5h28+VilKrrPoQVanw5PMw/TB0t5Ec4=
github.com/pelletier/go-toml/v2 v2.2.4/go.mod h1:2gIqNv+qfxSVS7cM2xJQKtLSTLUE9V8t9Stt+h56mCY=
github.com/pmezard/go-difflib v1.0.0 h1:4DBwDE0NGyQoBHbLQYPwSUPoCMWR5BEzIk/f1lZbAQM=
github.com/pmezard/go-difflib v1.0.0/go.mod h1:iKH77koFhYxTK1pcRnkKkqfTogsbg7gZNVY4sRDYZ/4=
github.com/quic-go/qpack v0.5.1 h1:giqksBPnT/HDtZ6VhtFKgoLOWmlyo9Ei6u9PqzIMbhI=
github.com/quic-go/qpack v0.5.1/go.mod h1:+PC4XFrEskIVkcLzpEkbLqq1uCoxPhQuvK5rH1ZgaEg=
github.com/quic-go/quic-go v0.54.0 h1:6s1YB9QotYI6Ospeiguknbp2Znb/jZYjZLRXn9kMQBg=
github.com/quic-go/quic-go v0.54.0/go.mod h1:e68ZEaCdyviluZmy44P6Iey98v/Wfz6HCjQEm+l8zTY=
github.com/rogpeppe/go-internal v1.8.0 h1:FCbCCtXNOY3UtUuHUYaghJg4y7Fd14rXifAYUAtL9R8=
github.com/rogpeppe/go-internal v1.8.0/go.mod h1:WmiCO8CzOY8rg0OYDC4/i/2WRWAB6poM+XZ2dLUbcbE=
github.com/sirupsen/logrus v1.9.3 h1:dueUQJ1C2q9oE3F7wvmSGAaVtTmUizReu6fjN8uqzbQ=
github.com/sirupsen/logrus v1.9.3/go.mod h1:naHLuLoDiP4jHNo9R0sCBMtWGeIprob74mVsIT4qYEQ=
github.com/stretchr/objx v0.1.0/go.mod h1:HFkY916IF+rwdDfMAkV7OtwuqBVzrE8GR6GFx+wExME=
github.com/stretchr/objx v0.4.0/go.mod h1:YvHI0jy2hoMjB+UWwv71VJQ9isScKT/TqJzVSSt89Yw=
github.com/stretchr/objx v0.5.0/go.mod h1:Yh+to48EsGEfYuaHDzXPcE3xhTkx73EhmCGUpEOglKo=
github.com/stretchr/testify v1.3.0/go.mod h1:M5WIy9Dh21IEIfnGCwXGc5bZfKNJtfHm1UVUgZn+9EI=
github.com/stretchr/testify v1.7.0/go.mod h1:6Fq8oRcR53rry900zMqJjRRixrwX3KX962/h/Wwjteg=
github.com/stretchr/testify v1.7.1/go.mod h1:6Fq8oRcR53rry900zMqJjRRixrwX3KX962/h/Wwjteg=
github.com/stretchr/testify v1.8.0/go.mod h1:yNjHg4UonilssWZ8iaSj1OCr/vHnekPRkoO+kdMU+MU=
github.com/stretchr/testify v1.8.1/go.mod h1:w2LPCIKwWwSfY2zedu0+kehJoqGctiVI29o6fzry7u4=
github.com/stretchr/testify v1.11.1 h1:7s2iGBzp5EwR7/aIZr8ao5+dra3wiQyKjjFuvgVKu7U=
github.com/stretchr/testify v1.11.1/go.mod h1:wZwfW3scLgRK+23gO65QZefKpKQRnfz6sD981Nm4B6U=
github.com/twitchyliquid64/golang-asm v0.15.1 h1:SU5vSMR7hnwNxj24w34ZyCi/FmDZTkS4MhqMhdFk5YI=
github.com/twitchyliquid64/golang-asm v0.15.1/go.mod h1:a1lVb/DtPvCB8fslRZhAngC2+aY1QWCk3Cedj/Gdt08=
github.com/ugorji/go/codec v1.3.0 h1:Qd2W2sQawAfG8XSvzwhBeoGq71zXOC/Q1E9y/wUcsUA=
github.com/ugorji/go/codec v1.3.0/go.mod h1:pRBVtBSKl77K30Bv8R2P+cLSGaTtex6fsA2Wjqmfxj4=
go.opentelemetry.io/auto/sdk v1.2.1 h1:jXsnJ4Lmnqd11kwkBV2LgLoFMZKizbCi5fNZ/ipaZ64=
go.opentelemetry.io/auto/sdk v1.2.1/go.mod h1:KRTj+aOaElaLi+wW1kO/DZRXwkF4C5xPbEe3ZiIhN7Y=
go.opentelemetry.io/otel v1.38.0 h1:RkfdswUDRimDg0m2Az18RKOsnI8UDzppJAtj01/Ymk8=
go.opentelemetry.io/otel v1.38.0/go.mod h1:zcmtmQ1+YmQM9wrNsTGV/q/uyusom3P8RxwExxkZhjM=
go.opentelemetry.io/otel/metric v1.38.0 h1:Kl6lzIYGAh5M159u9NgiRkmoMKjvbsKtYRwgfrA6WpA=
go.opentelemetry.io/otel/metric v1.38.0/go.mod h1:kB5n/QoRM8YwmUahxvI3bO34eVtQf2i4utNVLr9gEmI=
go.opentelemetry.io/otel/sdk v1.38.0 h1:l48sr5YbNf2hpCUj/FoGhW9yDkl+Ma+LrVl8qaM5b+E=
go.opentelemetry.io/otel/sdk v1.38.0/go.mod h1:ghmNdGlVemJI3+ZB5iDEuk4bWA3GkTpW+DOoZMYBVVg=
go.opentelemetry.io/otel/sdk/metric v1.38.0 h1:aSH66iL0aZqo//xXzQLYozmWrXxyFkBJ6qT5wthqPoM=
go.opentelemetry.io/otel/sdk/metric v1.38.0/go.mod h1:dg9PBnW9XdQ1Hd6ZnRz689CbtrUp0wMMs9iPcgT9EZA=
go.opentelemetry.io/otel/trace v1.38.0 h1:Fxk5bKrDZJUH+AMyyIXGcFAPah0oRcT+LuNtJrmcNLE=
go.opentelemetry.io/otel/trace v1.38.0/go.mod h1:j1P9ivuFsTceSWe1oY+EeW3sc+Pp42sO++GHkg4wwhs=
go.uber.org/mock v0.5.0 h1:KAMbZvZPyBPWgD14IrIQ38QCyjwpvVVV6K/bHl1IwQU=
go.uber.org/mock v0.5.0/go.mod h1:ge71pBPLYDk7QIi1LupWxdAykm7KIEFchiOqd6z7qMM=
golang.org/x/arch v0.20.0 h1:dx1zTU0MAE98U+TQ8BLl7XsJbgze2WnNKF/8tGp/Q6c=
golang.org/x/arch v0.20.0/go.mod h1:bdwinDaKcfZUGpH09BB7ZmOfhalA8lQdzl62l8gGWsk=
golang.org/x/crypto v0.43.0 h1:dduJYIi3A3KOfdGOHX8AVZ/jGiyPa3IbBozJ5kNuE04=
golang.org/x/crypto v0.43.0/go.mod h1:BFbav4mRNlXJL4wNeejLpWxB7wMbc79PdRGhWKncxR0=
golang.org/x/mod v0.28.0 h1:gQBtGhjxykdjY9YhZpSlZIsbnaE2+PgjfLWUQTnoZ1U=
golang.org/x/mod v0.28.0/go.mod h1:yfB/L0NOf/kmEbXjzCPOx1iK1fRutOydrCMsqRhEBxI=
golang.org/x/net v0.46.1-0.20251013234738-63d1a5100f82 h1:6/3JGEh1C88g7m+qzzTbl3A0FtsLguXieqofVLU/JAo=
golang.org/x/net v0.46.1-0.20251013234738-63d1a5100f82/go.mod h1:Q9BGdFy1y4nkUwiLvT5qtyhAnEHgnQ/zd8PfU6nc210=
golang.org/x/sync v0.17.0 h1:l60nONMj9l5drqw6jlhIELNv9I0A4OFgRsG9k2oT9Ug=
golang.org/x/sync v0.17.0/go.mod h1:9KTHXmSnoGruLpwFjVSX0lNNA75CykiMECbovNTZqGI=
golang.org/x/sys v0.0.0-20220715151400-c0bba94af5f8/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
golang.org/x/sys v0.6.0/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
golang.org/x/sys v0.37.0 h1:fdNQudmxPjkdUTPnLn5mdQv7Zwvbvpaxqs831goi9kQ=
golang.org/x/sys v0.37.0/go.mod h1:OgkHotnGiDImocRcuBABYBEXf8A9a87e/uXjp9XT3ks=
golang.org/x/text v0.30.0 h1:yznKA/E9zq54KzlzBEAWn1NXSQ8DIp/NYMy88xJjl4k=
golang.org/x/text v0.30.0/go.mod h1:yDdHFIX9t+tORqspjENWgzaCVXgk0yYnYuSZ8UzzBVM=
golang.org/x/tools v0.37.0 h1:DVSRzp7FwePZW356yEAChSdNcQo6Nsp+fex1SUW09lE=
golang.org/x/tools v0.37.0/go.mod h1:MBN5QPQtLMHVdvsbtarmTNukZDdgwdwlO5qGacAzF0w=
gonum.org/v1/gonum v0.16.0 h1:5+ul4Swaf3ESvrOnidPp4GZbzf0mxVQpDCYUQE7OJfk=
gonum.org/v1/gonum v0.16.0/go.mod h1:fef3am4MQ93R2HHpKnLk4/Tbh/s0+wqD5nfa6Pnwy4E=
google.golang.org/genproto/googleapis/rpc v0.0.0-20251022142026-3a174f9686a8 h1:M1rk8KBnUsBDg1oPGHNCxG4vc1f49epmTO7xscSajMk=
google.golang.org/genproto/googleapis/rpc v0.0.0-20251022142026-3a174f9686a8/go.mod h1:7i2o+ce6H/6BluujYR+kqX3GKH+dChPTQU19wjRPiGk=
google.golang.org/grpc v1.77.0 h1:wVVY6/8cGA6vvffn+wWK5ToddbgdU3d8MNENr4evgXM=
google.golang.org/grpc v1.77.0/go.mod h1:z0BY1iVj0q8E1uSQCjL9cppRj+gnZjzDnzV0dHhrNig=
google.golang.org/protobuf v1.36.10 h1:AYd7cD/uASjIL6Q9LiTjz8JLcrh/88q5UObnmY3aOOE=
google.golang.org/protobuf v1.36.10/go.mod h1:HTf+CrKn2C3g5S8VImy6tdcUvCska2kB7j23XfzDpco=
gopkg.in/check.v1 v0.0.0-20161208181325-20d25e280405/go.mod h1:Co6ibVJAznAaIkqp8huTwlJQCZ016jof/cbN4VW5Yz0=
gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c h1:Hei/4ADfdWqJk1ZMxUNpqntNwaWcugrBjAiHlqqRiVk=
gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c/go.mod h1:JHkPIbrfpd72SG/EVd6muEfDQjcINNoR0C8j2r3qZ4Q=
gopkg.in/yaml.v3 v3.0.0-20200313102051-9f266ea9e77c/go.mod h1:K4uyk7z7BCEPqu6E+C64Yfv1cQ7kz7rIZviUmN+EgEM=
gopkg.in/yaml.v3 v3.0.1 h1:fxVm/GzAzEWqLHuvctI91KS9hhNmmWOoWu0XTYJS7CA=
gopkg.in/yaml.v3 v3.0.1/go.mod h1:K4uyk7z7BCEPqu6E+C64Yfv1cQ7kz7rIZviUmN+EgEM=
olympos.io/encoding/edn v0.0.0-20201019073823-d3554ca0b0a3 h1:slmdOY3vp8a7KQbHkL+FLbvbkgMqmXojpFUO/jENuqQ=
olympos.io/encoding/edn v0.0.0-20201019073823-d3554ca0b0a3/go.mod h1:oVgVk4OWVDi43qWBEyGhXgYxt7+ED4iYNpTngSLX2Iw=

### .\backend\manager\go.sum END ###

### .\backend\manager\internal\app\app.go BEGIN ###
package app

import (
	httpapp "manager/internal/app/http"
	"manager/internal/config"
	botrepo "manager/internal/repository/bot"
	databaserepo "manager/internal/repository/database"
	"manager/internal/router"
	httpservice "manager/internal/services"

	"github.com/PrototypeSirius/ruglogger/logger"
	"github.com/sirupsen/logrus"
)

type App struct {
	HTTPApp *httpapp.HTTPApp
}

func New(log *logrus.Logger, cfg config.Config) *App {
	// 1. Initialize Repositories (gRPC Clients)
	botClient, err := botrepo.New(log, cfg.Client.Bot.Host, cfg.Client.Bot.Port)
	if err != nil {
		logger.FatalOnError(err, "error init bot client")
	}

	databaseClient, err := databaserepo.New(log, cfg.Client.Database.Host, cfg.Client.Database.Port)
	if err != nil {
		logger.FatalOnError(err, "error init database client")
	}

	// 2. Initialize Service (Business Logic)
	managerService := httpservice.New(log, botClient, databaseClient)

	// 3. Initialize Handler (Router)
	httpHandler := router.New(log, managerService)

	// 4. Initialize HTTP Server (Gin)
	httpServer, engine := httpapp.New(log, cfg.HttpServer.Port)

	// 5. Register Routes
	httpHandler.RouterRegister(engine)

	return &App{HTTPApp: httpServer}
}

### .\backend\manager\internal\app\app.go END ###

### .\backend\manager\internal\app\http\app.go BEGIN ###
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

### .\backend\manager\internal\app\http\app.go END ###

### .\backend\manager\internal\config\config.go BEGIN ###
package config

import (
	"flag"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env        string       `yaml:"env" env-default:"local"` // local, production
	HttpServer HttpConfig   `yaml:"httpserver"`              // http server config
	Client     ClientConfig `yaml:"client"`                  // client config
}

type ClientConfig struct {
	Bot      HttpConfig `yaml:"bot"`      // bot client config
	Database HttpConfig `yaml:"database"` // database client config
	ML       HttpConfig `yaml:"ml"`       // ml client config
}

type HttpConfig struct {
	Port int    `yaml:"port" env-required:"true"` // HTTP port
	Host string `yaml:"host" env-required:"true"` // HTTP host
}

func MustLoad() *Config {
	path := fechPathConfig()
	if path == "" {
		panic("config path is empty")
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file not found: " + path)
	}
	var cfg Config
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("failed to load config: " + err.Error())
	}
	return &cfg
}

func fechPathConfig() string {
	var res string
	//--config="path/to/config.yaml"
	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()
	return res
}

### .\backend\manager\internal\config\config.go END ###

### .\backend\manager\internal\models\models.go BEGIN ###
package model

type Transaction struct {
	ID          int64  `json:"id"`
	Date        int64  `json:"date"`
	Kategoria   string `json:"kategoria"`
	Type        string `json:"type"`
	Amount      int64  `json:"amount"`
	Description string `json:"description"`
}

type UserID struct {
	Uid int64 `json:"uid"`
}

type Metrics struct {
}

### .\backend\manager\internal\models\models.go END ###

### .\backend\manager\internal\repository\bot\client.go BEGIN ###
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

### .\backend\manager\internal\repository\bot\client.go END ###

### .\backend\manager\internal\repository\database\client.go BEGIN ###
package databaserepo

import (
	"context"
	"fmt"
	model "manager/internal/models"

	cyberdatabase "github.com/PrototypeSirius/protos_service/gen/cybergarden/database"
	"github.com/PrototypeSirius/ruglogger/apperror"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type DBClient struct {
	db  cyberdatabase.DatabaseClient
	log *logrus.Logger
}

func New(log *logrus.Logger, host string, port int) (*DBClient, error) {
	dbaddr := fmt.Sprintf("%s:%d", host, port)
	dcc, err := grpc.NewClient(dbaddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, apperror.SystemError(err, 1021, "error load database service")
	}
	return &DBClient{db: cyberdatabase.NewDatabaseClient(dcc), log: logrus.New()}, nil
}

func (c *DBClient) AddUser(ctx context.Context, uid int64) error {
	resp, err := c.db.AddUser(ctx, &cyberdatabase.AddUserRequest{UserID: uid})
	if err != nil {
		return apperror.SystemError(err, 1031, "grpc add user failed")
	}
	if resp.ErrorMes != "" && resp.ErrorMes != "success" {
		c.log.Warnf("AddUser response: %s", resp.ErrorMes)
	}
	return nil
}

func (c *DBClient) AddTransaction(ctx context.Context, uid int64, t model.Transaction) error {
	_, err := c.db.AddTransaction(ctx, &cyberdatabase.AddTransactionRequest{
		UserID:      uid,
		Transaction: mapModelToProto(t),
	})
	if err != nil {
		return apperror.SystemError(err, 1032, "grpc add transaction failed")
	}
	return nil
}

func mapModelToProto(t model.Transaction) *cyberdatabase.Transaction {
	return &cyberdatabase.Transaction{
		ID:          t.ID,
		Date:        t.Date,
		Kategoria:   t.Kategoria,
		Type:        t.Type,
		Amount:      t.Amount,
		Description: t.Description,
	}
}

func (c *DBClient) DeleteTransaction(ctx context.Context, uid, tid int64) error {
	_, err := c.db.DeleteTransaction(ctx, &cyberdatabase.DeleteTransactionRequest{
		UserID: uid,
		ID:     tid,
	})
	if err != nil {
		return apperror.SystemError(err, 1033, "grpc delete transaction failed")
	}
	return nil
}

func (c *DBClient) EditTransaction(ctx context.Context, uid int64, t model.Transaction) error {
	_, err := c.db.EditTransaction(ctx, &cyberdatabase.EditTransactionRequest{
		UserID:      uid,
		Transaction: mapModelToProto(t),
	})
	if err != nil {
		return apperror.SystemError(err, 1034, "grpc edit transaction failed")
	}
	return nil
}

func (c *DBClient) RequestUserTransactions(ctx context.Context, uid int64) ([]model.Transaction, error) {
	resp, err := c.db.RequestUserTransactions(ctx, &cyberdatabase.RequestUserTransactionsRequest{
		UserID: uid,
	})
	if err != nil {
		return nil, apperror.SystemError(err, 1035, "grpc request transactions failed")
	}
	var transactions []model.Transaction
	for _, t := range resp.GetTransactions() {
		transactions = append(transactions, model.Transaction{
			ID:          t.GetID(),
			Date:        t.GetDate(),
			Kategoria:   t.GetKategoria(),
			Type:        t.GetType(),
			Amount:      t.GetAmount(),
			Description: t.GetDescription(),
		})
	}

	return transactions, nil
}

### .\backend\manager\internal\repository\database\client.go END ###

### .\backend\manager\internal\router\router.go BEGIN ###
package router

import (
	"context"
	"errors"
	model "manager/internal/models"
	"net/http"
	"strings"

	"github.com/PrototypeSirius/ruglogger/apperror"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type ServiceAPI interface {
	AuthUser(ctx context.Context, initData string) (int64, error)
	AddTransaction(ctx context.Context, uid int64, t model.Transaction) error
	DeleteTransaction(ctx context.Context, uid, tid int64) error
	EditTransaction(ctx context.Context, uid int64, t model.Transaction) error
	GetHistory(ctx context.Context, uid int64) ([]model.Transaction, error)
}

type Handler struct {
	log     *logrus.Logger
	service ServiceAPI
}

func New(log *logrus.Logger, service ServiceAPI) *Handler {
	return &Handler{
		log:     log,
		service: service,
	}
}

func (h *Handler) RouterRegister(r *gin.Engine) {
	api := r.Group("/webapp")
	api.Use(h.authUser)
	{
		api.GET("/datametrics", h.requestDataMetrics)
		api.GET("/datahistory", h.requestDataHistory)
		api.POST("/addt", h.addTransaction)
		api.POST("/deletet", h.deleteTransaction)
		api.POST("/updatet", h.updateTransaction)
	}
}

func (h *Handler) requestDataMetrics(c *gin.Context) {

}

func (h *Handler) requestDataHistory(c *gin.Context) {
	uid := c.GetInt64("uid")

	history, err := h.service.GetHistory(c.Request.Context(), uid)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": history})
}

func (h *Handler) addTransaction(c *gin.Context) {
	uid := c.GetInt64("uid")

	var t model.Transaction
	if err := c.BindJSON(&t); err != nil {
		_ = c.Error(apperror.BadRequestError(err, 4001, "invalid json body"))
		return
	}

	if err := h.service.AddTransaction(c.Request.Context(), uid, t); err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "success"})
}

func (h *Handler) deleteTransaction(c *gin.Context) {
	uid := c.GetInt64("uid")

	var req struct {
		ID int64 `json:"id"`
	}
	if err := c.BindJSON(&req); err != nil {
		_ = c.Error(apperror.BadRequestError(err, 4002, "invalid json body"))
		return
	}

	if err := h.service.DeleteTransaction(c.Request.Context(), uid, req.ID); err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
func (h *Handler) updateTransaction(c *gin.Context) {
	uid := c.GetInt64("uid")

	var t model.Transaction
	if err := c.BindJSON(&t); err != nil {
		_ = c.Error(apperror.BadRequestError(err, 4003, "invalid json body"))
		return
	}

	if err := h.service.EditTransaction(c.Request.Context(), uid, t); err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (h *Handler) authUser(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	var authData string
	if authHeader != "" {
		authParts := strings.Split(authHeader, " ")
		if len(authParts) != 2 {
			appErr := apperror.BadRequestError(errors.New("invalid auth header"), 4004, "Invalid authorization header format. Expected 'tma <data>'.")
			_ = c.Error(appErr)
			return
		}
		authData = authParts[1]
	} else {
		raw := c.Request.URL.RawQuery
		const tmaPrefix = "tma"
		idx := strings.Index(raw, tmaPrefix)
		if idx == -1 {
			authData = ""
		} else {
			authData = raw[idx+len(tmaPrefix)+1:]
		}
	}

	if authData == "" {
		err := apperror.BadRequestError(errors.New("missing auth header"), 4005, "Missing authorization header")
		_ = c.Error(err)
		return
	}

	uid, err := h.service.AuthUser(c.Request.Context(), authData)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.Set("uid", uid)
	c.Next()
}

### .\backend\manager\internal\router\router.go END ###

### .\backend\manager\internal\services\http.go BEGIN ###
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
### .\backend\manager\internal\services\http.go END ###

### .\backend\README.md BEGIN ###
backend on go

rename this idk

### .\backend\README.md END ###

### .\docker-compose.yml BEGIN ###
services:

  frontend:
    build:
      context: ./backend/frontend
      dockerfile: Dockerfile
    restart: unless-stopped
    ports:
      - 3000:3000
    networks:
      - cybergarden-net

  db:
    image: postgres:16-alpine
    environment:
      POSTGRES_USER: ${POSTGRES_USER}
      POSTGRES_PASSWORD: ${POSTGRES_PASSWORD}
      POSTGRES_DB: ${POSTGRES_DB}
    restart: unless-stopped
    volumes:
      - ./backend/database_save:/var/lib/postgresql/data
    expose:
      - 5432
    networks:
      - cybergarden-net
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ${POSTGRES_USER} -d ${POSTGRES_DB}"]
      interval: 5s
      timeout: 5s
      retries: 5

  database-service:
    build:
      context: ./backend/database 
      dockerfile: Dockerfile
    restart: unless-stopped
    expose:
      - 2012
    networks:
      - cybergarden-net
    depends_on:
      db:
        condition: service_healthy

  bot:
    build:
      context: ./backend/bot
      dockerfile: Dockerfile
    restart: unless-stopped
    volumes:
      - ./backend/bot/app.log:/app/app.log
    ports:
      - 8081:8081
    expose:
      - 2011
    networks:
      - cybergarden-net

  manager:
    build:
      context: ./backend/manager
      dockerfile: Dockerfile
    restart: unless-stopped
    ports:
      - "8080:8080"
    networks:
      - cybergarden-net
    depends_on:
      - bot
      - database-service
### .\docker-compose.yml END ###

### .\flatten-rust.exe BEGIN ###
[Binary file skipped: .\flatten-rust.exe]
### .\flatten-rust.exe END ###

### .\frontend\components.json BEGIN ###
{
  "$schema": "https://ui.shadcn.com/schema.json",
  "style": "new-york",
  "rsc": true,
  "tsx": true,
  "tailwind": {
    "config": "",
    "css": "src/app/globals.css",
    "baseColor": "neutral",
    "cssVariables": true,
    "prefix": ""
  },
  "iconLibrary": "lucide",
  "aliases": {
    "components": "@/components",
    "utils": "@/lib/utils",
    "ui": "@/components/ui",
    "lib": "@/lib",
    "hooks": "@/hooks"
  },
  "registries": {}
}

### .\frontend\components.json END ###

### .\frontend\next.config.ts BEGIN ###
import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  /* config options here */
};

export default nextConfig;

### .\frontend\next.config.ts END ###

### .\frontend\package.json BEGIN ###
{
  "name": "frontend",
  "version": "0.1.0",
  "private": true,
  "scripts": {
    "dev": "next dev",
    "build": "next build",
    "start": "next start"
  },
  "dependencies": {
    "@hookform/resolvers": "^5.2.2",
    "@radix-ui/react-label": "^2.1.8",
    "@radix-ui/react-slot": "^1.2.4",
    "@tanstack/react-query": "^5.90.12",
    "@tanstack/react-table": "^8.21.3",
    "class-variance-authority": "^0.7.1",
    "clsx": "^2.1.1",
    "date-fns": "^4.1.0",
    "lucide-react": "^0.556.0",
    "next": "16.0.7",
    "next-auth": "^4.24.13",
    "react": "19.2.0",
    "react-dom": "19.2.0",
    "react-hook-form": "^7.68.0",
    "recharts": "^3.5.1",
    "tailwind-merge": "^3.4.0",
    "zod": "^4.1.13"
  },
  "devDependencies": {
    "@tailwindcss/postcss": "^4",
    "@types/node": "^20",
    "@types/react": "^19",
    "@types/react-dom": "^19",
    "tailwindcss": "^4",
    "tw-animate-css": "^1.4.0",
    "typescript": "^5"
  }
}

### .\frontend\package.json END ###

### .\frontend\pnpm-lock.yaml BEGIN ###
lockfileVersion: '9.0'

settings:
  autoInstallPeers: true
  excludeLinksFromLockfile: false

importers:

  .:
    dependencies:
      '@hookform/resolvers':
        specifier: ^5.2.2
        version: 5.2.2(react-hook-form@7.68.0(react@19.2.0))
      '@radix-ui/react-label':
        specifier: ^2.1.8
        version: 2.1.8(@types/react-dom@19.2.3(@types/react@19.2.7))(@types/react@19.2.7)(react-dom@19.2.0(react@19.2.0))(react@19.2.0)
      '@radix-ui/react-slot':
        specifier: ^1.2.4
        version: 1.2.4(@types/react@19.2.7)(react@19.2.0)
      '@tanstack/react-query':
        specifier: ^5.90.12
        version: 5.90.12(react@19.2.0)
      '@tanstack/react-table':
        specifier: ^8.21.3
        version: 8.21.3(react-dom@19.2.0(react@19.2.0))(react@19.2.0)
      class-variance-authority:
        specifier: ^0.7.1
        version: 0.7.1
      clsx:
        specifier: ^2.1.1
        version: 2.1.1
      date-fns:
        specifier: ^4.1.0
        version: 4.1.0
      lucide-react:
        specifier: ^0.556.0
        version: 0.556.0(react@19.2.0)
      next:
        specifier: 16.0.7
        version: 16.0.7(react-dom@19.2.0(react@19.2.0))(react@19.2.0)
      next-auth:
        specifier: ^4.24.13
        version: 4.24.13(next@16.0.7(react-dom@19.2.0(react@19.2.0))(react@19.2.0))(react-dom@19.2.0(react@19.2.0))(react@19.2.0)
      react:
        specifier: 19.2.0
        version: 19.2.0
      react-dom:
        specifier: 19.2.0
        version: 19.2.0(react@19.2.0)
      react-hook-form:
        specifier: ^7.68.0
        version: 7.68.0(react@19.2.0)
      recharts:
        specifier: ^3.5.1
        version: 3.5.1(@types/react@19.2.7)(react-dom@19.2.0(react@19.2.0))(react-is@19.2.1)(react@19.2.0)(redux@5.0.1)
      tailwind-merge:
        specifier: ^3.4.0
        version: 3.4.0
      zod:
        specifier: ^4.1.13
        version: 4.1.13
    devDependencies:
      '@tailwindcss/postcss':
        specifier: ^4
        version: 4.1.17
      '@types/node':
        specifier: ^20
        version: 20.19.25
      '@types/react':
        specifier: ^19
        version: 19.2.7
      '@types/react-dom':
        specifier: ^19
        version: 19.2.3(@types/react@19.2.7)
      tailwindcss:
        specifier: ^4
        version: 4.1.17
      tw-animate-css:
        specifier: ^1.4.0
        version: 1.4.0
      typescript:
        specifier: ^5
        version: 5.9.3

packages:

  '@alloc/quick-lru@5.2.0':
    resolution: {integrity: sha512-UrcABB+4bUrFABwbluTIBErXwvbsU/V7TZWfmbgJfbkwiBuziS9gxdODUyuiecfdGQ85jglMW6juS3+z5TsKLw==}
    engines: {node: '>=10'}

  '@babel/runtime@7.28.4':
    resolution: {integrity: sha512-Q/N6JNWvIvPnLDvjlE1OUBLPQHH6l3CltCEsHIujp45zQUSSh8K+gHnaEX45yAT1nyngnINhvWtzN+Nb9D8RAQ==}
    engines: {node: '>=6.9.0'}

  '@emnapi/runtime@1.7.1':
    resolution: {integrity: sha512-PVtJr5CmLwYAU9PZDMITZoR5iAOShYREoR45EyyLrbntV50mdePTgUn4AmOw90Ifcj+x2kRjdzr1HP3RrNiHGA==}

  '@hookform/resolvers@5.2.2':
    resolution: {integrity: sha512-A/IxlMLShx3KjV/HeTcTfaMxdwy690+L/ZADoeaTltLx+CVuzkeVIPuybK3jrRfw7YZnmdKsVVHAlEPIAEUNlA==}
    peerDependencies:
      react-hook-form: ^7.55.0

  '@img/colour@1.0.0':
    resolution: {integrity: sha512-A5P/LfWGFSl6nsckYtjw9da+19jB8hkJ6ACTGcDfEJ0aE+l2n2El7dsVM7UVHZQ9s2lmYMWlrS21YLy2IR1LUw==}
    engines: {node: '>=18'}

  '@img/sharp-darwin-arm64@0.34.5':
    resolution: {integrity: sha512-imtQ3WMJXbMY4fxb/Ndp6HBTNVtWCUI0WdobyheGf5+ad6xX8VIDO8u2xE4qc/fr08CKG/7dDseFtn6M6g/r3w==}
    engines: {node: ^18.17.0 || ^20.3.0 || >=21.0.0}
    cpu: [arm64]
    os: [darwin]

  '@img/sharp-darwin-x64@0.34.5':
    resolution: {integrity: sha512-YNEFAF/4KQ/PeW0N+r+aVVsoIY0/qxxikF2SWdp+NRkmMB7y9LBZAVqQ4yhGCm/H3H270OSykqmQMKLBhBJDEw==}
    engines: {node: ^18.17.0 || ^20.3.0 || >=21.0.0}
    cpu: [x64]
    os: [darwin]

  '@img/sharp-libvips-darwin-arm64@1.2.4':
    resolution: {integrity: sha512-zqjjo7RatFfFoP0MkQ51jfuFZBnVE2pRiaydKJ1G/rHZvnsrHAOcQALIi9sA5co5xenQdTugCvtb1cuf78Vf4g==}
    cpu: [arm64]
    os: [darwin]

  '@img/sharp-libvips-darwin-x64@1.2.4':
    resolution: {integrity: sha512-1IOd5xfVhlGwX+zXv2N93k0yMONvUlANylbJw1eTah8K/Jtpi15KC+WSiaX/nBmbm2HxRM1gZ0nSdjSsrZbGKg==}
    cpu: [x64]
    os: [darwin]

  '@img/sharp-libvips-linux-arm64@1.2.4':
    resolution: {integrity: sha512-excjX8DfsIcJ10x1Kzr4RcWe1edC9PquDRRPx3YVCvQv+U5p7Yin2s32ftzikXojb1PIFc/9Mt28/y+iRklkrw==}
    cpu: [arm64]
    os: [linux]

  '@img/sharp-libvips-linux-arm@1.2.4':
    resolution: {integrity: sha512-bFI7xcKFELdiNCVov8e44Ia4u2byA+l3XtsAj+Q8tfCwO6BQ8iDojYdvoPMqsKDkuoOo+X6HZA0s0q11ANMQ8A==}
    cpu: [arm]
    os: [linux]

  '@img/sharp-libvips-linux-ppc64@1.2.4':
    resolution: {integrity: sha512-FMuvGijLDYG6lW+b/UvyilUWu5Ayu+3r2d1S8notiGCIyYU/76eig1UfMmkZ7vwgOrzKzlQbFSuQfgm7GYUPpA==}
    cpu: [ppc64]
    os: [linux]

  '@img/sharp-libvips-linux-riscv64@1.2.4':
    resolution: {integrity: sha512-oVDbcR4zUC0ce82teubSm+x6ETixtKZBh/qbREIOcI3cULzDyb18Sr/Wcyx7NRQeQzOiHTNbZFF1UwPS2scyGA==}
    cpu: [riscv64]
    os: [linux]

  '@img/sharp-libvips-linux-s390x@1.2.4':
    resolution: {integrity: sha512-qmp9VrzgPgMoGZyPvrQHqk02uyjA0/QrTO26Tqk6l4ZV0MPWIW6LTkqOIov+J1yEu7MbFQaDpwdwJKhbJvuRxQ==}
    cpu: [s390x]
    os: [linux]

  '@img/sharp-libvips-linux-x64@1.2.4':
    resolution: {integrity: sha512-tJxiiLsmHc9Ax1bz3oaOYBURTXGIRDODBqhveVHonrHJ9/+k89qbLl0bcJns+e4t4rvaNBxaEZsFtSfAdquPrw==}
    cpu: [x64]
    os: [linux]

  '@img/sharp-libvips-linuxmusl-arm64@1.2.4':
    resolution: {integrity: sha512-FVQHuwx1IIuNow9QAbYUzJ+En8KcVm9Lk5+uGUQJHaZmMECZmOlix9HnH7n1TRkXMS0pGxIJokIVB9SuqZGGXw==}
    cpu: [arm64]
    os: [linux]

  '@img/sharp-libvips-linuxmusl-x64@1.2.4':
    resolution: {integrity: sha512-+LpyBk7L44ZIXwz/VYfglaX/okxezESc6UxDSoyo2Ks6Jxc4Y7sGjpgU9s4PMgqgjj1gZCylTieNamqA1MF7Dg==}
    cpu: [x64]
    os: [linux]

  '@img/sharp-linux-arm64@0.34.5':
    resolution: {integrity: sha512-bKQzaJRY/bkPOXyKx5EVup7qkaojECG6NLYswgktOZjaXecSAeCWiZwwiFf3/Y+O1HrauiE3FVsGxFg8c24rZg==}
    engines: {node: ^18.17.0 || ^20.3.0 || >=21.0.0}
    cpu: [arm64]
    os: [linux]

  '@img/sharp-linux-arm@0.34.5':
    resolution: {integrity: sha512-9dLqsvwtg1uuXBGZKsxem9595+ujv0sJ6Vi8wcTANSFpwV/GONat5eCkzQo/1O6zRIkh0m/8+5BjrRr7jDUSZw==}
    engines: {node: ^18.17.0 || ^20.3.0 || >=21.0.0}
    cpu: [arm]
    os: [linux]

  '@img/sharp-linux-ppc64@0.34.5':
    resolution: {integrity: sha512-7zznwNaqW6YtsfrGGDA6BRkISKAAE1Jo0QdpNYXNMHu2+0dTrPflTLNkpc8l7MUP5M16ZJcUvysVWWrMefZquA==}
    engines: {node: ^18.17.0 || ^20.3.0 || >=21.0.0}
    cpu: [ppc64]
    os: [linux]

  '@img/sharp-linux-riscv64@0.34.5':
    resolution: {integrity: sha512-51gJuLPTKa7piYPaVs8GmByo7/U7/7TZOq+cnXJIHZKavIRHAP77e3N2HEl3dgiqdD/w0yUfiJnII77PuDDFdw==}
    engines: {node: ^18.17.0 || ^20.3.0 || >=21.0.0}
    cpu: [riscv64]
    os: [linux]

  '@img/sharp-linux-s390x@0.34.5':
    resolution: {integrity: sha512-nQtCk0PdKfho3eC5MrbQoigJ2gd1CgddUMkabUj+rBevs8tZ2cULOx46E7oyX+04WGfABgIwmMC0VqieTiR4jg==}
    engines: {node: ^18.17.0 || ^20.3.0 || >=21.0.0}
    cpu: [s390x]
    os: [linux]

  '@img/sharp-linux-x64@0.34.5':
    resolution: {integrity: sha512-MEzd8HPKxVxVenwAa+JRPwEC7QFjoPWuS5NZnBt6B3pu7EG2Ge0id1oLHZpPJdn3OQK+BQDiw9zStiHBTJQQQQ==}
    engines: {node: ^18.17.0 || ^20.3.0 || >=21.0.0}
    cpu: [x64]
    os: [linux]

  '@img/sharp-linuxmusl-arm64@0.34.5':
    resolution: {integrity: sha512-fprJR6GtRsMt6Kyfq44IsChVZeGN97gTD331weR1ex1c1rypDEABN6Tm2xa1wE6lYb5DdEnk03NZPqA7Id21yg==}
    engines: {node: ^18.17.0 || ^20.3.0 || >=21.0.0}
    cpu: [arm64]
    os: [linux]

  '@img/sharp-linuxmusl-x64@0.34.5':
    resolution: {integrity: sha512-Jg8wNT1MUzIvhBFxViqrEhWDGzqymo3sV7z7ZsaWbZNDLXRJZoRGrjulp60YYtV4wfY8VIKcWidjojlLcWrd8Q==}
    engines: {node: ^18.17.0 || ^20.3.0 || >=21.0.0}
    cpu: [x64]
    os: [linux]

  '@img/sharp-wasm32@0.34.5':
    resolution: {integrity: sha512-OdWTEiVkY2PHwqkbBI8frFxQQFekHaSSkUIJkwzclWZe64O1X4UlUjqqqLaPbUpMOQk6FBu/HtlGXNblIs0huw==}
    engines: {node: ^18.17.0 || ^20.3.0 || >=21.0.0}
    cpu: [wasm32]

  '@img/sharp-win32-arm64@0.34.5':
    resolution: {integrity: sha512-WQ3AgWCWYSb2yt+IG8mnC6Jdk9Whs7O0gxphblsLvdhSpSTtmu69ZG1Gkb6NuvxsNACwiPV6cNSZNzt0KPsw7g==}
    engines: {node: ^18.17.0 || ^20.3.0 || >=21.0.0}
    cpu: [arm64]
    os: [win32]

  '@img/sharp-win32-ia32@0.34.5':
    resolution: {integrity: sha512-FV9m/7NmeCmSHDD5j4+4pNI8Cp3aW+JvLoXcTUo0IqyjSfAZJ8dIUmijx1qaJsIiU+Hosw6xM5KijAWRJCSgNg==}
    engines: {node: ^18.17.0 || ^20.3.0 || >=21.0.0}
    cpu: [ia32]
    os: [win32]

  '@img/sharp-win32-x64@0.34.5':
    resolution: {integrity: sha512-+29YMsqY2/9eFEiW93eqWnuLcWcufowXewwSNIT6UwZdUUCrM3oFjMWH/Z6/TMmb4hlFenmfAVbpWeup2jryCw==}
    engines: {node: ^18.17.0 || ^20.3.0 || >=21.0.0}
    cpu: [x64]
    os: [win32]

  '@jridgewell/gen-mapping@0.3.13':
    resolution: {integrity: sha512-2kkt/7niJ6MgEPxF0bYdQ6etZaA+fQvDcLKckhy1yIQOzaoKjBBjSj63/aLVjYE3qhRt5dvM+uUyfCg6UKCBbA==}

  '@jridgewell/remapping@2.3.5':
    resolution: {integrity: sha512-LI9u/+laYG4Ds1TDKSJW2YPrIlcVYOwi2fUC6xB43lueCjgxV4lffOCZCtYFiH6TNOX+tQKXx97T4IKHbhyHEQ==}

  '@jridgewell/resolve-uri@3.1.2':
    resolution: {integrity: sha512-bRISgCIjP20/tbWSPWMEi54QVPRZExkuD9lJL+UIxUKtwVJA8wW1Trb1jMs1RFXo1CBTNZ/5hpC9QvmKWdopKw==}
    engines: {node: '>=6.0.0'}

  '@jridgewell/sourcemap-codec@1.5.5':
    resolution: {integrity: sha512-cYQ9310grqxueWbl+WuIUIaiUaDcj7WOq5fVhEljNVgRfOUhY9fy2zTvfoqWsnebh8Sl70VScFbICvJnLKB0Og==}

  '@jridgewell/trace-mapping@0.3.31':
    resolution: {integrity: sha512-zzNR+SdQSDJzc8joaeP8QQoCQr8NuYx2dIIytl1QeBEZHJ9uW6hebsrYgbz8hJwUQao3TWCMtmfV8Nu1twOLAw==}

  '@next/env@16.0.7':
    resolution: {integrity: sha512-gpaNgUh5nftFKRkRQGnVi5dpcYSKGcZZkQffZ172OrG/XkrnS7UBTQ648YY+8ME92cC4IojpI2LqTC8sTDhAaw==}

  '@next/swc-darwin-arm64@16.0.7':
    resolution: {integrity: sha512-LlDtCYOEj/rfSnEn/Idi+j1QKHxY9BJFmxx7108A6D8K0SB+bNgfYQATPk/4LqOl4C0Wo3LACg2ie6s7xqMpJg==}
    engines: {node: '>= 10'}
    cpu: [arm64]
    os: [darwin]

  '@next/swc-darwin-x64@16.0.7':
    resolution: {integrity: sha512-rtZ7BhnVvO1ICf3QzfW9H3aPz7GhBrnSIMZyr4Qy6boXF0b5E3QLs+cvJmg3PsTCG2M1PBoC+DANUi4wCOKXpA==}
    engines: {node: '>= 10'}
    cpu: [x64]
    os: [darwin]

  '@next/swc-linux-arm64-gnu@16.0.7':
    resolution: {integrity: sha512-mloD5WcPIeIeeZqAIP5c2kdaTa6StwP4/2EGy1mUw8HiexSHGK/jcM7lFuS3u3i2zn+xH9+wXJs6njO7VrAqww==}
    engines: {node: '>= 10'}
    cpu: [arm64]
    os: [linux]

  '@next/swc-linux-arm64-musl@16.0.7':
    resolution: {integrity: sha512-+ksWNrZrthisXuo9gd1XnjHRowCbMtl/YgMpbRvFeDEqEBd523YHPWpBuDjomod88U8Xliw5DHhekBC3EOOd9g==}
    engines: {node: '>= 10'}
    cpu: [arm64]
    os: [linux]

  '@next/swc-linux-x64-gnu@16.0.7':
    resolution: {integrity: sha512-4WtJU5cRDxpEE44Ana2Xro1284hnyVpBb62lIpU5k85D8xXxatT+rXxBgPkc7C1XwkZMWpK5rXLXTh9PFipWsA==}
    engines: {node: '>= 10'}
    cpu: [x64]
    os: [linux]

  '@next/swc-linux-x64-musl@16.0.7':
    resolution: {integrity: sha512-HYlhqIP6kBPXalW2dbMTSuB4+8fe+j9juyxwfMwCe9kQPPeiyFn7NMjNfoFOfJ2eXkeQsoUGXg+O2SE3m4Qg2w==}
    engines: {node: '>= 10'}
    cpu: [x64]
    os: [linux]

  '@next/swc-win32-arm64-msvc@16.0.7':
    resolution: {integrity: sha512-EviG+43iOoBRZg9deGauXExjRphhuYmIOJ12b9sAPy0eQ6iwcPxfED2asb/s2/yiLYOdm37kPaiZu8uXSYPs0Q==}
    engines: {node: '>= 10'}
    cpu: [arm64]
    os: [win32]

  '@next/swc-win32-x64-msvc@16.0.7':
    resolution: {integrity: sha512-gniPjy55zp5Eg0896qSrf3yB1dw4F/3s8VK1ephdsZZ129j2n6e1WqCbE2YgcKhW9hPB9TVZENugquWJD5x0ug==}
    engines: {node: '>= 10'}
    cpu: [x64]
    os: [win32]

  '@panva/hkdf@1.2.1':
    resolution: {integrity: sha512-6oclG6Y3PiDFcoyk8srjLfVKyMfVCKJ27JwNPViuXziFpmdz+MZnZN/aKY0JGXgYuO/VghU0jcOAZgWXZ1Dmrw==}

  '@radix-ui/react-compose-refs@1.1.2':
    resolution: {integrity: sha512-z4eqJvfiNnFMHIIvXP3CY57y2WJs5g2v3X0zm9mEJkrkNv4rDxu+sg9Jh8EkXyeqBkB7SOcboo9dMVqhyrACIg==}
    peerDependencies:
      '@types/react': '*'
      react: ^16.8 || ^17.0 || ^18.0 || ^19.0 || ^19.0.0-rc
    peerDependenciesMeta:
      '@types/react':
        optional: true

  '@radix-ui/react-label@2.1.8':
    resolution: {integrity: sha512-FmXs37I6hSBVDlO4y764TNz1rLgKwjJMQ0EGte6F3Cb3f4bIuHB/iLa/8I9VKkmOy+gNHq8rql3j686ACVV21A==}
    peerDependencies:
      '@types/react': '*'
      '@types/react-dom': '*'
      react: ^16.8 || ^17.0 || ^18.0 || ^19.0 || ^19.0.0-rc
      react-dom: ^16.8 || ^17.0 || ^18.0 || ^19.0 || ^19.0.0-rc
    peerDependenciesMeta:
      '@types/react':
        optional: true
      '@types/react-dom':
        optional: true

  '@radix-ui/react-primitive@2.1.4':
    resolution: {integrity: sha512-9hQc4+GNVtJAIEPEqlYqW5RiYdrr8ea5XQ0ZOnD6fgru+83kqT15mq2OCcbe8KnjRZl5vF3ks69AKz3kh1jrhg==}
    peerDependencies:
      '@types/react': '*'
      '@types/react-dom': '*'
      react: ^16.8 || ^17.0 || ^18.0 || ^19.0 || ^19.0.0-rc
      react-dom: ^16.8 || ^17.0 || ^18.0 || ^19.0 || ^19.0.0-rc
    peerDependenciesMeta:
      '@types/react':
        optional: true
      '@types/react-dom':
        optional: true

  '@radix-ui/react-slot@1.2.4':
    resolution: {integrity: sha512-Jl+bCv8HxKnlTLVrcDE8zTMJ09R9/ukw4qBs/oZClOfoQk/cOTbDn+NceXfV7j09YPVQUryJPHurafcSg6EVKA==}
    peerDependencies:
      '@types/react': '*'
      react: ^16.8 || ^17.0 || ^18.0 || ^19.0 || ^19.0.0-rc
    peerDependenciesMeta:
      '@types/react':
        optional: true

  '@reduxjs/toolkit@2.11.0':
    resolution: {integrity: sha512-hBjYg0aaRL1O2Z0IqWhnTLytnjDIxekmRxm1snsHjHaKVmIF1HiImWqsq+PuEbn6zdMlkIj9WofK1vR8jjx+Xw==}
    peerDependencies:
      react: ^16.9.0 || ^17.0.0 || ^18 || ^19
      react-redux: ^7.2.1 || ^8.1.3 || ^9.0.0
    peerDependenciesMeta:
      react:
        optional: true
      react-redux:
        optional: true

  '@standard-schema/spec@1.0.0':
    resolution: {integrity: sha512-m2bOd0f2RT9k8QJx1JN85cZYyH1RqFBdlwtkSlf4tBDYLCiiZnv1fIIwacK6cqwXavOydf0NPToMQgpKq+dVlA==}

  '@standard-schema/utils@0.3.0':
    resolution: {integrity: sha512-e7Mew686owMaPJVNNLs55PUvgz371nKgwsc4vxE49zsODpJEnxgxRo2y/OKrqueavXgZNMDVj3DdHFlaSAeU8g==}

  '@swc/helpers@0.5.15':
    resolution: {integrity: sha512-JQ5TuMi45Owi4/BIMAJBoSQoOJu12oOk/gADqlcUL9JEdHB8vyjUSsxqeNXnmXHjYKMi2WcYtezGEEhqUI/E2g==}

  '@tailwindcss/node@4.1.17':
    resolution: {integrity: sha512-csIkHIgLb3JisEFQ0vxr2Y57GUNYh447C8xzwj89U/8fdW8LhProdxvnVH6U8M2Y73QKiTIH+LWbK3V2BBZsAg==}

  '@tailwindcss/oxide-android-arm64@4.1.17':
    resolution: {integrity: sha512-BMqpkJHgOZ5z78qqiGE6ZIRExyaHyuxjgrJ6eBO5+hfrfGkuya0lYfw8fRHG77gdTjWkNWEEm+qeG2cDMxArLQ==}
    engines: {node: '>= 10'}
    cpu: [arm64]
    os: [android]

  '@tailwindcss/oxide-darwin-arm64@4.1.17':
    resolution: {integrity: sha512-EquyumkQweUBNk1zGEU/wfZo2qkp/nQKRZM8bUYO0J+Lums5+wl2CcG1f9BgAjn/u9pJzdYddHWBiFXJTcxmOg==}
    engines: {node: '>= 10'}
    cpu: [arm64]
    os: [darwin]

  '@tailwindcss/oxide-darwin-x64@4.1.17':
    resolution: {integrity: sha512-gdhEPLzke2Pog8s12oADwYu0IAw04Y2tlmgVzIN0+046ytcgx8uZmCzEg4VcQh+AHKiS7xaL8kGo/QTiNEGRog==}
    engines: {node: '>= 10'}
    cpu: [x64]
    os: [darwin]

  '@tailwindcss/oxide-freebsd-x64@4.1.17':
    resolution: {integrity: sha512-hxGS81KskMxML9DXsaXT1H0DyA+ZBIbyG/sSAjWNe2EDl7TkPOBI42GBV3u38itzGUOmFfCzk1iAjDXds8Oh0g==}
    engines: {node: '>= 10'}
    cpu: [x64]
    os: [freebsd]

  '@tailwindcss/oxide-linux-arm-gnueabihf@4.1.17':
    resolution: {integrity: sha512-k7jWk5E3ldAdw0cNglhjSgv501u7yrMf8oeZ0cElhxU6Y2o7f8yqelOp3fhf7evjIS6ujTI3U8pKUXV2I4iXHQ==}
    engines: {node: '>= 10'}
    cpu: [arm]
    os: [linux]

  '@tailwindcss/oxide-linux-arm64-gnu@4.1.17':
    resolution: {integrity: sha512-HVDOm/mxK6+TbARwdW17WrgDYEGzmoYayrCgmLEw7FxTPLcp/glBisuyWkFz/jb7ZfiAXAXUACfyItn+nTgsdQ==}
    engines: {node: '>= 10'}
    cpu: [arm64]
    os: [linux]

  '@tailwindcss/oxide-linux-arm64-musl@4.1.17':
    resolution: {integrity: sha512-HvZLfGr42i5anKtIeQzxdkw/wPqIbpeZqe7vd3V9vI3RQxe3xU1fLjss0TjyhxWcBaipk7NYwSrwTwK1hJARMg==}
    engines: {node: '>= 10'}
    cpu: [arm64]
    os: [linux]

  '@tailwindcss/oxide-linux-x64-gnu@4.1.17':
    resolution: {integrity: sha512-M3XZuORCGB7VPOEDH+nzpJ21XPvK5PyjlkSFkFziNHGLc5d6g3di2McAAblmaSUNl8IOmzYwLx9NsE7bplNkwQ==}
    engines: {node: '>= 10'}
    cpu: [x64]
    os: [linux]

  '@tailwindcss/oxide-linux-x64-musl@4.1.17':
    resolution: {integrity: sha512-k7f+pf9eXLEey4pBlw+8dgfJHY4PZ5qOUFDyNf7SI6lHjQ9Zt7+NcscjpwdCEbYi6FI5c2KDTDWyf2iHcCSyyQ==}
    engines: {node: '>= 10'}
    cpu: [x64]
    os: [linux]

  '@tailwindcss/oxide-wasm32-wasi@4.1.17':
    resolution: {integrity: sha512-cEytGqSSoy7zK4JRWiTCx43FsKP/zGr0CsuMawhH67ONlH+T79VteQeJQRO/X7L0juEUA8ZyuYikcRBf0vsxhg==}
    engines: {node: '>=14.0.0'}
    cpu: [wasm32]
    bundledDependencies:
      - '@napi-rs/wasm-runtime'
      - '@emnapi/core'
      - '@emnapi/runtime'
      - '@tybys/wasm-util'
      - '@emnapi/wasi-threads'
      - tslib

  '@tailwindcss/oxide-win32-arm64-msvc@4.1.17':
    resolution: {integrity: sha512-JU5AHr7gKbZlOGvMdb4722/0aYbU+tN6lv1kONx0JK2cGsh7g148zVWLM0IKR3NeKLv+L90chBVYcJ8uJWbC9A==}
    engines: {node: '>= 10'}
    cpu: [arm64]
    os: [win32]

  '@tailwindcss/oxide-win32-x64-msvc@4.1.17':
    resolution: {integrity: sha512-SKWM4waLuqx0IH+FMDUw6R66Hu4OuTALFgnleKbqhgGU30DY20NORZMZUKgLRjQXNN2TLzKvh48QXTig4h4bGw==}
    engines: {node: '>= 10'}
    cpu: [x64]
    os: [win32]

  '@tailwindcss/oxide@4.1.17':
    resolution: {integrity: sha512-F0F7d01fmkQhsTjXezGBLdrl1KresJTcI3DB8EkScCldyKp3Msz4hub4uyYaVnk88BAS1g5DQjjF6F5qczheLA==}
    engines: {node: '>= 10'}

  '@tailwindcss/postcss@4.1.17':
    resolution: {integrity: sha512-+nKl9N9mN5uJ+M7dBOOCzINw94MPstNR/GtIhz1fpZysxL/4a+No64jCBD6CPN+bIHWFx3KWuu8XJRrj/572Dw==}

  '@tanstack/query-core@5.90.12':
    resolution: {integrity: sha512-T1/8t5DhV/SisWjDnaiU2drl6ySvsHj1bHBCWNXd+/T+Hh1cf6JodyEYMd5sgwm+b/mETT4EV3H+zCVczCU5hg==}

  '@tanstack/react-query@5.90.12':
    resolution: {integrity: sha512-graRZspg7EoEaw0a8faiUASCyJrqjKPdqJ9EwuDRUF9mEYJ1YPczI9H+/agJ0mOJkPCJDk0lsz5QTrLZ/jQ2rg==}
    peerDependencies:
      react: ^18 || ^19

  '@tanstack/react-table@8.21.3':
    resolution: {integrity: sha512-5nNMTSETP4ykGegmVkhjcS8tTLW6Vl4axfEGQN3v0zdHYbK4UfoqfPChclTrJ4EoK9QynqAu9oUf8VEmrpZ5Ww==}
    engines: {node: '>=12'}
    peerDependencies:
      react: '>=16.8'
      react-dom: '>=16.8'

  '@tanstack/table-core@8.21.3':
    resolution: {integrity: sha512-ldZXEhOBb8Is7xLs01fR3YEc3DERiz5silj8tnGkFZytt1abEvl/GhUmCE0PMLaMPTa3Jk4HbKmRlHmu+gCftg==}
    engines: {node: '>=12'}

  '@types/d3-array@3.2.2':
    resolution: {integrity: sha512-hOLWVbm7uRza0BYXpIIW5pxfrKe0W+D5lrFiAEYR+pb6w3N2SwSMaJbXdUfSEv+dT4MfHBLtn5js0LAWaO6otw==}

  '@types/d3-color@3.1.3':
    resolution: {integrity: sha512-iO90scth9WAbmgv7ogoq57O9YpKmFBbmoEoCHDB2xMBY0+/KVrqAaCDyCE16dUspeOvIxFFRI+0sEtqDqy2b4A==}

  '@types/d3-ease@3.0.2':
    resolution: {integrity: sha512-NcV1JjO5oDzoK26oMzbILE6HW7uVXOHLQvHshBUW4UMdZGfiY6v5BeQwh9a9tCzv+CeefZQHJt5SRgK154RtiA==}

  '@types/d3-interpolate@3.0.4':
    resolution: {integrity: sha512-mgLPETlrpVV1YRJIglr4Ez47g7Yxjl1lj7YKsiMCb27VJH9W8NVM6Bb9d8kkpG/uAQS5AmbA48q2IAolKKo1MA==}

  '@types/d3-path@3.1.1':
    resolution: {integrity: sha512-VMZBYyQvbGmWyWVea0EHs/BwLgxc+MKi1zLDCONksozI4YJMcTt8ZEuIR4Sb1MMTE8MMW49v0IwI5+b7RmfWlg==}

  '@types/d3-scale@4.0.9':
    resolution: {integrity: sha512-dLmtwB8zkAeO/juAMfnV+sItKjlsw2lKdZVVy6LRr0cBmegxSABiLEpGVmSJJ8O08i4+sGR6qQtb6WtuwJdvVw==}

  '@types/d3-shape@3.1.7':
    resolution: {integrity: sha512-VLvUQ33C+3J+8p+Daf+nYSOsjB4GXp19/S/aGo60m9h1v6XaxjiT82lKVWJCfzhtuZ3yD7i/TPeC/fuKLLOSmg==}

  '@types/d3-time@3.0.4':
    resolution: {integrity: sha512-yuzZug1nkAAaBlBBikKZTgzCeA+k1uy4ZFwWANOfKw5z5LRhV0gNA7gNkKm7HoK+HRN0wX3EkxGk0fpbWhmB7g==}

  '@types/d3-timer@3.0.2':
    resolution: {integrity: sha512-Ps3T8E8dZDam6fUyNiMkekK3XUsaUEik+idO9/YjPtfj2qruF8tFBXS7XhtE4iIXBLxhmLjP3SXpLhVf21I9Lw==}

  '@types/node@20.19.25':
    resolution: {integrity: sha512-ZsJzA5thDQMSQO788d7IocwwQbI8B5OPzmqNvpf3NY/+MHDAS759Wo0gd2WQeXYt5AAAQjzcrTVC6SKCuYgoCQ==}

  '@types/react-dom@19.2.3':
    resolution: {integrity: sha512-jp2L/eY6fn+KgVVQAOqYItbF0VY/YApe5Mz2F0aykSO8gx31bYCZyvSeYxCHKvzHG5eZjc+zyaS5BrBWya2+kQ==}
    peerDependencies:
      '@types/react': ^19.2.0

  '@types/react@19.2.7':
    resolution: {integrity: sha512-MWtvHrGZLFttgeEj28VXHxpmwYbor/ATPYbBfSFZEIRK0ecCFLl2Qo55z52Hss+UV9CRN7trSeq1zbgx7YDWWg==}

  '@types/use-sync-external-store@0.0.6':
    resolution: {integrity: sha512-zFDAD+tlpf2r4asuHEj0XH6pY6i0g5NeAHPn+15wk3BV6JA69eERFXC1gyGThDkVa1zCyKr5jox1+2LbV/AMLg==}

  caniuse-lite@1.0.30001759:
    resolution: {integrity: sha512-Pzfx9fOKoKvevQf8oCXoyNRQ5QyxJj+3O0Rqx2V5oxT61KGx8+n6hV/IUyJeifUci2clnmmKVpvtiqRzgiWjSw==}

  class-variance-authority@0.7.1:
    resolution: {integrity: sha512-Ka+9Trutv7G8M6WT6SeiRWz792K5qEqIGEGzXKhAE6xOWAY6pPH8U+9IY3oCMv6kqTmLsv7Xh/2w2RigkePMsg==}

  client-only@0.0.1:
    resolution: {integrity: sha512-IV3Ou0jSMzZrd3pZ48nLkT9DA7Ag1pnPzaiQhpW7c3RbcqqzvzzVu+L8gfqMp/8IM2MQtSiqaCxrrcfu8I8rMA==}

  clsx@2.1.1:
    resolution: {integrity: sha512-eYm0QWBtUrBWZWG0d386OGAw16Z995PiOVo2B7bjWSbHedGl5e0ZWaq65kOGgUSNesEIDkB9ISbTg/JK9dhCZA==}
    engines: {node: '>=6'}

  cookie@0.7.2:
    resolution: {integrity: sha512-yki5XnKuf750l50uGTllt6kKILY4nQ1eNIQatoXEByZ5dWgnKqbnqmTrBE5B4N7lrMJKQ2ytWMiTO2o0v6Ew/w==}
    engines: {node: '>= 0.6'}

  csstype@3.2.3:
    resolution: {integrity: sha512-z1HGKcYy2xA8AGQfwrn0PAy+PB7X/GSj3UVJW9qKyn43xWa+gl5nXmU4qqLMRzWVLFC8KusUX8T/0kCiOYpAIQ==}

  d3-array@3.2.4:
    resolution: {integrity: sha512-tdQAmyA18i4J7wprpYq8ClcxZy3SC31QMeByyCFyRt7BVHdREQZ5lpzoe5mFEYZUWe+oq8HBvk9JjpibyEV4Jg==}
    engines: {node: '>=12'}

  d3-color@3.1.0:
    resolution: {integrity: sha512-zg/chbXyeBtMQ1LbD/WSoW2DpC3I0mpmPdW+ynRTj/x2DAWYrIY7qeZIHidozwV24m4iavr15lNwIwLxRmOxhA==}
    engines: {node: '>=12'}

  d3-ease@3.0.1:
    resolution: {integrity: sha512-wR/XK3D3XcLIZwpbvQwQ5fK+8Ykds1ip7A2Txe0yxncXSdq1L9skcG7blcedkOX+ZcgxGAmLX1FrRGbADwzi0w==}
    engines: {node: '>=12'}

  d3-format@3.1.0:
    resolution: {integrity: sha512-YyUI6AEuY/Wpt8KWLgZHsIU86atmikuoOmCfommt0LYHiQSPjvX2AcFc38PX0CBpr2RCyZhjex+NS/LPOv6YqA==}
    engines: {node: '>=12'}

  d3-interpolate@3.0.1:
    resolution: {integrity: sha512-3bYs1rOD33uo8aqJfKP3JWPAibgw8Zm2+L9vBKEHJ2Rg+viTR7o5Mmv5mZcieN+FRYaAOWX5SJATX6k1PWz72g==}
    engines: {node: '>=12'}

  d3-path@3.1.0:
    resolution: {integrity: sha512-p3KP5HCf/bvjBSSKuXid6Zqijx7wIfNW+J/maPs+iwR35at5JCbLUT0LzF1cnjbCHWhqzQTIN2Jpe8pRebIEFQ==}
    engines: {node: '>=12'}

  d3-scale@4.0.2:
    resolution: {integrity: sha512-GZW464g1SH7ag3Y7hXjf8RoUuAFIqklOAq3MRl4OaWabTFJY9PN/E1YklhXLh+OQ3fM9yS2nOkCoS+WLZ6kvxQ==}
    engines: {node: '>=12'}

  d3-shape@3.2.0:
    resolution: {integrity: sha512-SaLBuwGm3MOViRq2ABk3eLoxwZELpH6zhl3FbAoJ7Vm1gofKx6El1Ib5z23NUEhF9AsGl7y+dzLe5Cw2AArGTA==}
    engines: {node: '>=12'}

  d3-time-format@4.1.0:
    resolution: {integrity: sha512-dJxPBlzC7NugB2PDLwo9Q8JiTR3M3e4/XANkreKSUxF8vvXKqm1Yfq4Q5dl8budlunRVlUUaDUgFt7eA8D6NLg==}
    engines: {node: '>=12'}

  d3-time@3.1.0:
    resolution: {integrity: sha512-VqKjzBLejbSMT4IgbmVgDjpkYrNWUYJnbCGo874u7MMKIWsILRX+OpX/gTk8MqjpT1A/c6HY2dCA77ZN0lkQ2Q==}
    engines: {node: '>=12'}

  d3-timer@3.0.1:
    resolution: {integrity: sha512-ndfJ/JxxMd3nw31uyKoY2naivF+r29V+Lc0svZxe1JvvIRmi8hUsrMvdOwgS1o6uBHmiz91geQ0ylPP0aj1VUA==}
    engines: {node: '>=12'}

  date-fns@4.1.0:
    resolution: {integrity: sha512-Ukq0owbQXxa/U3EGtsdVBkR1w7KOQ5gIBqdH2hkvknzZPYvBxb/aa6E8L7tmjFtkwZBu3UXBbjIgPo/Ez4xaNg==}

  decimal.js-light@2.5.1:
    resolution: {integrity: sha512-qIMFpTMZmny+MMIitAB6D7iVPEorVw6YQRWkvarTkT4tBeSLLiHzcwj6q0MmYSFCiVpiqPJTJEYIrpcPzVEIvg==}

  detect-libc@2.1.2:
    resolution: {integrity: sha512-Btj2BOOO83o3WyH59e8MgXsxEQVcarkUOpEYrubB0urwnN10yQ364rsiByU11nZlqWYZm05i/of7io4mzihBtQ==}
    engines: {node: '>=8'}

  enhanced-resolve@5.18.3:
    resolution: {integrity: sha512-d4lC8xfavMeBjzGr2vECC3fsGXziXZQyJxD868h2M/mBI3PwAuODxAkLkq5HYuvrPYcUtiLzsTo8U3PgX3Ocww==}
    engines: {node: '>=10.13.0'}

  es-toolkit@1.42.0:
    resolution: {integrity: sha512-SLHIyY7VfDJBM8clz4+T2oquwTQxEzu263AyhVK4jREOAwJ+8eebaa4wM3nlvnAqhDrMm2EsA6hWHaQsMPQ1nA==}

  eventemitter3@5.0.1:
    resolution: {integrity: sha512-GWkBvjiSZK87ELrYOSESUYeVIc9mvLLf/nXalMOS5dYrgZq9o5OVkbZAVM06CVxYsCwH9BDZFPlQTlPA1j4ahA==}

  graceful-fs@4.2.11:
    resolution: {integrity: sha512-RbJ5/jmFcNNCcDV5o9eTnBLJ/HszWV0P73bc+Ff4nS/rJj+YaS6IGyiOL0VoBYX+l1Wrl3k63h/KrH+nhJ0XvQ==}

  immer@10.2.0:
    resolution: {integrity: sha512-d/+XTN3zfODyjr89gM3mPq1WNX2B8pYsu7eORitdwyA2sBubnTl3laYlBk4sXY5FUa5qTZGBDPJICVbvqzjlbw==}

  immer@11.0.1:
    resolution: {integrity: sha512-naDCyggtcBWANtIrjQEajhhBEuL9b0Zg4zmlWK2CzS6xCWSE39/vvf4LqnMjUAWHBhot4m9MHCM/Z+mfWhUkiA==}

  internmap@2.0.3:
    resolution: {integrity: sha512-5Hh7Y1wQbvY5ooGgPbDaL5iYLAPzMTUrjMulskHLH6wnv/A+1q5rgEaiuqEjB+oxGXIVZs1FF+R/KPN3ZSQYYg==}
    engines: {node: '>=12'}

  jiti@2.6.1:
    resolution: {integrity: sha512-ekilCSN1jwRvIbgeg/57YFh8qQDNbwDb9xT/qu2DAHbFFZUicIl4ygVaAvzveMhMVr3LnpSKTNnwt8PoOfmKhQ==}
    hasBin: true

  jose@4.15.9:
    resolution: {integrity: sha512-1vUQX+IdDMVPj4k8kOxgUqlcK518yluMuGZwqlr44FS1ppZB/5GWh4rZG89erpOBOJjU/OBsnCVFfapsRz6nEA==}

  lightningcss-android-arm64@1.30.2:
    resolution: {integrity: sha512-BH9sEdOCahSgmkVhBLeU7Hc9DWeZ1Eb6wNS6Da8igvUwAe0sqROHddIlvU06q3WyXVEOYDZ6ykBZQnjTbmo4+A==}
    engines: {node: '>= 12.0.0'}
    cpu: [arm64]
    os: [android]

  lightningcss-darwin-arm64@1.30.2:
    resolution: {integrity: sha512-ylTcDJBN3Hp21TdhRT5zBOIi73P6/W0qwvlFEk22fkdXchtNTOU4Qc37SkzV+EKYxLouZ6M4LG9NfZ1qkhhBWA==}
    engines: {node: '>= 12.0.0'}
    cpu: [arm64]
    os: [darwin]

  lightningcss-darwin-x64@1.30.2:
    resolution: {integrity: sha512-oBZgKchomuDYxr7ilwLcyms6BCyLn0z8J0+ZZmfpjwg9fRVZIR5/GMXd7r9RH94iDhld3UmSjBM6nXWM2TfZTQ==}
    engines: {node: '>= 12.0.0'}
    cpu: [x64]
    os: [darwin]

  lightningcss-freebsd-x64@1.30.2:
    resolution: {integrity: sha512-c2bH6xTrf4BDpK8MoGG4Bd6zAMZDAXS569UxCAGcA7IKbHNMlhGQ89eRmvpIUGfKWNVdbhSbkQaWhEoMGmGslA==}
    engines: {node: '>= 12.0.0'}
    cpu: [x64]
    os: [freebsd]

  lightningcss-linux-arm-gnueabihf@1.30.2:
    resolution: {integrity: sha512-eVdpxh4wYcm0PofJIZVuYuLiqBIakQ9uFZmipf6LF/HRj5Bgm0eb3qL/mr1smyXIS1twwOxNWndd8z0E374hiA==}
    engines: {node: '>= 12.0.0'}
    cpu: [arm]
    os: [linux]

  lightningcss-linux-arm64-gnu@1.30.2:
    resolution: {integrity: sha512-UK65WJAbwIJbiBFXpxrbTNArtfuznvxAJw4Q2ZGlU8kPeDIWEX1dg3rn2veBVUylA2Ezg89ktszWbaQnxD/e3A==}
    engines: {node: '>= 12.0.0'}
    cpu: [arm64]
    os: [linux]

  lightningcss-linux-arm64-musl@1.30.2:
    resolution: {integrity: sha512-5Vh9dGeblpTxWHpOx8iauV02popZDsCYMPIgiuw97OJ5uaDsL86cnqSFs5LZkG3ghHoX5isLgWzMs+eD1YzrnA==}
    engines: {node: '>= 12.0.0'}
    cpu: [arm64]
    os: [linux]

  lightningcss-linux-x64-gnu@1.30.2:
    resolution: {integrity: sha512-Cfd46gdmj1vQ+lR6VRTTadNHu6ALuw2pKR9lYq4FnhvgBc4zWY1EtZcAc6EffShbb1MFrIPfLDXD6Xprbnni4w==}
    engines: {node: '>= 12.0.0'}
    cpu: [x64]
    os: [linux]

  lightningcss-linux-x64-musl@1.30.2:
    resolution: {integrity: sha512-XJaLUUFXb6/QG2lGIW6aIk6jKdtjtcffUT0NKvIqhSBY3hh9Ch+1LCeH80dR9q9LBjG3ewbDjnumefsLsP6aiA==}
    engines: {node: '>= 12.0.0'}
    cpu: [x64]
    os: [linux]

  lightningcss-win32-arm64-msvc@1.30.2:
    resolution: {integrity: sha512-FZn+vaj7zLv//D/192WFFVA0RgHawIcHqLX9xuWiQt7P0PtdFEVaxgF9rjM/IRYHQXNnk61/H/gb2Ei+kUQ4xQ==}
    engines: {node: '>= 12.0.0'}
    cpu: [arm64]
    os: [win32]

  lightningcss-win32-x64-msvc@1.30.2:
    resolution: {integrity: sha512-5g1yc73p+iAkid5phb4oVFMB45417DkRevRbt/El/gKXJk4jid+vPFF/AXbxn05Aky8PapwzZrdJShv5C0avjw==}
    engines: {node: '>= 12.0.0'}
    cpu: [x64]
    os: [win32]

  lightningcss@1.30.2:
    resolution: {integrity: sha512-utfs7Pr5uJyyvDETitgsaqSyjCb2qNRAtuqUeWIAKztsOYdcACf2KtARYXg2pSvhkt+9NfoaNY7fxjl6nuMjIQ==}
    engines: {node: '>= 12.0.0'}

  lru-cache@6.0.0:
    resolution: {integrity: sha512-Jo6dJ04CmSjuznwJSS3pUeWmd/H0ffTlkXXgwZi+eq1UCmqQwCh+eLsYOYCwY991i2Fah4h1BEMCx4qThGbsiA==}
    engines: {node: '>=10'}

  lucide-react@0.556.0:
    resolution: {integrity: sha512-iOb8dRk7kLaYBZhR2VlV1CeJGxChBgUthpSP8wom9jfj79qovgG6qcSdiy6vkoREKPnbUYzJsCn4o4PtG3Iy+A==}
    peerDependencies:
      react: ^16.5.1 || ^17.0.0 || ^18.0.0 || ^19.0.0

  magic-string@0.30.21:
    resolution: {integrity: sha512-vd2F4YUyEXKGcLHoq+TEyCjxueSeHnFxyyjNp80yg0XV4vUhnDer/lvvlqM/arB5bXQN5K2/3oinyCRyx8T2CQ==}

  nanoid@3.3.11:
    resolution: {integrity: sha512-N8SpfPUnUp1bK+PMYW8qSWdl9U+wwNWI4QKxOYDy9JAro3WMX7p2OeVRF9v+347pnakNevPmiHhNmZ2HbFA76w==}
    engines: {node: ^10 || ^12 || ^13.7 || ^14 || >=15.0.1}
    hasBin: true

  next-auth@4.24.13:
    resolution: {integrity: sha512-sgObCfcfL7BzIK76SS5TnQtc3yo2Oifp/yIpfv6fMfeBOiBJkDWF3A2y9+yqnmJ4JKc2C+nMjSjmgDeTwgN1rQ==}
    peerDependencies:
      '@auth/core': 0.34.3
      next: ^12.2.5 || ^13 || ^14 || ^15 || ^16
      nodemailer: ^7.0.7
      react: ^17.0.2 || ^18 || ^19
      react-dom: ^17.0.2 || ^18 || ^19
    peerDependenciesMeta:
      '@auth/core':
        optional: true
      nodemailer:
        optional: true

  next@16.0.7:
    resolution: {integrity: sha512-3mBRJyPxT4LOxAJI6IsXeFtKfiJUbjCLgvXO02fV8Wy/lIhPvP94Fe7dGhUgHXcQy4sSuYwQNcOLhIfOm0rL0A==}
    engines: {node: '>=20.9.0'}
    hasBin: true
    peerDependencies:
      '@opentelemetry/api': ^1.1.0
      '@playwright/test': ^1.51.1
      babel-plugin-react-compiler: '*'
      react: ^18.2.0 || 19.0.0-rc-de68d2f4-20241204 || ^19.0.0
      react-dom: ^18.2.0 || 19.0.0-rc-de68d2f4-20241204 || ^19.0.0
      sass: ^1.3.0
    peerDependenciesMeta:
      '@opentelemetry/api':
        optional: true
      '@playwright/test':
        optional: true
      babel-plugin-react-compiler:
        optional: true
      sass:
        optional: true

  oauth@0.9.15:
    resolution: {integrity: sha512-a5ERWK1kh38ExDEfoO6qUHJb32rd7aYmPHuyCu3Fta/cnICvYmgd2uhuKXvPD+PXB+gCEYYEaQdIRAjCOwAKNA==}

  object-hash@2.2.0:
    resolution: {integrity: sha512-gScRMn0bS5fH+IuwyIFgnh9zBdo4DV+6GhygmWM9HyNJSgS0hScp1f5vjtm7oIIOiT9trXrShAkLFSc2IqKNgw==}
    engines: {node: '>= 6'}

  oidc-token-hash@5.2.0:
    resolution: {integrity: sha512-6gj2m8cJZ+iSW8bm0FXdGF0YhIQbKrfP4yWTNzxc31U6MOjfEmB1rHvlYvxI1B7t7BCi1F2vYTT6YhtQRG4hxw==}
    engines: {node: ^10.13.0 || >=12.0.0}

  openid-client@5.7.1:
    resolution: {integrity: sha512-jDBPgSVfTnkIh71Hg9pRvtJc6wTwqjRkN88+gCFtYWrlP4Yx2Dsrow8uPi3qLr/aeymPF3o2+dS+wOpglK04ew==}

  picocolors@1.1.1:
    resolution: {integrity: sha512-xceH2snhtb5M9liqDsmEw56le376mTZkEX/jEb/RxNFyegNul7eNslCXP9FDj/Lcu0X8KEyMceP2ntpaHrDEVA==}

  postcss@8.4.31:
    resolution: {integrity: sha512-PS08Iboia9mts/2ygV3eLpY5ghnUcfLV/EXTOW1E2qYxJKGGBUtNjN76FYHnMs36RmARn41bC0AZmn+rR0OVpQ==}
    engines: {node: ^10 || ^12 || >=14}

  postcss@8.5.6:
    resolution: {integrity: sha512-3Ybi1tAuwAP9s0r1UQ2J4n5Y0G05bJkpUIO0/bI9MhwmD70S5aTWbXGBwxHrelT+XM1k6dM0pk+SwNkpTRN7Pg==}
    engines: {node: ^10 || ^12 || >=14}

  preact-render-to-string@5.2.6:
    resolution: {integrity: sha512-JyhErpYOvBV1hEPwIxc/fHWXPfnEGdRKxc8gFdAZ7XV4tlzyzG847XAyEZqoDnynP88akM4eaHcSOzNcLWFguw==}
    peerDependencies:
      preact: '>=10'

  preact@10.28.0:
    resolution: {integrity: sha512-rytDAoiXr3+t6OIP3WGlDd0ouCUG1iCWzkcY3++Nreuoi17y6T5i/zRhe6uYfoVcxq6YU+sBtJouuRDsq8vvqA==}

  pretty-format@3.8.0:
    resolution: {integrity: sha512-WuxUnVtlWL1OfZFQFuqvnvs6MiAGk9UNsBostyBOB0Is9wb5uRESevA6rnl/rkksXaGX3GzZhPup5d6Vp1nFew==}

  react-dom@19.2.0:
    resolution: {integrity: sha512-UlbRu4cAiGaIewkPyiRGJk0imDN2T3JjieT6spoL2UeSf5od4n5LB/mQ4ejmxhCFT1tYe8IvaFulzynWovsEFQ==}
    peerDependencies:
      react: ^19.2.0

  react-hook-form@7.68.0:
    resolution: {integrity: sha512-oNN3fjrZ/Xo40SWlHf1yCjlMK417JxoSJVUXQjGdvdRCU07NTFei1i1f8ApUAts+IVh14e4EdakeLEA+BEAs/Q==}
    engines: {node: '>=18.0.0'}
    peerDependencies:
      react: ^16.8.0 || ^17 || ^18 || ^19

  react-is@19.2.1:
    resolution: {integrity: sha512-L7BnWgRbMwzMAubQcS7sXdPdNLmKlucPlopgAzx7FtYbksWZgEWiuYM5x9T6UqS2Ne0rsgQTq5kY2SGqpzUkYA==}

  react-redux@9.2.0:
    resolution: {integrity: sha512-ROY9fvHhwOD9ySfrF0wmvu//bKCQ6AeZZq1nJNtbDC+kk5DuSuNX/n6YWYF/SYy7bSba4D4FSz8DJeKY/S/r+g==}
    peerDependencies:
      '@types/react': ^18.2.25 || ^19
      react: ^18.0 || ^19
      redux: ^5.0.0
    peerDependenciesMeta:
      '@types/react':
        optional: true
      redux:
        optional: true

  react@19.2.0:
    resolution: {integrity: sha512-tmbWg6W31tQLeB5cdIBOicJDJRR2KzXsV7uSK9iNfLWQ5bIZfxuPEHp7M8wiHyHnn0DD1i7w3Zmin0FtkrwoCQ==}
    engines: {node: '>=0.10.0'}

  recharts@3.5.1:
    resolution: {integrity: sha512-+v+HJojK7gnEgG6h+b2u7k8HH7FhyFUzAc4+cPrsjL4Otdgqr/ecXzAnHciqlzV1ko064eNcsdzrYOM78kankA==}
    engines: {node: '>=18'}
    peerDependencies:
      react: ^16.8.0 || ^17.0.0 || ^18.0.0 || ^19.0.0
      react-dom: ^16.0.0 || ^17.0.0 || ^18.0.0 || ^19.0.0
      react-is: ^16.8.0 || ^17.0.0 || ^18.0.0 || ^19.0.0

  redux-thunk@3.1.0:
    resolution: {integrity: sha512-NW2r5T6ksUKXCabzhL9z+h206HQw/NJkcLm1GPImRQ8IzfXwRGqjVhKJGauHirT0DAuyy6hjdnMZaRoAcy0Klw==}
    peerDependencies:
      redux: ^5.0.0

  redux@5.0.1:
    resolution: {integrity: sha512-M9/ELqF6fy8FwmkpnF0S3YKOqMyoWJ4+CS5Efg2ct3oY9daQvd/Pc71FpGZsVsbl3Cpb+IIcjBDUnnyBdQbq4w==}

  reselect@5.1.1:
    resolution: {integrity: sha512-K/BG6eIky/SBpzfHZv/dd+9JBFiS4SWV7FIujVyJRux6e45+73RaUHXLmIR1f7WOMaQ0U1km6qwklRQxpJJY0w==}

  scheduler@0.27.0:
    resolution: {integrity: sha512-eNv+WrVbKu1f3vbYJT/xtiF5syA5HPIMtf9IgY/nKg0sWqzAUEvqY/xm7OcZc/qafLx/iO9FgOmeSAp4v5ti/Q==}

  semver@7.7.3:
    resolution: {integrity: sha512-SdsKMrI9TdgjdweUSR9MweHA4EJ8YxHn8DFaDisvhVlUOe4BF1tLD7GAj0lIqWVl+dPb/rExr0Btby5loQm20Q==}
    engines: {node: '>=10'}
    hasBin: true

  sharp@0.34.5:
    resolution: {integrity: sha512-Ou9I5Ft9WNcCbXrU9cMgPBcCK8LiwLqcbywW3t4oDV37n1pzpuNLsYiAV8eODnjbtQlSDwZ2cUEeQz4E54Hltg==}
    engines: {node: ^18.17.0 || ^20.3.0 || >=21.0.0}

  source-map-js@1.2.1:
    resolution: {integrity: sha512-UXWMKhLOwVKb728IUtQPXxfYU+usdybtUrK/8uGE8CQMvrhOpwvzDBwj0QhSL7MQc7vIsISBG8VQ8+IDQxpfQA==}
    engines: {node: '>=0.10.0'}

  styled-jsx@5.1.6:
    resolution: {integrity: sha512-qSVyDTeMotdvQYoHWLNGwRFJHC+i+ZvdBRYosOFgC+Wg1vx4frN2/RG/NA7SYqqvKNLf39P2LSRA2pu6n0XYZA==}
    engines: {node: '>= 12.0.0'}
    peerDependencies:
      '@babel/core': '*'
      babel-plugin-macros: '*'
      react: '>= 16.8.0 || 17.x.x || ^18.0.0-0 || ^19.0.0-0'
    peerDependenciesMeta:
      '@babel/core':
        optional: true
      babel-plugin-macros:
        optional: true

  tailwind-merge@3.4.0:
    resolution: {integrity: sha512-uSaO4gnW+b3Y2aWoWfFpX62vn2sR3skfhbjsEnaBI81WD1wBLlHZe5sWf0AqjksNdYTbGBEd0UasQMT3SNV15g==}

  tailwindcss@4.1.17:
    resolution: {integrity: sha512-j9Ee2YjuQqYT9bbRTfTZht9W/ytp5H+jJpZKiYdP/bpnXARAuELt9ofP0lPnmHjbga7SNQIxdTAXCmtKVYjN+Q==}

  tapable@2.3.0:
    resolution: {integrity: sha512-g9ljZiwki/LfxmQADO3dEY1CbpmXT5Hm2fJ+QaGKwSXUylMybePR7/67YW7jOrrvjEgL1Fmz5kzyAjWVWLlucg==}
    engines: {node: '>=6'}

  tiny-invariant@1.3.3:
    resolution: {integrity: sha512-+FbBPE1o9QAYvviau/qC5SE3caw21q3xkvWKBtja5vgqOWIHHJ3ioaq1VPfn/Szqctz2bU/oYeKd9/z5BL+PVg==}

  tslib@2.8.1:
    resolution: {integrity: sha512-oJFu94HQb+KVduSUQL7wnpmqnfmLsOA/nAh6b6EH0wCEoK0/mPeXU6c3wKDV83MkOuHPRHtSXKKU99IBazS/2w==}

  tw-animate-css@1.4.0:
    resolution: {integrity: sha512-7bziOlRqH0hJx80h/3mbicLW7o8qLsH5+RaLR2t+OHM3D0JlWGODQKQ4cxbK7WlvmUxpcj6Kgu6EKqjrGFe3QQ==}

  typescript@5.9.3:
    resolution: {integrity: sha512-jl1vZzPDinLr9eUt3J/t7V6FgNEw9QjvBPdysz9KfQDD41fQrC2Y4vKQdiaUpFT4bXlb1RHhLpp8wtm6M5TgSw==}
    engines: {node: '>=14.17'}
    hasBin: true

  undici-types@6.21.0:
    resolution: {integrity: sha512-iwDZqg0QAGrg9Rav5H4n0M64c3mkR59cJ6wQp+7C4nI0gsmExaedaYLNO44eT4AtBBwjbTiGPMlt2Md0T9H9JQ==}

  use-sync-external-store@1.6.0:
    resolution: {integrity: sha512-Pp6GSwGP/NrPIrxVFAIkOQeyw8lFenOHijQWkUTrDvrF4ALqylP2C/KCkeS9dpUM3KvYRQhna5vt7IL95+ZQ9w==}
    peerDependencies:
      react: ^16.8.0 || ^17.0.0 || ^18.0.0 || ^19.0.0

  uuid@8.3.2:
    resolution: {integrity: sha512-+NYs2QeMWy+GWFOEm9xnn6HCDp0l7QBD7ml8zLUmJ+93Q5NF0NocErnwkTkXVFNiX3/fpC6afS8Dhb/gz7R7eg==}
    hasBin: true

  victory-vendor@37.3.6:
    resolution: {integrity: sha512-SbPDPdDBYp+5MJHhBCAyI7wKM3d5ivekigc2Dk2s7pgbZ9wIgIBYGVw4zGHBml/qTFbexrofXW6Gu4noGxrOwQ==}

  yallist@4.0.0:
    resolution: {integrity: sha512-3wdGidZyq5PB084XLES5TpOSRA3wjXAlIWMhum2kRcv/41Sn2emQ0dycQW4uZXLejwKvg6EsvbdlVL+FYEct7A==}

  zod@4.1.13:
    resolution: {integrity: sha512-AvvthqfqrAhNH9dnfmrfKzX5upOdjUVJYFqNSlkmGf64gRaTzlPwz99IHYnVs28qYAybvAlBV+H7pn0saFY4Ig==}

snapshots:

  '@alloc/quick-lru@5.2.0': {}

  '@babel/runtime@7.28.4': {}

  '@emnapi/runtime@1.7.1':
    dependencies:
      tslib: 2.8.1
    optional: true

  '@hookform/resolvers@5.2.2(react-hook-form@7.68.0(react@19.2.0))':
    dependencies:
      '@standard-schema/utils': 0.3.0
      react-hook-form: 7.68.0(react@19.2.0)

  '@img/colour@1.0.0':
    optional: true

  '@img/sharp-darwin-arm64@0.34.5':
    optionalDependencies:
      '@img/sharp-libvips-darwin-arm64': 1.2.4
    optional: true

  '@img/sharp-darwin-x64@0.34.5':
    optionalDependencies:
      '@img/sharp-libvips-darwin-x64': 1.2.4
    optional: true

  '@img/sharp-libvips-darwin-arm64@1.2.4':
    optional: true

  '@img/sharp-libvips-darwin-x64@1.2.4':
    optional: true

  '@img/sharp-libvips-linux-arm64@1.2.4':
    optional: true

  '@img/sharp-libvips-linux-arm@1.2.4':
    optional: true

  '@img/sharp-libvips-linux-ppc64@1.2.4':
    optional: true

  '@img/sharp-libvips-linux-riscv64@1.2.4':
    optional: true

  '@img/sharp-libvips-linux-s390x@1.2.4':
    optional: true

  '@img/sharp-libvips-linux-x64@1.2.4':
    optional: true

  '@img/sharp-libvips-linuxmusl-arm64@1.2.4':
    optional: true

  '@img/sharp-libvips-linuxmusl-x64@1.2.4':
    optional: true

  '@img/sharp-linux-arm64@0.34.5':
    optionalDependencies:
      '@img/sharp-libvips-linux-arm64': 1.2.4
    optional: true

  '@img/sharp-linux-arm@0.34.5':
    optionalDependencies:
      '@img/sharp-libvips-linux-arm': 1.2.4
    optional: true

  '@img/sharp-linux-ppc64@0.34.5':
    optionalDependencies:
      '@img/sharp-libvips-linux-ppc64': 1.2.4
    optional: true

  '@img/sharp-linux-riscv64@0.34.5':
    optionalDependencies:
      '@img/sharp-libvips-linux-riscv64': 1.2.4
    optional: true

  '@img/sharp-linux-s390x@0.34.5':
    optionalDependencies:
      '@img/sharp-libvips-linux-s390x': 1.2.4
    optional: true

  '@img/sharp-linux-x64@0.34.5':
    optionalDependencies:
      '@img/sharp-libvips-linux-x64': 1.2.4
    optional: true

  '@img/sharp-linuxmusl-arm64@0.34.5':
    optionalDependencies:
      '@img/sharp-libvips-linuxmusl-arm64': 1.2.4
    optional: true

  '@img/sharp-linuxmusl-x64@0.34.5':
    optionalDependencies:
      '@img/sharp-libvips-linuxmusl-x64': 1.2.4
    optional: true

  '@img/sharp-wasm32@0.34.5':
    dependencies:
      '@emnapi/runtime': 1.7.1
    optional: true

  '@img/sharp-win32-arm64@0.34.5':
    optional: true

  '@img/sharp-win32-ia32@0.34.5':
    optional: true

  '@img/sharp-win32-x64@0.34.5':
    optional: true

  '@jridgewell/gen-mapping@0.3.13':
    dependencies:
      '@jridgewell/sourcemap-codec': 1.5.5
      '@jridgewell/trace-mapping': 0.3.31

  '@jridgewell/remapping@2.3.5':
    dependencies:
      '@jridgewell/gen-mapping': 0.3.13
      '@jridgewell/trace-mapping': 0.3.31

  '@jridgewell/resolve-uri@3.1.2': {}

  '@jridgewell/sourcemap-codec@1.5.5': {}

  '@jridgewell/trace-mapping@0.3.31':
    dependencies:
      '@jridgewell/resolve-uri': 3.1.2
      '@jridgewell/sourcemap-codec': 1.5.5

  '@next/env@16.0.7': {}

  '@next/swc-darwin-arm64@16.0.7':
    optional: true

  '@next/swc-darwin-x64@16.0.7':
    optional: true

  '@next/swc-linux-arm64-gnu@16.0.7':
    optional: true

  '@next/swc-linux-arm64-musl@16.0.7':
    optional: true

  '@next/swc-linux-x64-gnu@16.0.7':
    optional: true

  '@next/swc-linux-x64-musl@16.0.7':
    optional: true

  '@next/swc-win32-arm64-msvc@16.0.7':
    optional: true

  '@next/swc-win32-x64-msvc@16.0.7':
    optional: true

  '@panva/hkdf@1.2.1': {}

  '@radix-ui/react-compose-refs@1.1.2(@types/react@19.2.7)(react@19.2.0)':
    dependencies:
      react: 19.2.0
    optionalDependencies:
      '@types/react': 19.2.7

  '@radix-ui/react-label@2.1.8(@types/react-dom@19.2.3(@types/react@19.2.7))(@types/react@19.2.7)(react-dom@19.2.0(react@19.2.0))(react@19.2.0)':
    dependencies:
      '@radix-ui/react-primitive': 2.1.4(@types/react-dom@19.2.3(@types/react@19.2.7))(@types/react@19.2.7)(react-dom@19.2.0(react@19.2.0))(react@19.2.0)
      react: 19.2.0
      react-dom: 19.2.0(react@19.2.0)
    optionalDependencies:
      '@types/react': 19.2.7
      '@types/react-dom': 19.2.3(@types/react@19.2.7)

  '@radix-ui/react-primitive@2.1.4(@types/react-dom@19.2.3(@types/react@19.2.7))(@types/react@19.2.7)(react-dom@19.2.0(react@19.2.0))(react@19.2.0)':
    dependencies:
      '@radix-ui/react-slot': 1.2.4(@types/react@19.2.7)(react@19.2.0)
      react: 19.2.0
      react-dom: 19.2.0(react@19.2.0)
    optionalDependencies:
      '@types/react': 19.2.7
      '@types/react-dom': 19.2.3(@types/react@19.2.7)

  '@radix-ui/react-slot@1.2.4(@types/react@19.2.7)(react@19.2.0)':
    dependencies:
      '@radix-ui/react-compose-refs': 1.1.2(@types/react@19.2.7)(react@19.2.0)
      react: 19.2.0
    optionalDependencies:
      '@types/react': 19.2.7

  '@reduxjs/toolkit@2.11.0(react-redux@9.2.0(@types/react@19.2.7)(react@19.2.0)(redux@5.0.1))(react@19.2.0)':
    dependencies:
      '@standard-schema/spec': 1.0.0
      '@standard-schema/utils': 0.3.0
      immer: 11.0.1
      redux: 5.0.1
      redux-thunk: 3.1.0(redux@5.0.1)
      reselect: 5.1.1
    optionalDependencies:
      react: 19.2.0
      react-redux: 9.2.0(@types/react@19.2.7)(react@19.2.0)(redux@5.0.1)

  '@standard-schema/spec@1.0.0': {}

  '@standard-schema/utils@0.3.0': {}

  '@swc/helpers@0.5.15':
    dependencies:
      tslib: 2.8.1

  '@tailwindcss/node@4.1.17':
    dependencies:
      '@jridgewell/remapping': 2.3.5
      enhanced-resolve: 5.18.3
      jiti: 2.6.1
      lightningcss: 1.30.2
      magic-string: 0.30.21
      source-map-js: 1.2.1
      tailwindcss: 4.1.17

  '@tailwindcss/oxide-android-arm64@4.1.17':
    optional: true

  '@tailwindcss/oxide-darwin-arm64@4.1.17':
    optional: true

  '@tailwindcss/oxide-darwin-x64@4.1.17':
    optional: true

  '@tailwindcss/oxide-freebsd-x64@4.1.17':
    optional: true

  '@tailwindcss/oxide-linux-arm-gnueabihf@4.1.17':
    optional: true

  '@tailwindcss/oxide-linux-arm64-gnu@4.1.17':
    optional: true

  '@tailwindcss/oxide-linux-arm64-musl@4.1.17':
    optional: true

  '@tailwindcss/oxide-linux-x64-gnu@4.1.17':
    optional: true

  '@tailwindcss/oxide-linux-x64-musl@4.1.17':
    optional: true

  '@tailwindcss/oxide-wasm32-wasi@4.1.17':
    optional: true

  '@tailwindcss/oxide-win32-arm64-msvc@4.1.17':
    optional: true

  '@tailwindcss/oxide-win32-x64-msvc@4.1.17':
    optional: true

  '@tailwindcss/oxide@4.1.17':
    optionalDependencies:
      '@tailwindcss/oxide-android-arm64': 4.1.17
      '@tailwindcss/oxide-darwin-arm64': 4.1.17
      '@tailwindcss/oxide-darwin-x64': 4.1.17
      '@tailwindcss/oxide-freebsd-x64': 4.1.17
      '@tailwindcss/oxide-linux-arm-gnueabihf': 4.1.17
      '@tailwindcss/oxide-linux-arm64-gnu': 4.1.17
      '@tailwindcss/oxide-linux-arm64-musl': 4.1.17
      '@tailwindcss/oxide-linux-x64-gnu': 4.1.17
      '@tailwindcss/oxide-linux-x64-musl': 4.1.17
      '@tailwindcss/oxide-wasm32-wasi': 4.1.17
      '@tailwindcss/oxide-win32-arm64-msvc': 4.1.17
      '@tailwindcss/oxide-win32-x64-msvc': 4.1.17

  '@tailwindcss/postcss@4.1.17':
    dependencies:
      '@alloc/quick-lru': 5.2.0
      '@tailwindcss/node': 4.1.17
      '@tailwindcss/oxide': 4.1.17
      postcss: 8.5.6
      tailwindcss: 4.1.17

  '@tanstack/query-core@5.90.12': {}

  '@tanstack/react-query@5.90.12(react@19.2.0)':
    dependencies:
      '@tanstack/query-core': 5.90.12
      react: 19.2.0

  '@tanstack/react-table@8.21.3(react-dom@19.2.0(react@19.2.0))(react@19.2.0)':
    dependencies:
      '@tanstack/table-core': 8.21.3
      react: 19.2.0
      react-dom: 19.2.0(react@19.2.0)

  '@tanstack/table-core@8.21.3': {}

  '@types/d3-array@3.2.2': {}

  '@types/d3-color@3.1.3': {}

  '@types/d3-ease@3.0.2': {}

  '@types/d3-interpolate@3.0.4':
    dependencies:
      '@types/d3-color': 3.1.3

  '@types/d3-path@3.1.1': {}

  '@types/d3-scale@4.0.9':
    dependencies:
      '@types/d3-time': 3.0.4

  '@types/d3-shape@3.1.7':
    dependencies:
      '@types/d3-path': 3.1.1

  '@types/d3-time@3.0.4': {}

  '@types/d3-timer@3.0.2': {}

  '@types/node@20.19.25':
    dependencies:
      undici-types: 6.21.0

  '@types/react-dom@19.2.3(@types/react@19.2.7)':
    dependencies:
      '@types/react': 19.2.7

  '@types/react@19.2.7':
    dependencies:
      csstype: 3.2.3

  '@types/use-sync-external-store@0.0.6': {}

  caniuse-lite@1.0.30001759: {}

  class-variance-authority@0.7.1:
    dependencies:
      clsx: 2.1.1

  client-only@0.0.1: {}

  clsx@2.1.1: {}

  cookie@0.7.2: {}

  csstype@3.2.3: {}

  d3-array@3.2.4:
    dependencies:
      internmap: 2.0.3

  d3-color@3.1.0: {}

  d3-ease@3.0.1: {}

  d3-format@3.1.0: {}

  d3-interpolate@3.0.1:
    dependencies:
      d3-color: 3.1.0

  d3-path@3.1.0: {}

  d3-scale@4.0.2:
    dependencies:
      d3-array: 3.2.4
      d3-format: 3.1.0
      d3-interpolate: 3.0.1
      d3-time: 3.1.0
      d3-time-format: 4.1.0

  d3-shape@3.2.0:
    dependencies:
      d3-path: 3.1.0

  d3-time-format@4.1.0:
    dependencies:
      d3-time: 3.1.0

  d3-time@3.1.0:
    dependencies:
      d3-array: 3.2.4

  d3-timer@3.0.1: {}

  date-fns@4.1.0: {}

  decimal.js-light@2.5.1: {}

  detect-libc@2.1.2: {}

  enhanced-resolve@5.18.3:
    dependencies:
      graceful-fs: 4.2.11
      tapable: 2.3.0

  es-toolkit@1.42.0: {}

  eventemitter3@5.0.1: {}

  graceful-fs@4.2.11: {}

  immer@10.2.0: {}

  immer@11.0.1: {}

  internmap@2.0.3: {}

  jiti@2.6.1: {}

  jose@4.15.9: {}

  lightningcss-android-arm64@1.30.2:
    optional: true

  lightningcss-darwin-arm64@1.30.2:
    optional: true

  lightningcss-darwin-x64@1.30.2:
    optional: true

  lightningcss-freebsd-x64@1.30.2:
    optional: true

  lightningcss-linux-arm-gnueabihf@1.30.2:
    optional: true

  lightningcss-linux-arm64-gnu@1.30.2:
    optional: true

  lightningcss-linux-arm64-musl@1.30.2:
    optional: true

  lightningcss-linux-x64-gnu@1.30.2:
    optional: true

  lightningcss-linux-x64-musl@1.30.2:
    optional: true

  lightningcss-win32-arm64-msvc@1.30.2:
    optional: true

  lightningcss-win32-x64-msvc@1.30.2:
    optional: true

  lightningcss@1.30.2:
    dependencies:
      detect-libc: 2.1.2
    optionalDependencies:
      lightningcss-android-arm64: 1.30.2
      lightningcss-darwin-arm64: 1.30.2
      lightningcss-darwin-x64: 1.30.2
      lightningcss-freebsd-x64: 1.30.2
      lightningcss-linux-arm-gnueabihf: 1.30.2
      lightningcss-linux-arm64-gnu: 1.30.2
      lightningcss-linux-arm64-musl: 1.30.2
      lightningcss-linux-x64-gnu: 1.30.2
      lightningcss-linux-x64-musl: 1.30.2
      lightningcss-win32-arm64-msvc: 1.30.2
      lightningcss-win32-x64-msvc: 1.30.2

  lru-cache@6.0.0:
    dependencies:
      yallist: 4.0.0

  lucide-react@0.556.0(react@19.2.0):
    dependencies:
      react: 19.2.0

  magic-string@0.30.21:
    dependencies:
      '@jridgewell/sourcemap-codec': 1.5.5

  nanoid@3.3.11: {}

  next-auth@4.24.13(next@16.0.7(react-dom@19.2.0(react@19.2.0))(react@19.2.0))(react-dom@19.2.0(react@19.2.0))(react@19.2.0):
    dependencies:
      '@babel/runtime': 7.28.4
      '@panva/hkdf': 1.2.1
      cookie: 0.7.2
      jose: 4.15.9
      next: 16.0.7(react-dom@19.2.0(react@19.2.0))(react@19.2.0)
      oauth: 0.9.15
      openid-client: 5.7.1
      preact: 10.28.0
      preact-render-to-string: 5.2.6(preact@10.28.0)
      react: 19.2.0
      react-dom: 19.2.0(react@19.2.0)
      uuid: 8.3.2

  next@16.0.7(react-dom@19.2.0(react@19.2.0))(react@19.2.0):
    dependencies:
      '@next/env': 16.0.7
      '@swc/helpers': 0.5.15
      caniuse-lite: 1.0.30001759
      postcss: 8.4.31
      react: 19.2.0
      react-dom: 19.2.0(react@19.2.0)
      styled-jsx: 5.1.6(react@19.2.0)
    optionalDependencies:
      '@next/swc-darwin-arm64': 16.0.7
      '@next/swc-darwin-x64': 16.0.7
      '@next/swc-linux-arm64-gnu': 16.0.7
      '@next/swc-linux-arm64-musl': 16.0.7
      '@next/swc-linux-x64-gnu': 16.0.7
      '@next/swc-linux-x64-musl': 16.0.7
      '@next/swc-win32-arm64-msvc': 16.0.7
      '@next/swc-win32-x64-msvc': 16.0.7
      sharp: 0.34.5
    transitivePeerDependencies:
      - '@babel/core'
      - babel-plugin-macros

  oauth@0.9.15: {}

  object-hash@2.2.0: {}

  oidc-token-hash@5.2.0: {}

  openid-client@5.7.1:
    dependencies:
      jose: 4.15.9
      lru-cache: 6.0.0
      object-hash: 2.2.0
      oidc-token-hash: 5.2.0

  picocolors@1.1.1: {}

  postcss@8.4.31:
    dependencies:
      nanoid: 3.3.11
      picocolors: 1.1.1
      source-map-js: 1.2.1

  postcss@8.5.6:
    dependencies:
      nanoid: 3.3.11
      picocolors: 1.1.1
      source-map-js: 1.2.1

  preact-render-to-string@5.2.6(preact@10.28.0):
    dependencies:
      preact: 10.28.0
      pretty-format: 3.8.0

  preact@10.28.0: {}

  pretty-format@3.8.0: {}

  react-dom@19.2.0(react@19.2.0):
    dependencies:
      react: 19.2.0
      scheduler: 0.27.0

  react-hook-form@7.68.0(react@19.2.0):
    dependencies:
      react: 19.2.0

  react-is@19.2.1: {}

  react-redux@9.2.0(@types/react@19.2.7)(react@19.2.0)(redux@5.0.1):
    dependencies:
      '@types/use-sync-external-store': 0.0.6
      react: 19.2.0
      use-sync-external-store: 1.6.0(react@19.2.0)
    optionalDependencies:
      '@types/react': 19.2.7
      redux: 5.0.1

  react@19.2.0: {}

  recharts@3.5.1(@types/react@19.2.7)(react-dom@19.2.0(react@19.2.0))(react-is@19.2.1)(react@19.2.0)(redux@5.0.1):
    dependencies:
      '@reduxjs/toolkit': 2.11.0(react-redux@9.2.0(@types/react@19.2.7)(react@19.2.0)(redux@5.0.1))(react@19.2.0)
      clsx: 2.1.1
      decimal.js-light: 2.5.1
      es-toolkit: 1.42.0
      eventemitter3: 5.0.1
      immer: 10.2.0
      react: 19.2.0
      react-dom: 19.2.0(react@19.2.0)
      react-is: 19.2.1
      react-redux: 9.2.0(@types/react@19.2.7)(react@19.2.0)(redux@5.0.1)
      reselect: 5.1.1
      tiny-invariant: 1.3.3
      use-sync-external-store: 1.6.0(react@19.2.0)
      victory-vendor: 37.3.6
    transitivePeerDependencies:
      - '@types/react'
      - redux

  redux-thunk@3.1.0(redux@5.0.1):
    dependencies:
      redux: 5.0.1

  redux@5.0.1: {}

  reselect@5.1.1: {}

  scheduler@0.27.0: {}

  semver@7.7.3:
    optional: true

  sharp@0.34.5:
    dependencies:
      '@img/colour': 1.0.0
      detect-libc: 2.1.2
      semver: 7.7.3
    optionalDependencies:
      '@img/sharp-darwin-arm64': 0.34.5
      '@img/sharp-darwin-x64': 0.34.5
      '@img/sharp-libvips-darwin-arm64': 1.2.4
      '@img/sharp-libvips-darwin-x64': 1.2.4
      '@img/sharp-libvips-linux-arm': 1.2.4
      '@img/sharp-libvips-linux-arm64': 1.2.4
      '@img/sharp-libvips-linux-ppc64': 1.2.4
      '@img/sharp-libvips-linux-riscv64': 1.2.4
      '@img/sharp-libvips-linux-s390x': 1.2.4
      '@img/sharp-libvips-linux-x64': 1.2.4
      '@img/sharp-libvips-linuxmusl-arm64': 1.2.4
      '@img/sharp-libvips-linuxmusl-x64': 1.2.4
      '@img/sharp-linux-arm': 0.34.5
      '@img/sharp-linux-arm64': 0.34.5
      '@img/sharp-linux-ppc64': 0.34.5
      '@img/sharp-linux-riscv64': 0.34.5
      '@img/sharp-linux-s390x': 0.34.5
      '@img/sharp-linux-x64': 0.34.5
      '@img/sharp-linuxmusl-arm64': 0.34.5
      '@img/sharp-linuxmusl-x64': 0.34.5
      '@img/sharp-wasm32': 0.34.5
      '@img/sharp-win32-arm64': 0.34.5
      '@img/sharp-win32-ia32': 0.34.5
      '@img/sharp-win32-x64': 0.34.5
    optional: true

  source-map-js@1.2.1: {}

  styled-jsx@5.1.6(react@19.2.0):
    dependencies:
      client-only: 0.0.1
      react: 19.2.0

  tailwind-merge@3.4.0: {}

  tailwindcss@4.1.17: {}

  tapable@2.3.0: {}

  tiny-invariant@1.3.3: {}

  tslib@2.8.1: {}

  tw-animate-css@1.4.0: {}

  typescript@5.9.3: {}

  undici-types@6.21.0: {}

  use-sync-external-store@1.6.0(react@19.2.0):
    dependencies:
      react: 19.2.0

  uuid@8.3.2: {}

  victory-vendor@37.3.6:
    dependencies:
      '@types/d3-array': 3.2.2
      '@types/d3-ease': 3.0.2
      '@types/d3-interpolate': 3.0.4
      '@types/d3-scale': 4.0.9
      '@types/d3-shape': 3.1.7
      '@types/d3-time': 3.0.4
      '@types/d3-timer': 3.0.2
      d3-array: 3.2.4
      d3-ease: 3.0.1
      d3-interpolate: 3.0.1
      d3-scale: 4.0.2
      d3-shape: 3.2.0
      d3-time: 3.1.0
      d3-timer: 3.0.1

  yallist@4.0.0: {}

  zod@4.1.13: {}

### .\frontend\pnpm-lock.yaml END ###

### .\frontend\postcss.config.mjs BEGIN ###
const config = {
  plugins: {
    "@tailwindcss/postcss": {},
  },
};

export default config;

### .\frontend\postcss.config.mjs END ###

### .\frontend\public\file.svg BEGIN ###
<svg fill="none" viewBox="0 0 16 16" xmlns="http://www.w3.org/2000/svg"><path d="M14.5 13.5V5.41a1 1 0 0 0-.3-.7L9.8.29A1 1 0 0 0 9.08 0H1.5v13.5A2.5 2.5 0 0 0 4 16h8a2.5 2.5 0 0 0 2.5-2.5m-1.5 0v-7H8v-5H3v12a1 1 0 0 0 1 1h8a1 1 0 0 0 1-1M9.5 5V2.12L12.38 5zM5.13 5h-.62v1.25h2.12V5zm-.62 3h7.12v1.25H4.5zm.62 3h-.62v1.25h7.12V11z" clip-rule="evenodd" fill="#666" fill-rule="evenodd"/></svg>
### .\frontend\public\file.svg END ###

### .\frontend\public\globe.svg BEGIN ###
<svg fill="none" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16"><g clip-path="url(#a)"><path fill-rule="evenodd" clip-rule="evenodd" d="M10.27 14.1a6.5 6.5 0 0 0 3.67-3.45q-1.24.21-2.7.34-.31 1.83-.97 3.1M8 16A8 8 0 1 0 8 0a8 8 0 0 0 0 16m.48-1.52a7 7 0 0 1-.96 0H7.5a4 4 0 0 1-.84-1.32q-.38-.89-.63-2.08a40 40 0 0 0 3.92 0q-.25 1.2-.63 2.08a4 4 0 0 1-.84 1.31zm2.94-4.76q1.66-.15 2.95-.43a7 7 0 0 0 0-2.58q-1.3-.27-2.95-.43a18 18 0 0 1 0 3.44m-1.27-3.54a17 17 0 0 1 0 3.64 39 39 0 0 1-4.3 0 17 17 0 0 1 0-3.64 39 39 0 0 1 4.3 0m1.1-1.17q1.45.13 2.69.34a6.5 6.5 0 0 0-3.67-3.44q.65 1.26.98 3.1M8.48 1.5l.01.02q.41.37.84 1.31.38.89.63 2.08a40 40 0 0 0-3.92 0q.25-1.2.63-2.08a4 4 0 0 1 .85-1.32 7 7 0 0 1 .96 0m-2.75.4a6.5 6.5 0 0 0-3.67 3.44 29 29 0 0 1 2.7-.34q.31-1.83.97-3.1M4.58 6.28q-1.66.16-2.95.43a7 7 0 0 0 0 2.58q1.3.27 2.95.43a18 18 0 0 1 0-3.44m.17 4.71q-1.45-.12-2.69-.34a6.5 6.5 0 0 0 3.67 3.44q-.65-1.27-.98-3.1" fill="#666"/></g><defs><clipPath id="a"><path fill="#fff" d="M0 0h16v16H0z"/></clipPath></defs></svg>
### .\frontend\public\globe.svg END ###

### .\frontend\public\next.svg BEGIN ###
<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 394 80"><path fill="#000" d="M262 0h68.5v12.7h-27.2v66.6h-13.6V12.7H262V0ZM149 0v12.7H94v20.4h44.3v12.6H94v21h55v12.6H80.5V0h68.7zm34.3 0h-17.8l63.8 79.4h17.9l-32-39.7 32-39.6h-17.9l-23 28.6-23-28.6zm18.3 56.7-9-11-27.1 33.7h17.8l18.3-22.7z"/><path fill="#000" d="M81 79.3 17 0H0v79.3h13.6V17l50.2 62.3H81Zm252.6-.4c-1 0-1.8-.4-2.5-1s-1.1-1.6-1.1-2.6.3-1.8 1-2.5 1.6-1 2.6-1 1.8.3 2.5 1a3.4 3.4 0 0 1 .6 4.3 3.7 3.7 0 0 1-3 1.8zm23.2-33.5h6v23.3c0 2.1-.4 4-1.3 5.5a9.1 9.1 0 0 1-3.8 3.5c-1.6.8-3.5 1.3-5.7 1.3-2 0-3.7-.4-5.3-1s-2.8-1.8-3.7-3.2c-.9-1.3-1.4-3-1.4-5h6c.1.8.3 1.6.7 2.2s1 1.2 1.6 1.5c.7.4 1.5.5 2.4.5 1 0 1.8-.2 2.4-.6a4 4 0 0 0 1.6-1.8c.3-.8.5-1.8.5-3V45.5zm30.9 9.1a4.4 4.4 0 0 0-2-3.3 7.5 7.5 0 0 0-4.3-1.1c-1.3 0-2.4.2-3.3.5-.9.4-1.6 1-2 1.6a3.5 3.5 0 0 0-.3 4c.3.5.7.9 1.3 1.2l1.8 1 2 .5 3.2.8c1.3.3 2.5.7 3.7 1.2a13 13 0 0 1 3.2 1.8 8.1 8.1 0 0 1 3 6.5c0 2-.5 3.7-1.5 5.1a10 10 0 0 1-4.4 3.5c-1.8.8-4.1 1.2-6.8 1.2-2.6 0-4.9-.4-6.8-1.2-2-.8-3.4-2-4.5-3.5a10 10 0 0 1-1.7-5.6h6a5 5 0 0 0 3.5 4.6c1 .4 2.2.6 3.4.6 1.3 0 2.5-.2 3.5-.6 1-.4 1.8-1 2.4-1.7a4 4 0 0 0 .8-2.4c0-.9-.2-1.6-.7-2.2a11 11 0 0 0-2.1-1.4l-3.2-1-3.8-1c-2.8-.7-5-1.7-6.6-3.2a7.2 7.2 0 0 1-2.4-5.7 8 8 0 0 1 1.7-5 10 10 0 0 1 4.3-3.5c2-.8 4-1.2 6.4-1.2 2.3 0 4.4.4 6.2 1.2 1.8.8 3.2 2 4.3 3.4 1 1.4 1.5 3 1.5 5h-5.8z"/></svg>
### .\frontend\public\next.svg END ###

### .\frontend\public\vercel.svg BEGIN ###
<svg fill="none" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 1155 1000"><path d="m577.3 0 577.4 1000H0z" fill="#fff"/></svg>
### .\frontend\public\vercel.svg END ###

### .\frontend\public\window.svg BEGIN ###
<svg fill="none" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16"><path fill-rule="evenodd" clip-rule="evenodd" d="M1.5 2.5h13v10a1 1 0 0 1-1 1h-11a1 1 0 0 1-1-1zM0 1h16v11.5a2.5 2.5 0 0 1-2.5 2.5h-11A2.5 2.5 0 0 1 0 12.5zm3.75 4.5a.75.75 0 1 0 0-1.5.75.75 0 0 0 0 1.5M7 4.75a.75.75 0 1 1-1.5 0 .75.75 0 0 1 1.5 0m1.75.75a.75.75 0 1 0 0-1.5.75.75 0 0 0 0 1.5" fill="#666"/></svg>
### .\frontend\public\window.svg END ###

### .\frontend\README.md BEGIN ###
This is a [Next.js](https://nextjs.org) project bootstrapped with [`create-next-app`](https://nextjs.org/docs/app/api-reference/cli/create-next-app).

## Getting Started

First, run the development server:

```bash
npm run dev
# or
yarn dev
# or
pnpm dev
# or
bun dev
```

Open [http://localhost:3000](http://localhost:3000) with your browser to see the result.

You can start editing the page by modifying `app/page.tsx`. The page auto-updates as you edit the file.

This project uses [`next/font`](https://nextjs.org/docs/app/building-your-application/optimizing/fonts) to automatically optimize and load [Geist](https://vercel.com/font), a new font family for Vercel.

## Learn More

To learn more about Next.js, take a look at the following resources:

- [Next.js Documentation](https://nextjs.org/docs) - learn about Next.js features and API.
- [Learn Next.js](https://nextjs.org/learn) - an interactive Next.js tutorial.

You can check out [the Next.js GitHub repository](https://github.com/vercel/next.js) - your feedback and contributions are welcome!

## Deploy on Vercel

The easiest way to deploy your Next.js app is to use the [Vercel Platform](https://vercel.com/new?utm_medium=default-template&filter=next.js&utm_source=create-next-app&utm_campaign=create-next-app-readme) from the creators of Next.js.

Check out our [Next.js deployment documentation](https://nextjs.org/docs/app/building-your-application/deploying) for more details.

### .\frontend\README.md END ###

### .\frontend\src\app\favicon.ico BEGIN ###
         (  F          (  n  00     (-  �         �  �F  (                                                           $   ]   �   �   ]   $                                       �   �   �   �   �   �   �   �                           8   �   �   �   �   �   �   �   �   �   �   8                  �   �   �   �   �   �   �   �   �   �   �   �              �   �   �   �   �   �   �   �   �   �   �   �   �   �       #   �   �   �OOO�������������������������ggg�   �   �   �   #   Y   �   �   ��������������������������555�   �   �   �   Y   �   �   �   �   �kkk���������������������   �   �   �   �   �   �   �   �   �   �			������������������   �   �   �   �   �   Y   �   �   �   �   �JJJ���������kkk�   �   �   �   �   �   Y   #   �   �   �   �   ����������			�   �   �   �   �   �   #       �   �   �   �   �   �111�DDD�   �   �   �   �   �   �              �   �   �   �   �   �   �   �   �   �   �   �                  8   �   �   �   �   �   �   �   �   �   �   8                           �   �   �   �   �   �   �   �                                       $   ]   �   �   ]   $                                                                                                                                                                                                                                                                                    (       @                                                                               ,   U   �   �   �   �   U   ,                                                                                      *   �   �   �   �   �   �   �   �   �   �   �   �   *                                                                      �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �                                                          Q   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   Q                                               r   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   r                                       r   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   r                               O   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   O                          �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �                      �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �               (   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   '           �   �   �   �   �   �   �888���������������������������������������������������������___�   �   �   �   �   �   �   �          �   �   �   �   �   �   ����������������������������������������������������������SSS�   �   �   �   �   �   �   �      +   �   �   �   �   �   �   �   �hhh�����������������������������������������������������   �   �   �   �   �   �   �   �   +   T   �   �   �   �   �   �   �   ��������������������������������������������������,,,�   �   �   �   �   �   �   �   �   T   �   �   �   �   �   �   �   �   �   �GGG���������������������������������������������   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   ������������������������������������������   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �+++���������������������������������jjj�   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   ����������������������������������   �   �   �   �   �   �   �   �   �   �   �   T   �   �   �   �   �   �   �   �   �   �   ��������������������������III�   �   �   �   �   �   �   �   �   �   �   �   T   +   �   �   �   �   �   �   �   �   �   �   �   �hhh����������������������   �   �   �   �   �   �   �   �   �   �   �   +      �   �   �   �   �   �   �   �   �   �   �   ������������������,,,�   �   �   �   �   �   �   �   �   �   �   �   �          �   �   �   �   �   �   �   �   �   �   �   �   �GGG�������������   �   �   �   �   �   �   �   �   �   �   �   �   �           '   �   �   �   �   �   �   �   �   �   �   �   �   ����������   �   �   �   �   �   �   �   �   �   �   �   �   (               �   �   �   �   �   �   �   �   �   �   �   �   �333�___�   �   �   �   �   �   �   �   �   �   �   �   �   �                      �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �                          O   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   O                               r   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   r                                       r   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   r                                               Q   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   Q                                                          �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �                                                                      *   �   �   �   �   �   �   �   �   �   �   �   �   *                                                                                      ,   U   �   �   �   �   U   ,                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                               (   0   `           -                                                                                             	   (   L   j   �   �   �   �   j   K   (   	                                                                                                                                          V   �   �   �   �   �   �   �   �   �   �   �   �   �   �   U                                                                                                                      %   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   &                                                                                                      �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �                                                                                          Q   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   R                                                                              �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �                                                                     �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �                                                             �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �                                                     �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �                                              �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �                                       P   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   O                                  �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �                              �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �                       #   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   #                   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �                  �   �   �   �   �   �   �   �   �   �$$$�hhh�eee�eee�eee�eee�eee�eee�eee�eee�eee�eee�eee�eee�eee�eee�eee�eee�eee�eee�eee�eee�eee�PPP��   �   �   �   �   �   �   �   �   �              U   �   �   �   �   �   �   �   �   �   ������������������������������������������������������������������������������������������sss�   �   �   �   �   �   �   �   �   �   �   U           �   �   �   �   �   �   �   �   �   �   �   �eee��������������������������������������������������������������������������������������   �   �   �   �   �   �   �   �   �   �   �       	   �   �   �   �   �   �   �   �   �   �   �   ����������������������������������������������������������������������������������HHH�   �   �   �   �   �   �   �   �   �   �   �   �   	   (   �   �   �   �   �   �   �   �   �   �   �   �   �EEE�����������������������������������������������������������������������������   �   �   �   �   �   �   �   �   �   �   �   �   �   (   K   �   �   �   �   �   �   �   �   �   �   �   �   �   �������������������������������������������������������������������������,,,�   �   �   �   �   �   �   �   �   �   �   �   �   �   L   j   �   �   �   �   �   �   �   �   �   �   �   �   �   �)))���������������������������������������������������������������������   �   �   �   �   �   �   �   �   �   �   �   �   �   �   j   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   ������������������������������������������������������������������   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   ����������������������������������������������������������iii�   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �eee������������������������������������������������������   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   ��������������������������������������������������HHH�   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   j   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �EEE���������������������������������������������   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   j   L   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �����������������������������������������,,,�   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   K   (   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �)))�������������������������������������   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   (   	   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   ����������������������������������   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   	       �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   ��������������������������iii�   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �           U   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �eee����������������������   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   U              �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   ������������������HHH�   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �                  �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �EEE�������������   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �                   #   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   ���������,,,�   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   #                       �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �222�}}}�   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �                              �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �                                  O   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   P                                       �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �                                              �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �                                                     �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �                                                             �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �                                                                     �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �                                                                              R   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   Q                                                                                          �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �                                                                                                      &   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   �   %                                                                                                                      U   �   �   �   �   �   �   �   �   �   �   �   �   �   �   V                                                                                                                                          	   (   K   j   �   �   �   �   j   L   (   	                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                        �PNG

   IHDR         \r�f   sRGB ���   8eXIfMM *    �i            �       �           D"8s  IDATx�]	�ՙn�]<QVA���h$	�N��13*�q��d�č�I���D�L2��(�(Ԙ2�ę�G	��q_@屈���xț�Џ��{o�������U�{}�O��;������9�d���(Dg��8	��N �]��@�hx�?v �N�3�=`;�6�.�&��u��  ��6�P��н��@�àR� P�iZq�^DN���wp����X�hИHg@��
:��|�5` p"@�'�ɲ�s{�p�*�2����� d ү���|(0�
0 ��>K�
�xX�6 IJ� �C|?$KEN�}ϓ|������h $	2 ��|/� . Nz �#���W�e�
�5������ܶ���;�y �� �g�s�h^  I�� DL(�;�8��Hjg�cH|x�1��R"�a���Ӂ� G��@��9`/`%0�H�@j�~,���K
�,t).��I���D�T�O�)~��V�u$b 誛�U%�7������ _�$b 8A������J�3` 510wQ�?��vr���:�2�K�@ ��v*{%#��A�Z�咁^(��=�g \��W�����!:��,`�6��643�:@�c.Fٟ����u?�<��'������_܏vp: �8Q��
I�Ł�p{3���kHȢ�G�����c�Ѽ<�62&�
��2uC�����敭��T�3�
�����;���d�/~m��.��X�@{�w.��d]G�� {lK��Eb���(P�RuM�T�C�����d��])��_Lm�=��=@b���K��GUk�^�U�������)1����g�T���m`9�\����Q��@����Ⱆ6�:ڞ�^�w�����E�D�� �	�5����F�,��
�X"�d�m�<�nB~��@����t�t�x���;�f�>����I8����8��C1۪$B���e���+��jl��EZ��& ��S:�:�6�m����\G1��`���!�nl�l�Ɗ�^�Q`��@Oc�S��@e�ͷ���qb�p���S��@up���F�D@�Г������2@#����L3 �A��$H2� _h��FH#rq(��O�D�򤬈���runGOWa�b� &�SgD�3�ED�to�*Ǥ����9k��~)���,$� x�R�1�v�K ��9�D䍁U(�w�&LE��ꩻ�S)��3�Y8x8 $.i�(��K�ŀY����a�]����4��ǀ	c����@3�f����4� Ƣ���/*b��� ���$!I�~��7�B*-1`	o � �	�$��ǡD�����L������ �J"���OQ��)��2@#�x4�"$e ���I�8��Oi��8�"� �G��8[x�t<�.��7&�m&؎R�^��tq� ؕ�.���Y�-2� �d� ��*_��&d|j\�W�b ��G����*g�� ��釁�F4�"I�؃�/ b1q�N����Y�D��p���9���p�}w\� �Ԥ���1 j`��O���xK=��H�� �A��1�#�
D:U8j���t���$b b�A||�U�Q��26%��)1 ��_�ꢳ!~D��� ��+b >A��:]�E$��50��GDhR�t����ݻwR�)��P� ��n$� 3���@bS�Nu�,Y�j�ʲ��:����;�����@�`�|�-[)�'OV��Ն�sFxڮ��ۥ�n}͛7�����~��ƺ�:���Q��J_��UKj8�q0x���;v4 ̞=[�hW=�	��	�&�!e5�8hѢE��w�]�����6���_�iW}�SZ�?	�/`�;vl�}��2 <�h�" ����A�܁�X,�m۶�+V�(��<�w���#F�^���;���aH�c ���)S�*�{a���p��c89(�^����4�&E��oÆ��W�/��u�=�^���*?{k^�_E�����z���g�� UI-���{WU*
�:p�9.tڷo(/ݺus>��3�'�^�Rg���ڞG��I_D�������~~� ��{���?N0�7�S��.ƍ׸�~?}/y]nA;�أ���2 ]�FOB2C?�_I����[�:�:�=#�OzK�-� ��ϣ�%����?j��I���P�ۯ��{N�-hU��t�:������� ,���G�K�-hU���c�hP7 ����@�n?�\�-�k�.���2�:�� �`��F��=�-�V�_�G��܂V� ��}�0 WI����F��ʭ���sM�rZ�8pJ�Q�*@OK8���
rZ��ݖa, ��w� �S�W^y����.��5�at7��ݏ���Tv#�~7n��A"�����+��W��pM��/�hK8����g��F/^������M{e ��R�|�)q��7�t��?8'���K��P~���瞰�\��r��>�ǷUk �eP��|�^x����
�/V/��v���������*�p�v�� ����ʟ]J��}��k8(������ĉ�ѣGǗ�O�mڴq,X�o���e.�^ �Qx���p�t����4^_�N�{�����y�2 �s����� �-عsg�s���i�v��Z8
!~PJ?�c�������|�] �ܽ{��z�긓R��1pn���z�����tlp�9�f�r�v�jT殿�z�4*O�L�~����ԕ3��4�~~�r�;�m�xY�+���������3 r�;�m�x�4���:7]ՁqL�4)U��!r�1��u�6���$��7����8�w��̙3Ǹ|5�>?�\z��O���͆� ��,�E����3�����2���[����2Wu:E�����^p.H1cJ�t�]}��B�u��SOu�����Ic�O�����%�  �AZ������k����D?�5 �@Q�����3�w�+��"��T��S��Uޥ�13��?��5 M'݋��>p��Z�j�~fj�׈�סԐ�n�����>� ��i5D�[bf ��~a�'�`Xc��� -�1�k����āI�������k��Q�ů|�k�M��(92�@�t�����݂X-�Lדa��N4��qܞ'$f0@�@V�nA�ܘY�L9:�|/^s� ��	��)0`�j��T\w�uZ-����¨\�	@�:��c�t���{�-��Rb��1%� �I,Y%T���~��r�1����C��,�$��*ˀ���f<��0z����h�F���� ����|���8Z-�CR����Tg� �HRf��glY����s��-��p��'+����m�_ؒg������C�{ �	����Ȫ�ϏΙ3g�-�GR|׹7`G��񥡘�0�U��_ٵZЏ�د�D�)���\>����ʗ������z N���@��~~��-��P��{rs���@�<����|.]�Ը|��m|g����_��y�W�KD1�b�M���%�s\����r�1��n�\�ƒ�"-� �`.4��~%3��I}[0A��$��= -�>BH"G�ۏ�^r��<�EBG�i �%���9�@^�~~@�����1����@� t�-[����{%@C�$�mAg���Κ5kʆх����/双O��l��ӿ��B�@.X���u�p�O��6��x�9MPn�`߷o_���^n�`t�
��(�����\r��s�A�y���ۂ�T��@h
�E0l�0��;�tڵӘkƸN����Y�jU��
S#�|^㽺- |��p�N�.���ޥ`�^{�zL�6��4 �ě�b��e�]&"�d�sΜ9Uޥ�U0�!��*nP�*`���o֨v����i8G�����hh��m������ɓ�s�=�{J�U0�Ղ���wZ������������8bEz���,Y�D��![C�>}��7:k׮�no��f� >jvR?#b��X�(��F�AT�F��i��[�{��zv��>��C���a+�[0B2�D��=��G~�(
�ĺ������LO�\s�܂>"8|�`[)
&Lp8�'��������4 oGe�#�ۏ�lْ_\�D̀܂�2Z�l��i�9��t�ȑ9f ޢ�-����=���Y�y��n?uQ�}Xͬ�sA�i >=��1�=R��+� +�܂��.2 ��K������CƢۃ20h� �˫%53�5@�MA�%���̣������j[��9�;�� _(�����0��~r���\�{�m�P����x#TT9��n?����N#��ץ&�}� ��)
�T�VL�!���j���`�p �8@Rr�UAV�A����=��-����pLH�`@n�*Ȋ1�܂U���?}w ]�H2@�ߴi��V���[�˯%�������5 �8�)Э
T`��|rZbZ-�.�!da+@� ���ߞ�Z�gf�[0p���� �� I��gr�$��o%P�_rCy�V�|߽����"m�Y���-�[ l��k xA� ��ۯ9]�[pҤI�Ȩ�pP���k ��Feِ���gHE�d�nAm"Z�$��5}���z�8����2r�X�|� ��Sܻw��r�J�s�J�~�T�f�z{ �ͫ ��x�j?j��Q�E�n� �js���|G�xз�<dXt(��Q�E�.�p�47 ��)���;��ys�_�V�D���-XTi����?� �~�薜����� �`Q�=V�?���^�
������.]�|X�
�m�B~��?���J� �D�������~�h r�����ER���A݀�B���~w�q�Ӿ}���<�ŕ[й5�d��-�`�5 ?�Kq�~l4��0@��)����/I��(����؋���n��9���Y�4�!�Cو2ח*w9���GKݐ�s�&�r�e��s��?�6�8J� |(�uwO䴁d�&K)�nA��?R���n@7,��8�=���r�e����n�M�69k��M7�����J��R�]�e�n��9���Z���� /?នo>��󕾤�rzr�� ��`���V{���u��4448�V��ra��p� ��QRZ�<{�dK.F9��#~T���s.����N%*� ���Ýu�8G&����/W:*x%�{�}@� ��l���Nc#�AI�������i����*?�د�0}�g���C"Āpۯ������4薒ҏ(b�8�_Q�Y� ���r7'���`��� �j �6�� *��3�W�g��"��l��1�:�Sg}%� �	��P?����1`�����Y� ��"��D�0b@�� �����9������[t��F1���p`k�\U�`��R��A#W81 e`)R�ZM��� ��[u��F0�	rq.����� #^�=C"Ā9P'�R~f�� �
pn�zdC"�e���?�\K����@&$b }jz�3۵� x/{��1 Ra�#�|��ƟUK�= &�^��TM�n�2�9�5)?s���{O'�D��D���o [kM�oK0�x�� �Td�_@]b r� �G�����; ����D��D���1�gaR�`��'`0�  �>\��/���f��������ŀ����!fn�Z�|b����U�.t���ट���r�9�+��������	�b rnE�Dk�= ��8�����!b R�Cl�P�E�`�܌�K�'~�@���}*�!`�@��6 L��;��	$b@D��?#��g�F�
��V��1�v��;�Es��Q����=ɮ�4���b@T��n��!��3q�0^�V�� c ��1�ܶ��[����M�=8I����1@�څ@Cu��`N�o�� WJĀ� W����e��I�� n��N�mீ��ܴ�_d��(�4`E܅I� ���"̵�1 *3�+\�E� �\M���)g	r���
���8�>��p�?vI� �0�ǀ~�!b������$'�%"I����R��i�1 �0��? S~&�� �r�����{ n�_�����L�?��T�e��Ǝ�7�C"r��OQ~"qI� ��O 8�?$b �܋r�#@�_�v�J̙��/��3�'d�/����W[����o'N��l��-2� ���@j�O~��0���2` H�@�؄��+����pOB� �uO��(l�S�ԕ���9����~�c�:x/�Xd�.���Ɣ�d ��V�y@F $H2� ����+M*�i��l8O@F $H2� ���2�4& r�PO��֢����7N�YS ����Y�1`��;�JS3n� g[�'��@W@"la`32�n?'�HB2p
�hām�mu �����j@F@��V����Z!��xI���H�y�ѱ)��>��Z!6 ���a�`�����dDV$9f���	pM�6�I�!LG:\LdrwPy�~�P�%��L3��7�TK��Am�mo|�6��	3��-�h J3��?�67 �yr���"����g��4. $�1���_�[*��&���S/�dq�������C��h �3��>�6Ŷ%������\�#�RZq��=lK|ŔX��X�WS�e j5 /����$���:��v@������8���d��1(�z2~F�)���3��͋���l��C�������#����=�.\Lt? %� N$9b�%�:���2��u	 �1|-�	ld�����t $b��@?���@� �F�c��ρ^�D�d�[9�ࠐz�����:
H�@ ��P2v )~���@����z5��|����R�ֵ���|`#�W39؂��<�"-�0��\<�d��u�oGLz 1��Gp����e�倯d� .�jH�@j�F�3��@ c{s<��J&	�@�����b���w��  �� ��n���v��< �����,M;��*p>p!0hH��{=�����x�]I�� DLh����<'��h8�@V �#��J���f� I�� �Hn����W�}�N�t[u�$�������� �@� 2 	�]&)�� #�3���,	=%�T���k�&�  I�����I��ӳ� �[8	�	�L�]�]t�T�g���6�-@b2 U�OV��: A?��} .i�|	�xC���rv�w; ��#�>�i 8_b82 �WP����� �� {'n���8�z;�Ƥy��s� ��@���P��o|�S�ih $3��@߹j��    IEND�B`�
### .\frontend\src\app\favicon.ico END ###

### .\frontend\src\app\globals.css BEGIN ###
@import "tailwindcss";
@import "tw-animate-css";

@custom-variant dark (&:is(.dark *));

@theme inline {
  --color-background: var(--background);
  --color-foreground: var(--foreground);
  --font-sans: var(--font-geist-sans);
  --font-mono: var(--font-geist-mono);
  --color-sidebar-ring: var(--sidebar-ring);
  --color-sidebar-border: var(--sidebar-border);
  --color-sidebar-accent-foreground: var(--sidebar-accent-foreground);
  --color-sidebar-accent: var(--sidebar-accent);
  --color-sidebar-primary-foreground: var(--sidebar-primary-foreground);
  --color-sidebar-primary: var(--sidebar-primary);
  --color-sidebar-foreground: var(--sidebar-foreground);
  --color-sidebar: var(--sidebar);
  --color-chart-5: var(--chart-5);
  --color-chart-4: var(--chart-4);
  --color-chart-3: var(--chart-3);
  --color-chart-2: var(--chart-2);
  --color-chart-1: var(--chart-1);
  --color-ring: var(--ring);
  --color-input: var(--input);
  --color-border: var(--border);
  --color-destructive: var(--destructive);
  --color-accent-foreground: var(--accent-foreground);
  --color-accent: var(--accent);
  --color-muted-foreground: var(--muted-foreground);
  --color-muted: var(--muted);
  --color-secondary-foreground: var(--secondary-foreground);
  --color-secondary: var(--secondary);
  --color-primary-foreground: var(--primary-foreground);
  --color-primary: var(--primary);
  --color-popover-foreground: var(--popover-foreground);
  --color-popover: var(--popover);
  --color-card-foreground: var(--card-foreground);
  --color-card: var(--card);
  --radius-sm: calc(var(--radius) - 4px);
  --radius-md: calc(var(--radius) - 2px);
  --radius-lg: var(--radius);
  --radius-xl: calc(var(--radius) + 4px);
}

:root {
  --radius: 0.625rem;
  --background: oklch(1 0 0);
  --foreground: oklch(0.145 0 0);
  --card: oklch(1 0 0);
  --card-foreground: oklch(0.145 0 0);
  --popover: oklch(1 0 0);
  --popover-foreground: oklch(0.145 0 0);
  --primary: oklch(0.205 0 0);
  --primary-foreground: oklch(0.985 0 0);
  --secondary: oklch(0.97 0 0);
  --secondary-foreground: oklch(0.205 0 0);
  --muted: oklch(0.97 0 0);
  --muted-foreground: oklch(0.556 0 0);
  --accent: oklch(0.97 0 0);
  --accent-foreground: oklch(0.205 0 0);
  --destructive: oklch(0.577 0.245 27.325);
  --border: oklch(0.922 0 0);
  --input: oklch(0.922 0 0);
  --ring: oklch(0.708 0 0);
  --chart-1: oklch(0.646 0.222 41.116);
  --chart-2: oklch(0.6 0.118 184.704);
  --chart-3: oklch(0.398 0.07 227.392);
  --chart-4: oklch(0.828 0.189 84.429);
  --chart-5: oklch(0.769 0.188 70.08);
  --sidebar: oklch(0.985 0 0);
  --sidebar-foreground: oklch(0.145 0 0);
  --sidebar-primary: oklch(0.205 0 0);
  --sidebar-primary-foreground: oklch(0.985 0 0);
  --sidebar-accent: oklch(0.97 0 0);
  --sidebar-accent-foreground: oklch(0.205 0 0);
  --sidebar-border: oklch(0.922 0 0);
  --sidebar-ring: oklch(0.708 0 0);
}

.dark {
  --background: oklch(0.145 0 0);
  --foreground: oklch(0.985 0 0);
  --card: oklch(0.205 0 0);
  --card-foreground: oklch(0.985 0 0);
  --popover: oklch(0.205 0 0);
  --popover-foreground: oklch(0.985 0 0);
  --primary: oklch(0.922 0 0);
  --primary-foreground: oklch(0.205 0 0);
  --secondary: oklch(0.269 0 0);
  --secondary-foreground: oklch(0.985 0 0);
  --muted: oklch(0.269 0 0);
  --muted-foreground: oklch(0.708 0 0);
  --accent: oklch(0.269 0 0);
  --accent-foreground: oklch(0.985 0 0);
  --destructive: oklch(0.704 0.191 22.216);
  --border: oklch(1 0 0 / 10%);
  --input: oklch(1 0 0 / 15%);
  --ring: oklch(0.556 0 0);
  --chart-1: oklch(0.488 0.243 264.376);
  --chart-2: oklch(0.696 0.17 162.48);
  --chart-3: oklch(0.769 0.188 70.08);
  --chart-4: oklch(0.627 0.265 303.9);
  --chart-5: oklch(0.645 0.246 16.439);
  --sidebar: oklch(0.205 0 0);
  --sidebar-foreground: oklch(0.985 0 0);
  --sidebar-primary: oklch(0.488 0.243 264.376);
  --sidebar-primary-foreground: oklch(0.985 0 0);
  --sidebar-accent: oklch(0.269 0 0);
  --sidebar-accent-foreground: oklch(0.985 0 0);
  --sidebar-border: oklch(1 0 0 / 10%);
  --sidebar-ring: oklch(0.556 0 0);
}

@layer base {
  * {
    @apply border-border outline-ring/50;
  }
  body {
    @apply bg-background text-foreground;
  }
}

### .\frontend\src\app\globals.css END ###

### .\frontend\src\app\layout.tsx BEGIN ###
import type { Metadata } from "next";
import { Geist, Geist_Mono } from "next/font/google";
import "./globals.css";

const geistSans = Geist({
  variable: "--font-geist-sans",
  subsets: ["latin"],
});

const geistMono = Geist_Mono({
  variable: "--font-geist-mono",
  subsets: ["latin"],
});

export const metadata: Metadata = {
  title: "Create Next App",
  description: "Generated by create next app",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body
        className={`${geistSans.variable} ${geistMono.variable} antialiased`}
      >
        {children}
      </body>
    </html>
  );
}

### .\frontend\src\app\layout.tsx END ###

### .\frontend\src\app\page.tsx BEGIN ###
import Image from "next/image";

export default function Home() {
  return (
    <div className="flex min-h-screen items-center justify-center bg-zinc-50 font-sans dark:bg-black">
      <main className="flex min-h-screen w-full max-w-3xl flex-col items-center justify-between py-32 px-16 bg-white dark:bg-black sm:items-start">
        <Image
          className="dark:invert"
          src="/next.svg"
          alt="Next.js logo"
          width={100}
          height={20}
          priority
        />
        <div className="flex flex-col items-center gap-6 text-center sm:items-start sm:text-left">
          <h1 className="max-w-xs text-3xl font-semibold leading-10 tracking-tight text-black dark:text-zinc-50">
            To get started, edit the page.tsx file.
          </h1>
          <p className="max-w-md text-lg leading-8 text-zinc-600 dark:text-zinc-400">
            Looking for a starting point or more instructions? Head over to{" "}
            <a
              href="https://vercel.com/templates?framework=next.js&utm_source=create-next-app&utm_medium=appdir-template-tw&utm_campaign=create-next-app"
              className="font-medium text-zinc-950 dark:text-zinc-50"
            >
              Templates
            </a>{" "}
            or the{" "}
            <a
              href="https://nextjs.org/learn?utm_source=create-next-app&utm_medium=appdir-template-tw&utm_campaign=create-next-app"
              className="font-medium text-zinc-950 dark:text-zinc-50"
            >
              Learning
            </a>{" "}
            center.
          </p>
        </div>
        <div className="flex flex-col gap-4 text-base font-medium sm:flex-row">
          <a
            className="flex h-12 w-full items-center justify-center gap-2 rounded-full bg-foreground px-5 text-background transition-colors hover:bg-[#383838] dark:hover:bg-[#ccc] md:w-[158px]"
            href="https://vercel.com/new?utm_source=create-next-app&utm_medium=appdir-template-tw&utm_campaign=create-next-app"
            target="_blank"
            rel="noopener noreferrer"
          >
            <Image
              className="dark:invert"
              src="/vercel.svg"
              alt="Vercel logomark"
              width={16}
              height={16}
            />
            Deploy Now
          </a>
          <a
            className="flex h-12 w-full items-center justify-center rounded-full border border-solid border-black/[.08] px-5 transition-colors hover:border-transparent hover:bg-black/[.04] dark:border-white/[.145] dark:hover:bg-[#1a1a1a] md:w-[158px]"
            href="https://nextjs.org/docs?utm_source=create-next-app&utm_medium=appdir-template-tw&utm_campaign=create-next-app"
            target="_blank"
            rel="noopener noreferrer"
          >
            Documentation
          </a>
        </div>
      </main>
    </div>
  );
}

### .\frontend\src\app\page.tsx END ###

### .\frontend\src\components\ui\button.tsx BEGIN ###
import * as React from "react"
import { Slot } from "@radix-ui/react-slot"
import { cva, type VariantProps } from "class-variance-authority"

import { cn } from "@/lib/utils"

const buttonVariants = cva(
  "inline-flex items-center justify-center gap-2 whitespace-nowrap rounded-md text-sm font-medium transition-all disabled:pointer-events-none disabled:opacity-50 [&_svg]:pointer-events-none [&_svg:not([class*='size-'])]:size-4 shrink-0 [&_svg]:shrink-0 outline-none focus-visible:border-ring focus-visible:ring-ring/50 focus-visible:ring-[3px] aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 aria-invalid:border-destructive",
  {
    variants: {
      variant: {
        default: "bg-primary text-primary-foreground hover:bg-primary/90",
        destructive:
          "bg-destructive text-white hover:bg-destructive/90 focus-visible:ring-destructive/20 dark:focus-visible:ring-destructive/40 dark:bg-destructive/60",
        outline:
          "border bg-background shadow-xs hover:bg-accent hover:text-accent-foreground dark:bg-input/30 dark:border-input dark:hover:bg-input/50",
        secondary:
          "bg-secondary text-secondary-foreground hover:bg-secondary/80",
        ghost:
          "hover:bg-accent hover:text-accent-foreground dark:hover:bg-accent/50",
        link: "text-primary underline-offset-4 hover:underline",
      },
      size: {
        default: "h-9 px-4 py-2 has-[>svg]:px-3",
        sm: "h-8 rounded-md gap-1.5 px-3 has-[>svg]:px-2.5",
        lg: "h-10 rounded-md px-6 has-[>svg]:px-4",
        icon: "size-9",
        "icon-sm": "size-8",
        "icon-lg": "size-10",
      },
    },
    defaultVariants: {
      variant: "default",
      size: "default",
    },
  }
)

function Button({
  className,
  variant,
  size,
  asChild = false,
  ...props
}: React.ComponentProps<"button"> &
  VariantProps<typeof buttonVariants> & {
    asChild?: boolean
  }) {
  const Comp = asChild ? Slot : "button"

  return (
    <Comp
      data-slot="button"
      className={cn(buttonVariants({ variant, size, className }))}
      {...props}
    />
  )
}

export { Button, buttonVariants }

### .\frontend\src\components\ui\button.tsx END ###

### .\frontend\src\components\ui\form.tsx BEGIN ###
"use client"

import * as React from "react"
import * as LabelPrimitive from "@radix-ui/react-label"
import { Slot } from "@radix-ui/react-slot"
import {
  Controller,
  FormProvider,
  useFormContext,
  useFormState,
  type ControllerProps,
  type FieldPath,
  type FieldValues,
} from "react-hook-form"

import { cn } from "@/lib/utils"
import { Label } from "@/components/ui/label"

const Form = FormProvider

type FormFieldContextValue<
  TFieldValues extends FieldValues = FieldValues,
  TName extends FieldPath<TFieldValues> = FieldPath<TFieldValues>,
> = {
  name: TName
}

const FormFieldContext = React.createContext<FormFieldContextValue>(
  {} as FormFieldContextValue
)

const FormField = <
  TFieldValues extends FieldValues = FieldValues,
  TName extends FieldPath<TFieldValues> = FieldPath<TFieldValues>,
>({
  ...props
}: ControllerProps<TFieldValues, TName>) => {
  return (
    <FormFieldContext.Provider value={{ name: props.name }}>
      <Controller {...props} />
    </FormFieldContext.Provider>
  )
}

const useFormField = () => {
  const fieldContext = React.useContext(FormFieldContext)
  const itemContext = React.useContext(FormItemContext)
  const { getFieldState } = useFormContext()
  const formState = useFormState({ name: fieldContext.name })
  const fieldState = getFieldState(fieldContext.name, formState)

  if (!fieldContext) {
    throw new Error("useFormField should be used within <FormField>")
  }

  const { id } = itemContext

  return {
    id,
    name: fieldContext.name,
    formItemId: `${id}-form-item`,
    formDescriptionId: `${id}-form-item-description`,
    formMessageId: `${id}-form-item-message`,
    ...fieldState,
  }
}

type FormItemContextValue = {
  id: string
}

const FormItemContext = React.createContext<FormItemContextValue>(
  {} as FormItemContextValue
)

function FormItem({ className, ...props }: React.ComponentProps<"div">) {
  const id = React.useId()

  return (
    <FormItemContext.Provider value={{ id }}>
      <div
        data-slot="form-item"
        className={cn("grid gap-2", className)}
        {...props}
      />
    </FormItemContext.Provider>
  )
}

function FormLabel({
  className,
  ...props
}: React.ComponentProps<typeof LabelPrimitive.Root>) {
  const { error, formItemId } = useFormField()

  return (
    <Label
      data-slot="form-label"
      data-error={!!error}
      className={cn("data-[error=true]:text-destructive", className)}
      htmlFor={formItemId}
      {...props}
    />
  )
}

function FormControl({ ...props }: React.ComponentProps<typeof Slot>) {
  const { error, formItemId, formDescriptionId, formMessageId } = useFormField()

  return (
    <Slot
      data-slot="form-control"
      id={formItemId}
      aria-describedby={
        !error
          ? `${formDescriptionId}`
          : `${formDescriptionId} ${formMessageId}`
      }
      aria-invalid={!!error}
      {...props}
    />
  )
}

function FormDescription({ className, ...props }: React.ComponentProps<"p">) {
  const { formDescriptionId } = useFormField()

  return (
    <p
      data-slot="form-description"
      id={formDescriptionId}
      className={cn("text-muted-foreground text-sm", className)}
      {...props}
    />
  )
}

function FormMessage({ className, ...props }: React.ComponentProps<"p">) {
  const { error, formMessageId } = useFormField()
  const body = error ? String(error?.message ?? "") : props.children

  if (!body) {
    return null
  }

  return (
    <p
      data-slot="form-message"
      id={formMessageId}
      className={cn("text-destructive text-sm", className)}
      {...props}
    >
      {body}
    </p>
  )
}

export {
  useFormField,
  Form,
  FormItem,
  FormLabel,
  FormControl,
  FormDescription,
  FormMessage,
  FormField,
}

### .\frontend\src\components\ui\form.tsx END ###

### .\frontend\src\components\ui\label.tsx BEGIN ###
"use client"

import * as React from "react"
import * as LabelPrimitive from "@radix-ui/react-label"

import { cn } from "@/lib/utils"

function Label({
  className,
  ...props
}: React.ComponentProps<typeof LabelPrimitive.Root>) {
  return (
    <LabelPrimitive.Root
      data-slot="label"
      className={cn(
        "flex items-center gap-2 text-sm leading-none font-medium select-none group-data-[disabled=true]:pointer-events-none group-data-[disabled=true]:opacity-50 peer-disabled:cursor-not-allowed peer-disabled:opacity-50",
        className
      )}
      {...props}
    />
  )
}

export { Label }

### .\frontend\src\components\ui\label.tsx END ###

### .\frontend\src\components\ui\table.tsx BEGIN ###
"use client"

import * as React from "react"

import { cn } from "@/lib/utils"

function Table({ className, ...props }: React.ComponentProps<"table">) {
  return (
    <div
      data-slot="table-container"
      className="relative w-full overflow-x-auto"
    >
      <table
        data-slot="table"
        className={cn("w-full caption-bottom text-sm", className)}
        {...props}
      />
    </div>
  )
}

function TableHeader({ className, ...props }: React.ComponentProps<"thead">) {
  return (
    <thead
      data-slot="table-header"
      className={cn("[&_tr]:border-b", className)}
      {...props}
    />
  )
}

function TableBody({ className, ...props }: React.ComponentProps<"tbody">) {
  return (
    <tbody
      data-slot="table-body"
      className={cn("[&_tr:last-child]:border-0", className)}
      {...props}
    />
  )
}

function TableFooter({ className, ...props }: React.ComponentProps<"tfoot">) {
  return (
    <tfoot
      data-slot="table-footer"
      className={cn(
        "bg-muted/50 border-t font-medium [&>tr]:last:border-b-0",
        className
      )}
      {...props}
    />
  )
}

function TableRow({ className, ...props }: React.ComponentProps<"tr">) {
  return (
    <tr
      data-slot="table-row"
      className={cn(
        "hover:bg-muted/50 data-[state=selected]:bg-muted border-b transition-colors",
        className
      )}
      {...props}
    />
  )
}

function TableHead({ className, ...props }: React.ComponentProps<"th">) {
  return (
    <th
      data-slot="table-head"
      className={cn(
        "text-foreground h-10 px-2 text-left align-middle font-medium whitespace-nowrap [&:has([role=checkbox])]:pr-0 [&>[role=checkbox]]:translate-y-[2px]",
        className
      )}
      {...props}
    />
  )
}

function TableCell({ className, ...props }: React.ComponentProps<"td">) {
  return (
    <td
      data-slot="table-cell"
      className={cn(
        "p-2 align-middle whitespace-nowrap [&:has([role=checkbox])]:pr-0 [&>[role=checkbox]]:translate-y-[2px]",
        className
      )}
      {...props}
    />
  )
}

function TableCaption({
  className,
  ...props
}: React.ComponentProps<"caption">) {
  return (
    <caption
      data-slot="table-caption"
      className={cn("text-muted-foreground mt-4 text-sm", className)}
      {...props}
    />
  )
}

export {
  Table,
  TableHeader,
  TableBody,
  TableFooter,
  TableHead,
  TableRow,
  TableCell,
  TableCaption,
}

### .\frontend\src\components\ui\table.tsx END ###

### .\frontend\src\lib\utils.ts BEGIN ###
import { clsx, type ClassValue } from "clsx"
import { twMerge } from "tailwind-merge"

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

### .\frontend\src\lib\utils.ts END ###

### .\frontend\tsconfig.json BEGIN ###
{
  "compilerOptions": {
    "target": "ES2017",
    "lib": ["dom", "dom.iterable", "esnext"],
    "allowJs": true,
    "skipLibCheck": true,
    "strict": true,
    "noEmit": true,
    "esModuleInterop": true,
    "module": "esnext",
    "moduleResolution": "bundler",
    "resolveJsonModule": true,
    "isolatedModules": true,
    "jsx": "react-jsx",
    "incremental": true,
    "plugins": [
      {
        "name": "next"
      }
    ],
    "paths": {
      "@/*": ["./src/*"]
    }
  },
  "include": [
    "next-env.d.ts",
    "**/*.ts",
    "**/*.tsx",
    ".next/types/**/*.ts",
    ".next/dev/types/**/*.ts",
    "**/*.mts"
  ],
  "exclude": ["node_modules"]
}

### .\frontend\tsconfig.json END ###

### .\ml\README.md BEGIN ###
machine learning stuff here

### .\ml\README.md END ###

### .\project.md BEGIN ###
### DIRECTORY . FOLDER STRUCTURE ###
FILE .env
FILE .gitignore
DIR backend/
    DIR bot/
        FILE app.log
        DIR cmd/
            FILE main.go
        DIR config/
            FILE config.yaml
        FILE Dockerfile
        FILE go.mod
        FILE go.sum
        DIR internal/
            DIR app/
                FILE app.go
                DIR bot/
                    FILE app.go
                DIR grpc/
                    FILE app.go
            DIR config/
                FILE config.go
            DIR handlers/
                DIR bot/
                    FILE server.go
                DIR grpc/
                    FILE server.go
            DIR models/
                FILE model.go
            DIR repository/
                DIR bot/
                    FILE bot.go
            DIR services/
                FILE bot.go
    DIR database/
        DIR cmd/
            FILE main.go
        DIR config/
            FILE config.yaml
        FILE Dockerfile
        FILE go.mod
        FILE go.sum
        DIR internal/
            DIR app/
                FILE app.go
                DIR grpc/
                    FILE app.go
            DIR config/
                FILE config.go
            DIR handlers/
                DIR grpc/
                    FILE server.go
            DIR migrations/
                FILE 1_init_schema.down.sql
                FILE 1_init_schema.up.sql
            DIR migrator/
                FILE migrator.go
            DIR models/
                FILE models.go
            DIR repository/
                DIR database/
                    FILE db.go
            DIR services/
                FILE db.go
    DIR database_save/
    DIR manager/
        DIR cmd/
            FILE main.go
        DIR config/
            FILE config.yaml
        FILE Dockerfile
        FILE go.mod
        FILE go.sum
        DIR internal/
            DIR app/
                FILE app.go
                DIR http/
                    FILE app.go
            DIR config/
                FILE config.go
            DIR models/
                FILE models.go
            DIR repository/
                DIR bot/
                    FILE client.go
                DIR database/
                    FILE client.go
            DIR router/
                FILE router.go
            DIR services/
                FILE http.go
    FILE README.md
FILE docker-compose.yml
DIR frontend/
    FILE .gitignore
    FILE components.json
    FILE next.config.ts
    FILE package.json
    FILE pnpm-lock.yaml
    FILE postcss.config.mjs
    DIR public/
        FILE file.svg
        FILE globe.svg
        FILE next.svg
        FILE vercel.svg
        FILE window.svg
    FILE README.md
    DIR src/
        DIR app/
            FILE favicon.ico
            FILE globals.css
            FILE layout.tsx
            FILE page.tsx
        DIR components/
            DIR ui/
                FILE button.tsx
                FILE form.tsx
                FILE label.tsx
                FILE table.tsx
        DIR lib/
            FILE utils.ts
    FILE tsconfig.json
DIR ml/
    FILE README.md
FILE project.md
FILE README.md
### DIRECTORY . FOLDER STRUCTURE ###

### DIRECTORY . FLATTENED CONTENT ###

### .\project.md END ###

### .\README.md BEGIN ###
# cybergarden2025-all-in
Кейс Центр Инвеста

### .\README.md END ###

### DIRECTORY . FLATTENED CONTENT ###

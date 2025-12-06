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

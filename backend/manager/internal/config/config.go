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

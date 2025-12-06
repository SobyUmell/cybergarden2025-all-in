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

package config

import (
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string        `yaml:"env" env-default:"local"`
	StoragePath string        `yaml:"storage_path" env-required:"true"`
	HTTPServer  HTTPServer    `yaml:"http_server"`
	Clients     ClientsConfig `yaml:"clients"`
	AppSecret   string        `yaml:"app_secret" env-required:"true" env:"APP_SECRET"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
	User        string        `yaml:"user" env-required:"true"`
	Password    string        `yaml:"password" env-required:"true" env:"HTTP_SERVER_PASSWORD"`
}

type Clients struct {
	Address      string        `yaml:"address"`
	Timeout      time.Duration `yaml:"timeout"`
	RetriesCount int           `yaml:"retriesCount"`
}

type ClientsConfig struct {
	SSO Clients `yaml:"sso"`
}

func MustLoad() (*Config, error) {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "./config/local.yaml"
		fmt.Println("⚙️  CONFIG_PATH not set, using default:", configPath)
	}
	//check if file exist

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file does not exist: %v", configPath)
	}

	var cfg Config
	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		return nil, fmt.Errorf("cannot read condig: %w", err)
	}

	return &cfg, nil
}

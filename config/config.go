package config

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		App      App      `yaml:"app"`
		HTTP     HTTP     `yaml:"http"`
		Postgres Postgres `yaml:"postgres"`
		Log      Log      `yaml:"logger"`
	}

	App struct {
		Name    string `env-required:"true" yaml:"name" env:"APP_NAME"`
		Versoin string `env-required:"true" yaml:"version" env:"APP_VERSION"`
	}

	HTTP struct {
		Port string `env-required:"true" yaml:"port" env:"SERVER_PORT"`
	}

	Postgres struct {
		URL            string        `env-required:"true" yaml:"url" env:"POSTGRES_URL"`
		ConnectTimeout time.Duration `env-required:"true" yaml:"connect_timeout" env:"POSTGRES_CONNECT_TIMEOUT"`
	}
	Log struct {
		Level string `env-required:"true" yaml:"level" env:"LOG_LEVEL"`
	}
)

func New(configPath string) (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig(configPath, cfg)
	if err != nil {
		return nil, fmt.Errorf("config - NewConfig - cleanenv.ReadConfig: %w", err)
	}

	err = cleanenv.UpdateEnv(cfg)
	if err != nil {
		return nil, fmt.Errorf("config - NewConfig - cleanenv.UpdateEnv: %w", err)
	}

	return cfg, nil
}

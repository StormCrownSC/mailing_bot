package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"path"
)

type Config struct {
	Postgres Postgres `yaml:"postgres"`
	Logger   ZeroLog  `yaml:"zero_log"`
	Telegram string   `env:"TG_TOKEN"`
}

type Postgres struct {
	Host     string `env:"POSTGRES_HOST" env-required:"true"`
	Port     string `env:"POSTGRES_PORT" env-required:"true"`
	User     string `env:"POSTGRES_USER" env-required:"true"`
	Password string `env:"POSTGRES_PASSWORD" env-required:"true"`
	DBName   string `env:"POSTGRES_DB" env-required:"true"`
	SSLMode  string `env:"ssl_mode"`
	MaxConns int32  `yaml:"max_conns"`
	MinConns int32  `yaml:"min_conns"`
}

type ZeroLog struct {
	Level          string   `yaml:"level"`
	SkipFrameCount int      `yaml:"skip_frame_count"`
	InTG           bool     `yaml:"in_tg"`
	ChatID         int64    `yaml:"chat_id"`
	TGToken        string   `env:"TG_TOKEN"`
	AlertUsers     []string `yaml:"alert_users"`
}

func NewConfig(configFile string) (*Config, error) {
	var cfg Config

	err := cleanenv.ReadConfig(path.Join("./", configFile), &cfg)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	err = cleanenv.UpdateEnv(&cfg)
	if err != nil {
		return nil, fmt.Errorf("error updating env: %w", err)
	}

	return &cfg, nil
}

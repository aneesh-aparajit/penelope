package config

import (
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/go-viper/mapstructure/v2"
	"github.com/spf13/viper"
)

type Config struct {
	Entry  string `json:"entry"`
	Limits struct {
		MaxHops     uint8         `json:"maxHops"`
		StepTimeout time.Duration `json:"stepTimeout"`
		IdleTimeout time.Duration `json:"idleTimeout"`
	} `json:"limits"`
	Agents map[string]AgentDef `json:"agents"`
}

var (
	once sync.Once
	cfg  *Config = nil
	err  error   = nil
)

func LoadConfig(path string) (*Config, error) {
	once.Do(func() {
		cfg, err = read(path)
		if err != nil {
			slog.Error("failed to read config", slog.String("path", path), slog.Any("err", err))
		}
	})
	return cfg, err
}

func read(path string) (*Config, error) {
	v := viper.New()
	v.SetConfigFile(path)
	v.SetConfigType("yaml")

	v.SetDefault("limits.stepTimeout", "30m")
	v.SetDefault("limits.idleTimeout", "5m")

	if err := v.ReadInConfig(); err != nil {
		slog.Error("failed to read config", slog.String("path", path), slog.Any("err", err))
		return nil, fmt.Errorf("read %s: %w", path, err)
	}

	cfg := &Config{}
	if err := v.UnmarshalExact(cfg, viper.DecodeHook(mapstructure.ComposeDecodeHookFunc(
		mapstructure.StringToTimeDurationHookFunc(),
		mapstructure.TextUnmarshallerHookFunc(),
	))); err != nil {
		return nil, fmt.Errorf("unmarshal %s: %w", path, err)
	}
	return cfg, nil
}

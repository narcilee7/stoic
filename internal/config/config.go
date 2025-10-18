package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	LLM LLMConfig `mapstructure:"llm"`
	UI  UIConfig  `mapstructure:"ui"`
}

type LLMConfig struct {
	Provider string       `mapstructure:"provider"`
	Ollama   OllamaConfig `mapstructure:"ollama"`
}

type OllamaConfig struct {
	BaseURL string `mapstructure:"base_url"`
	Model   string `mapstructure:"model"`
}

type UIConfig struct {
	Theme    string `mapstructure:"theme"`
	Language string `mapstructure:"language"`
}

func Load() (*Config, error) {
	viper.SetConfigName("stoic")
	viper.SetConfigType("yaml")

	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.stoic")
	viper.AddConfigPath("/etc/stoic")

	setDefaults()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}
	return &config, nil
}

func setDefaults() {
	viper.SetDefault("ai.provider", "ollama")
	viper.SetDefault("ai.ollama.base_url", "http://localhost:11434")
	viper.SetDefault("ai.ollama.model", "llama2")
	viper.SetDefault("ui.theme", "calm")
	viper.SetDefault("ui.language", "zh")
}

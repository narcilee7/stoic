package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	LLM          LLMConfig          `mapstructure:"llm"`
	UI           UIConfig           `mapstructure:"ui"`
	Philosophers PhilosophersConfig `mapstructure:"philosophers"`
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

type PhilosopherConfig struct {
	ID          string `yaml:"id" mapstructure:"id"`
	Name        string `yaml:"name" mapstructure:"name"`
	School      string `yaml:"school" mapstructure:"school"`
	Description string `yaml:"description" mapstructure:"description"`
	Personality string `yaml:"personality" mapstructure:"personality"`
	Style       string `yaml:"style" mapstructure:"style"`
	Prompt      string `yaml:"prompt" mapstructure:"prompt"`
	Emoji       string `yaml:"emoji" mapstructure:"emoji"`
	Enabled     bool   `yaml:"enabled" mapstructure:"enabled"`
}

type PhilosophersConfig struct {
	Custom   []PhilosopherConfig `mapstructure:"custom"`
	Disabled []string            `mapstructure:"disabled"`
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
	viper.SetDefault("llm.provider", "ollama")
	viper.SetDefault("llm.ollama.base_url", "http://localhost:11434")
	viper.SetDefault("llm.ollama.model", "llama2")
	viper.SetDefault("ui.theme", "calm")
	viper.SetDefault("ui.language", "zh")
	viper.SetDefault("philosophers.custom", []interface{}{})
	viper.SetDefault("philosophers.disabled", []interface{}{})
}

func (c *PhilosopherConfig) Validate() error {
	if c.ID == "" {
		return fmt.Errorf("id is required")
	}
	if c.Name == "" {
		return fmt.Errorf("name is required")
	}
	if c.School == "" {
		return fmt.Errorf("school is required")
	}
	if c.Description == "" {
		return fmt.Errorf("description is required")
	}
	if c.Personality == "" {
		return fmt.Errorf("personality is required")
	}
	if c.Style == "" {
		return fmt.Errorf("style is required")
	}
	if c.Prompt == "" {
		return fmt.Errorf("prompt is required")
	}
	if c.Emoji == "" {
		return fmt.Errorf("emoji is required")
	}
	return nil
}

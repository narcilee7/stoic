package config

import (
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
)

func TestLoadConfig(t *testing.T) {
	root, err := filepath.Abs("../..")
	if err != nil {
		t.Fatalf("Failed to get root path: %v", err)
	}

	v := viper.New()
	v.SetConfigName("stoic")
	v.AddConfigPath(root)
	err = v.ReadInConfig()
	if err != nil {
		t.Errorf("Failed to load config: %v", err)
	}

	if v.GetString("llm.provider") != "ollama" {
		t.Error("Expected llm.provider to be 'ollama', but got different value")
	}
}

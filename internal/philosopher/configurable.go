package philosopher

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/stoic/internal/config"
	"github.com/stoic/provider"
)

type ConfigurablePhilosopher struct {
	config      config.PhilosopherConfig
	provider    provider.Provider
	memory      []ConversationMemory
	mood        string
	preferences map[string]interface{}
}

type ConversationMemory struct {
	UserMessage         string    `json:"user_message"`
	PhilosopherResponse string    `json:"philosopher_response"`
	Timestamp           time.Time `json:"timestamp"`
	Emotion             string    `json:"emotion"`
}

func NewConfigurablePhilosopher(config config.PhilosopherConfig, provider provider.Provider) *ConfigurablePhilosopher {
	return &ConfigurablePhilosopher{
		config:   config,
		provider: provider,
	}
}

// Name 返回哲学家名字
func (c *ConfigurablePhilosopher) Name() string {
	return c.config.Name
}

// School 返回哲学流派
func (c *ConfigurablePhilosopher) School() string {
	return c.config.School
}

// Description 返回描述
func (c *ConfigurablePhilosopher) Description() string {
	return c.config.Description
}

// Emoji 返回表情符号
func (c *ConfigurablePhilosopher) Emoji() string {
	return c.config.Emoji
}

// Respond 生成哲学回应
func (c *ConfigurablePhilosopher) Respond(ctx context.Context, message string) (string, error) {
	prompt := fmt.Sprintf(`%s

User's message: %s

Respond as %s would, with %s wisdom:`,
		c.config.Prompt,
		message,
		c.config.Name,
		c.config.School)

	response, err := c.provider.Generate(ctx, prompt)
	if err != nil {
		return "", fmt.Errorf("error generating response: %w", err)
	}

	return response, nil
}

func (c *ConfigurablePhilosopher) AddMemory(userMsg, response, emotion string) {
	c.memory = append(c.memory, ConversationMemory{
		UserMessage:         userMsg,
		PhilosopherResponse: response,
		Timestamp:           time.Now(),
		Emotion:             emotion,
	})

	// 进程内存缓存10条
	if len(c.memory) > 10 {
		c.memory = c.memory[1:]
	}
}

func (c *ConfigurablePhilosopher) GetRelevantMemories(query string) []ConversationMemory {
	var relevant []ConversationMemory
	for _, memory := range c.memory {
		if strings.Contains(strings.ToLower(memory.UserMessage), strings.ToLower(query)) {
			relevant = append(relevant, memory)
		}
	}
	return relevant
}

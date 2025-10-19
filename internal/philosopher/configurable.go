package philosopher

import (
	"context"
	"fmt"

	"github.com/narcilee7/stoic/provider"
)

// ConfigurablePhilosopher 可配置的哲学家
type ConfigurablePhilosopher struct {
	config   PhilosopherConfig
	provider provider.Provider
}

// NewConfigurablePhilosopher 创建可配置哲学家
func NewConfigurablePhilosopher(config PhilosopherConfig, provider provider.Provider) *ConfigurablePhilosopher {
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

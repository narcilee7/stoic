package philosopher

import "fmt"

type PhilosopherConfig struct {
	ID          string `yaml:"id"`          // 唯一标识符
	Name        string `yaml:"name"`        // 显示名称
	School      string `yaml:"school"`      // 哲学流派
	Description string `yaml:"description"` // 描述
	Personality string `yaml:"personality"` // 性格特征
	Style       string `yaml:"style"`       // 语言风格
	Prompt      string `yaml:"prompt"`      // AI提示词模板
	Emoji       string `yaml:"emoji"`       // 表情符号
	Enabled     bool   `yaml:"enabled"`     // 是否启用
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

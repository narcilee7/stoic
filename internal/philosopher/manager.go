package philosopher

import (
	"context"
	"fmt"

	"github.com/narcilee7/stoic/internal/config"
	"github.com/narcilee7/stoic/provider"
)

type Philosopher interface {
	Name() string
	School() string
	Respond(ctx context.Context, message string) (string, error)
}

type Manager struct {
	config       *config.Config
	provider     provider.Provider
	philosophers map[string]Philosopher
	current      Philosopher
}

func NewManager(cfg *config.Config) (*Manager, error) {
	llmProvider, err := createProvider(cfg)
	if err != nil {
		return nil, fmt.Errorf("error creating AI provider: %w", err)
	}

	manager := &Manager{
		config:       cfg,
		provider:     llmProvider,
		philosophers: make(map[string]Philosopher),
	}

	// 注册默认哲学家
	if err := manager.registerDefaultPhilosophers(); err != nil {
		return nil, fmt.Errorf("error registering default philosophers: %w", err)
	}

	return manager, nil
}

func createProvider(cfg *config.Config) (provider.Provider, error) {
	switch cfg.LLM.Provider {
	case "ollama":
		return provider.NewOllamaProvider(cfg.LLM.Ollama.BaseURL, cfg.LLM.Ollama.Model), nil
	default:
		return nil, fmt.Errorf("unsupported AI provider: %s", cfg.LLM.Provider)
	}
}

func (m *Manager) registerDefaultPhilosophers() error {
	stoic := NewStoicPhilosopher(m.provider)
	m.philosophers["stoic"] = stoic

	taoist := NewTaoistPhilosopher(m.provider)
	m.philosophers["taoist"] = taoist

	m.current = stoic

	return nil
}

func (m *Manager) GetPhilosopher(name string) (Philosopher, error) {
	phil, ok := m.philosophers[name]
	if !ok {
		return nil, fmt.Errorf("philosopher not found: %s", name)
	}
	return phil, nil
}

func (m *Manager) SetCurrentPhilosopher(name string) error {
	phil, err := m.GetPhilosopher(name)
	if err != nil {
		return err
	}
	m.current = phil
	return nil
}

func (m *Manager) GetCurrentPhilosopher() Philosopher {
	return m.current
}

func (m *Manager) ListPhilosophers() []string {
	names := make([]string, 0, len(m.philosophers))
	for name := range m.philosophers {
		names = append(names, name)
	}
	return names
}

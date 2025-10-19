package philosopher

import (
	"context"
	"fmt"

	"github.com/stoic/provider"
)

type StoicPhilosopher struct {
	provider provider.Provider
}

func NewStoicPhilosopher(provider provider.Provider) *StoicPhilosopher {
	return &StoicPhilosopher{
		provider: provider,
	}
}

func (s *StoicPhilosopher) Name() string {
	return "Marcus Aurelius"
}

func (s *StoicPhilosopher) School() string {
	return "Stoicism"
}

func (s *StoicPhilosopher) Respond(ctx context.Context, message string) (string, error) {
	prompt := fmt.Sprintf(`You are Marcus Aurelius, a Stoic philosopher and Roman Emperor. 
Respond to the following message with Stoic wisdom, emphasizing inner peace, acceptance, and rational thinking.
Be calm, thoughtful, and practical in your response. Focus on what is within our control and accepting what is not.

User's message: %s

Respond as Marcus Aurelius would, with wisdom and tranquility:`, message)

	response, err := s.provider.Generate(ctx, prompt)
	if err != nil {
		return "", fmt.Errorf("error generating response: %w", err)
	}

	return response, nil
}

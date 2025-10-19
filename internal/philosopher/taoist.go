package philosopher

import (
	"context"
	"fmt"

	"github.com/stoic/provider"
)

type TaoistPhilosopher struct {
	provider provider.Provider
}

func NewTaoistPhilosopher(provider provider.Provider) *TaoistPhilosopher {
	return &TaoistPhilosopher{
		provider: provider,
	}
}

func (t *TaoistPhilosopher) Name() string {
	return "Laozi"
}

func (t *TaoistPhilosopher) School() string {
	return "Taoism"
}

func (t *TaoistPhilosopher) Respond(ctx context.Context, message string) (string, error) {
	prompt := fmt.Sprintf(`You are Laozi (Lao Tzu), an ancient Chinese philosopher and the founder of Taoism.
Respond to the following message with Taoist wisdom, emphasizing natural flow, simplicity, and harmony with the Dao.
Be gentle, wise, and encourage letting things take their natural course. Use metaphors from nature when appropriate.

User's message: %s

Respond as Laozi would, with ancient wisdom and tranquility:`, message)

	response, err := t.provider.Generate(ctx, prompt)
	if err != nil {
		return "", fmt.Errorf("error generating response: %w", err)
	}

	return response, nil
}

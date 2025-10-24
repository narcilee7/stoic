package provider

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

type GenerateOptions struct {
	Temperature float64
	MaxTokens   int
	TopP        float64
	Stream      bool
}

type Provider interface {
	Generate(ctx context.Context, prompt string) (string, error)
	GenerateWithOptions(ctx context.Context, prompt string, opts GenerateOptions) (string, error)
	StreamGenerate(ctx context.Context, prompt string, opts GenerateOptions) (<-chan string, <-chan error)
}

type BaseProvider struct {
	Client  *http.Client
	Model   string
	BaseURL string
	Timeout time.Duration
}

func NewBaseProvider(baseURL, model string) *BaseProvider {
	return &BaseProvider{
		Client: &http.Client{
			Timeout: 30 * time.Second,
		},
		Model:   model,
		BaseURL: baseURL,
		Timeout: 30 * time.Second,
	}
}

func (b *BaseProvider) PrepareRequest(ctx context.Context, method, endpoint string, body io.Reader) (*http.Request, error) {
	url := b.BaseURL + endpoint
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

func (b *BaseProvider) HandleResponse(resp *http.Response) (string, error) {
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func (b *BaseProvider) Generate(ctx context.Context, prompt string) (string, error) {
	return b.GenerateWithOptions(ctx, prompt, GenerateOptions{})
}

func (b *BaseProvider) GenerateWithOptions(ctx context.Context, prompt string, opts GenerateOptions) (string, error) {
	return "", fmt.Errorf("GenerateWithOptions not implemented")
}

func (b *BaseProvider) StreamGenerate(ctx context.Context, prompt string, opts GenerateOptions) (<-chan string, <-chan error) {
	respChan := make(chan string)
	errChan := make(chan error, 1)
	errChan <- fmt.Errorf("StreamGenerate not implemented")
	close(respChan)
	close(errChan)
	return respChan, errChan
}

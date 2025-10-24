package cloud

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/stoic/provider"
)

type MoonshotProvider struct {
	provider.BaseProvider
	APIKey string
}

type moonshotRequest struct {
	Model    string    `json:"model"`
	Messages []message `json:"messages"`
	Stream   bool      `json:"stream,omitempty"`
}

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type moonshotResponse struct {
	Choices []choice `json:"choices"`
}

type choice struct {
	Message      message `json:"message"`
	Delta        message `json:"delta,omitempty"` // For streaming
	FinishReason string  `json:"finish_reason,omitempty"`
}

func NewMoonshotProvider(apiKey, baseURL, model string) *MoonshotProvider {
	base := provider.NewBaseProvider(baseURL, model)
	return &MoonshotProvider{
		BaseProvider: *base,
		APIKey:       apiKey,
	}
}

func (m *MoonshotProvider) GenerateWithOptions(ctx context.Context, prompt string, opts provider.GenerateOptions) (string, error) {
	if opts.Stream {
		return "", fmt.Errorf("use StreamGenerate for streaming")
	}

	reqBody := moonshotRequest{
		Model: m.Model,
		Messages: []message{
			{Role: "user", Content: prompt},
		},
		Stream: false,
		// Map opts to request fields if needed
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := m.PrepareRequest(ctx, "POST", "/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+m.APIKey)

	resp, err := m.Client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result moonshotResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if len(result.Choices) == 0 {
		return "", fmt.Errorf("no response choices")
	}

	return result.Choices[0].Message.Content, nil
}

func (m *MoonshotProvider) StreamGenerate(ctx context.Context, prompt string, opts provider.GenerateOptions) (<-chan string, <-chan error) {
	respChan := make(chan string)
	errChan := make(chan error, 1)

	go func() {
		defer close(respChan)
		defer close(errChan)

		reqBody := moonshotRequest{
			Model: m.Model,
			Messages: []message{
				{Role: "user", Content: prompt},
			},
			Stream: true,
		}

		jsonData, err := json.Marshal(reqBody)
		if err != nil {
			errChan <- err
			return
		}

		req, err := m.PrepareRequest(ctx, "POST", "/v1/chat/completions", bytes.NewBuffer(jsonData))
		if err != nil {
			errChan <- err
			return
		}
		req.Header.Set("Authorization", "Bearer "+m.APIKey)
		req.Header.Set("Accept", "text/event-stream")

		resp, err := m.Client.Do(req)
		if err != nil {
			errChan <- err
			return
		}
		defer resp.Body.Close()

		reader := bufio.NewReader(resp.Body)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				if err != io.EOF {
					errChan <- err
				}
				return
			}

			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}

			if strings.HasPrefix(line, "data: ") {
				data := line[6:]
				if data == "[DONE]" {
					return
				}

				var chunk moonshotResponse
				if err := json.Unmarshal([]byte(data), &chunk); err != nil {
					errChan <- err
					return
				}

				if len(chunk.Choices) > 0 {
					respChan <- chunk.Choices[0].Delta.Content
				}
			}
		}
	}()

	return respChan, errChan
}

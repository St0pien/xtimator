package xtimators

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/St0pien/ratio-xtimator/internal/core"
	"github.com/St0pien/ratio-xtimator/internal/prompts"
	"github.com/ollama/ollama/api"
)

type OllamaXtimator struct {
	client *api.Client
	model  string
	prompt prompts.Prompt
}

type OllamaXtimatorConfig struct {
	Model  string
	Prompt prompts.Prompt
}

func NewOllamaXtimator(config OllamaXtimatorConfig) (OllamaXtimator, error) {
	client, err := api.ClientFromEnvironment()
	if err != nil {
		return OllamaXtimator{client: nil}, fmt.Errorf("Couldn't instantiate ollama client: %w", err)
	}

	return OllamaXtimator{
		client: client,
		model:  config.Model,
		prompt: config.Prompt,
	}, nil
}

func processImpression(text string) TweetImpression {
	return TweetImpression(strings.ToLower(strings.Trim(text, " \t\n")))
}

func (ollama OllamaXtimator) Xtimate(post core.Post, response core.Post) (TweetImpression, error) {
	promptText, err := ollama.prompt.Build(post, response)
	if err != nil {
		return "", fmt.Errorf("Failed to build prompt for post %v: %w", post, err)
	}

	req := &api.ChatRequest{
		Model: ollama.model,
		Messages: []api.Message{
			{Role: "system", Content: ollama.prompt.GetSystemPrompt()},
			{Role: "user", Content: promptText},
		},
	}

	var llmResponse strings.Builder
	var wg sync.WaitGroup
	wg.Add(1)

	err = ollama.client.Chat(context.Background(), req, func(cr api.ChatResponse) error {
		llmResponse.WriteString(cr.Message.Content)

		if cr.Done {
			wg.Done()
		}

		return nil
	})

	if err != nil {
		return "", fmt.Errorf("Ollama query failed: %w", err)
	}

	return processImpression(llmResponse.String()), nil
}

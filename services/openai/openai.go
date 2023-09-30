package openai

import (
	"bytes"
	"encoding/json"

	"io"
	"net/http"

	"github.com/j-dunham/openai-cli/config"
)

type Service interface {
	GetCompletion(prompt string) (string, error)
}

type service struct {
	cfg *config.Config
}

func (s service) GetCompletion(prompt string) (string, error) {
	data := Data{
		Model:       s.cfg.OpenAiModel,
		MaxTokens:   s.cfg.OpenAiMaxTokens,
		Temperature: s.cfg.OpenAiTemperature,
		Messages:    []Message{{Role: "user", Content: prompt}},
	}

	payload, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	// maybe use functional options for setting target url?
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(payload))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.cfg.OpenAiToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var result Response
	json.Unmarshal(body, &result)
	return result.Choices[0].Message.Content, nil
}

func NewService(cfg *config.Config) Service {
	return service{}
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Data struct {
	Model       string    `json:"model"`
	MaxTokens   int       `json:"max_tokens"`
	Temperature float64   `json:"temperature"`
	Messages    []Message `json:"messages"`
}

type Response struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index        int     `json:"index"`
		Message      Message `json:"message"`
		FinishReason string  `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

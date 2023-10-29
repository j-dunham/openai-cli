package openai

import (
	"bytes"
	"encoding/json"
	"log"

	"io"
	"net/http"
	"net/url"

	"github.com/j-dunham/openai-cli/config"
)

type Service interface {
	GetCompletion(prompt []Message) (string, error)
}

type service struct {
	cfg *config.Config
}

type RequestOption func(*http.Request)

func WithURL(u string) RequestOption {
	return func(req *http.Request) {
		req.URL, _ = url.Parse(u)
	}
}

func WithMethod(m string) RequestOption {
	return func(req *http.Request) {
		req.Method = m
	}
}

func WithHeader(key, value string) RequestOption {
	return func(req *http.Request) {
		if req.Header == nil {
			req.Header = make(http.Header)
		}
		req.Header.Add(key, value)
	}
}

func WithBody(body []byte) RequestOption {
	return func(req *http.Request) {
		req.Body = io.NopCloser(bytes.NewReader(body))
	}
}

func NewRequest(opts ...RequestOption) *http.Request {
	req := &http.Request{}
	for _, opt := range opts {
		opt(req)
	}
	return req
}

func (s service) GetCompletion(messages []Message) (string, error) {
	if s.cfg == nil {
		log.Fatal("config is nil")
	}

	data := Data{
		Model:       s.cfg.OpenAiModel,
		MaxTokens:   s.cfg.OpenAiMaxTokens,
		Temperature: s.cfg.OpenAiTemperature,
		Messages:    messages,
	}

	payload, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	req := NewRequest(
		WithURL("https://api.openai.com/v1/chat/completions"),
		WithMethod("POST"),
		WithHeader("Content-Type", "application/json"),
		WithHeader("Authorization", "Bearer "+s.cfg.OpenAiToken),
		WithBody(payload),
	)

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
	return service{
		cfg: cfg,
	}
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

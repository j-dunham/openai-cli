package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	OpenAiToken       string
	OpenAiModel       string
	OpenAiMaxTokens   int
	OpenAiTemperature float64
	DBFile            string
}

func NewConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		fmt.Println("Loading config from ", homeDir)
		envPath := fmt.Sprintf("%s/.openai-cli", homeDir)
		if err := godotenv.Load(envPath); err != nil {
			return nil, err
		}
	}

	if os.Getenv("OPEN_AI_TOKEN") == "" {
		return nil, fmt.Errorf("OPEN_AI_TOKEN is not set")
	}

	if os.Getenv("OPEN_AI_MODEL") == "" {
		return nil, fmt.Errorf("OPEN_AI_MODEL is not set")
	}

	if os.Getenv("OPEN_AI_MAX_TOKENS") == "" {
		return nil, fmt.Errorf("OPEN_AI_MAX_TOKENS is not set")
	}

	if os.Getenv("OPEN_AI_TEMPERATURE") == "" {
		return nil, fmt.Errorf("OPEN_AI_TEMPERATURE is not set")
	}

	maxTokens, err := strconv.Atoi(os.Getenv("OPEN_AI_MAX_TOKENS"))
	if err != nil {
		return nil, err
	}

	temperature, err := strconv.ParseFloat(os.Getenv("OPEN_AI_TEMPERATURE"), 64)
	if err != nil {
		return nil, err
	}

	return &Config{
		OpenAiToken:       os.Getenv("OPEN_AI_TOKEN"),
		OpenAiModel:       os.Getenv("OPEN_AI_MODEL"),
		OpenAiMaxTokens:   maxTokens,
		OpenAiTemperature: temperature,
		DBFile:            "Prompt.db",
	}, nil
}

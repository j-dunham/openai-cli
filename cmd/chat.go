package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/fatih/color"
)

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
	ID        string `json:"id"`
	Object    string `json:"object"`
	Created   int    `json:"created"`
	Model     string `json:"model"`
	Choices   []struct {
			Index        int `json:"index"`
			Message      struct {
					Role    string `json:"role"`
					Content string `json:"content"`
			} `json:"message"`
			FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
			PromptTokens     int `json:"prompt_tokens"`
			CompletionTokens int `json:"completion_tokens"`
			TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

func createRequestData(prompt string) ([]byte, error) {
	data := Data{
		Model:       "gpt-3.5-turbo",
		MaxTokens:   100,
		Temperature: 0.9,
		Messages:    []Message{{Role: "user", Content: prompt}},
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return jsonData, nil
}

func createRequest(prompt string) (*http.Request, error) {

	jsonData, _ := createRequestData(prompt)
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	token := os.Getenv("OPEN_AI_TOKEN")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	return req, nil
}

func printResponse(prompt string, response http.Response){
	body, _ := io.ReadAll(response.Body)
	var resp Response
	json.Unmarshal(body, &resp)
	blue := color.New(color.FgBlue).PrintlnFunc()
	green := color.New(color.FgGreen).PrintlnFunc()

	blue("Prompt:")
	green(prompt, "\n")
	blue("Response:")
	green(resp.Choices[0].Message.Content)
}

func Execute(prompt string) {
	client := &http.Client{}
	req, _ := createRequest(prompt)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()
	printResponse(prompt, *resp)
}

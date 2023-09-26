package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/fatih/color"
	"github.com/j-dunham/openai-cli/util"
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
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
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
	maxTokens, err := strconv.Atoi(os.Getenv("OPEN_AI_MAX_TOKENS"))
	if err != nil {
		return nil, err
	}

	temperatureStr := os.Getenv("OPEN_AI_TEMPERATURE")
	temperature, err := strconv.ParseFloat(temperatureStr, 64)
	if err != nil {
		return nil, err
	}

	data := Data{
		Model:       os.Getenv("OPEN_AI_MODEL"),
		MaxTokens:   maxTokens,
		Temperature: temperature,
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

func parseResponse(res http.Response) Response {
	body, _ := io.ReadAll(res.Body)
	var resp Response
	json.Unmarshal(body, &resp)
	return resp
}

func printResponse(prompt string, response Response) {

	blue := color.New(color.FgBlue).PrintlnFunc()
	green := color.New(color.FgGreen).PrintlnFunc()

	fmt.Println("====================================")
	blue("| Prompt:")
	green("| ", prompt)
	blue("| Response:")
	green(util.WrapText(response.Choices[0].Message.Content, 100, "|  "))
	fmt.Println("====================================")
}

func Execute(prompt string, save bool) {
	client := &http.Client{}
	req, _ := createRequest(prompt)

	done := make(chan bool)
	go util.LoadingAnimation("thinking", done)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()
	response := parseResponse(*resp)
	// Stop the loading animation
	done <- true
	if save {
		util.CreateTable()
		util.InsertPrompt(prompt, response.Choices[0].Message.Content)
	}
	printResponse(prompt, response)
}

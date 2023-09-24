package main

import (
	"fmt"
	"flag"

	"github.com/j-dunham/openai-cli/cmd"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file:", err)
		return
	}
	prompt := flag.String("prompt", "", "The prompt to ask ChatGPT.")
    flag.Parse()
	if *prompt == "" {
		fmt.Println("You must provide a prompt to ask ChatGPT.")
		return
	}

	cmd.Execute(*prompt)
}

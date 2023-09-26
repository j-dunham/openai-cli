package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/j-dunham/openai-cli/cmd"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file:", err)
		os.Exit(1)
	}
	prompt := flag.String("prompt", "", "The prompt to ask ChatGPT.")
	flag.Parse()
	if *prompt == "" {
		fmt.Println("You must provide a prompt to ask ChatGPT.")
		os.Exit(2)
	}

	cmd.Execute(*prompt)
}

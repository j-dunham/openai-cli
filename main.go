package main

import (
	"fmt"
	"os"

	"github.com/j-dunham/openai-cli/cmd"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file:", err)
		return
	}

	cmd.Execute(os.Args[1])
}

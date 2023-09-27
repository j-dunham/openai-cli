# OpenAI CLI

This is a command-line interface (CLI) for OpenAI's ChatGPT API. It allows you to interact with the API from your terminal, without needing to use a web browser.

## Setup
1) Create a OpenAI API token at https://platform.openai.com/account/api-keys
2) Create a `.env` file using the `.env.template` file and set `OPEN_AI_TOKEN` to the generated token for step 1

## Usage

```console
go run . -prompt "<PROMPT HERE>" 
```

Example output
```console
You: Tell me a joke about developers?      
OpenAI: Why do developers prefer dark mode?
                                           
Because light attracts bugs!               
                                           
                                           
                                                                              

┃ What is your Prompt?                         
┃                                              


```

# OpenAI CLI

This is a command-line interface (CLI) for OpenAI's ChatGPT API. It allows you to interact with the API from your terminal, without needing to use a web browser.

## Setup
1) Create a OpenAI API token at https://platform.openai.com/account/api-keys
2) Create a `.env` file using the `.env.template` file and set `OPEN_AI_TOKEN` to the generated token for step 1

## Usage

```console
go run . "<PROMPT HERE>" 
```

Example output
```console
====================================
| Prompt:
|  Tell me a joke about a software developer?
| Response:
|  Why don't software developers like nature? Because there are too many bugs! 
====================================
```

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
|  What is ChatGPT?
| Response:
|  ChatGPT is a language model developed by OpenAI. It is designed to generate human-like text 
|  responses given a prompt or a conversation. It is trained using Reinforcement Learning from Human 
|  Feedback (RLHF), combining both supervised fine-tuning and reinforcement learning. ChatGPT is 
|  capable of engaging in conversations on a wide range of topics, providing detailed responses based 
|  on the input it receives. It has been trained with a diverse dataset from the internet, but it may 
|  sometimes produce incorrect or nonsensical 
====================================
```

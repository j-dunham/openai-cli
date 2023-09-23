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
Prompt:
What is ChatGPT? 

Response:
ChatGPT is an advanced language model developed by OpenAI. It uses deep learning techniques to generate human-like responses in natural language conversations. ChatGPT is trained using a method called Reinforcement Learning from Human Feedback (RLHF), where human AI trainers provide responses and rank different model-generated alternatives. This training process helps ChatGPT improve its responses over time. It can be used for a variety of tasks like drafting content, answering questions, creating conversational agents, and more.
```

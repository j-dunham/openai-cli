# OpenAI CLI

This is a command-line interface (CLI) for OpenAI's ChatGPT API. It allows you to interact with the API from your terminal, without needing to use a web browser.

### Chat with OpenAI 
![Demo](demo.gif)
### List Previous Prompts
![List Chat](list.gif)

## Setup
1) Create a OpenAI API token at https://platform.openai.com/account/api-keys
2) Create a `.env` file using the `.env.template` file and set `OPEN_AI_TOKEN` to the generated token for step 1

## Special Prompt Commands
- Pefixing the prompt with `/system` will make the prompt a system prompt, which is used as a way to set the behavior of the AI.

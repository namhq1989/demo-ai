package openai

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	oai "github.com/sashabaranov/go-openai"
)

type GenerateTermResult struct {
	Prompt string `json:"prompt"`
}

const generatePrompt = `
	Generate a prompt for text-to-image generation, with the following information:
	- description: {{description}}
	- style: {{style}}
	- colorScheme: {{colorScheme}}
	- text: {{text}}
	- textStyle: {{textStyle}}
	- layout: {{layout}}
	- theme: {{theme}}
	- additionalElements: {{additionalElements}}

	If a value is empty or is "-", just skip it.
	If the description is related to human-render, add an additional instruction to prevent abnormalities render.
	Must in 4k resolution.
	Respond a JSON only with the following format:
	- "prompt": "<generated prompt>"
`

func (o OpenAI) GeneratePrompt(description, style, colorScheme, text, textStyle, layout, theme, additionalElements string) string {
	prompt := strings.ReplaceAll(generatePrompt, "{{description}}", description)
	prompt = strings.ReplaceAll(prompt, "{{style}}", style)
	prompt = strings.ReplaceAll(prompt, "{{colorScheme}}", colorScheme)
	prompt = strings.ReplaceAll(prompt, "{{text}}", text)
	prompt = strings.ReplaceAll(prompt, "{{textStyle}}", textStyle)
	prompt = strings.ReplaceAll(prompt, "{{layout}}", layout)
	prompt = strings.ReplaceAll(prompt, "{{theme}}", theme)
	prompt = strings.ReplaceAll(prompt, "{{additionalElements}}", additionalElements)

	resp, err := o.client.CreateChatCompletion(context.Background(), oai.ChatCompletionRequest{
		Model:       oai.GPT3Dot5Turbo,
		Messages:    []oai.ChatCompletionMessage{{Role: oai.ChatMessageRoleUser, Content: prompt}},
		MaxTokens:   150,
		Temperature: 0.5,
	})

	if err != nil {
		return ""
	}

	var result GenerateTermResult
	if err = json.Unmarshal([]byte(resp.Choices[0].Message.Content), &result); err != nil {
		fmt.Println("data", resp.Choices[0].Message.Content)
		return ""
	}

	return result.Prompt
}

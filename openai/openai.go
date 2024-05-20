package openai

import (
	"fmt"

	oai "github.com/sashabaranov/go-openai"
)

type OpenAI struct {
	client *oai.Client
}

func NewOpenAIClient(token string) *OpenAI {
	fmt.Printf("⚡️ [openai]: connected \n")

	return &OpenAI{client: oai.NewClient(token)}
}

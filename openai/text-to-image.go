package openai

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/namhq1989/demo-ai/util"
	oai "github.com/sashabaranov/go-openai"
)

type TextToImagePayload struct {
	Prompt         string `json:"prompt"`
	Model          string `json:"model"`
	NumOfImages    int    `json:"n"`
	ResponseFormat string `json:"response_format"`
	Size           string `json:"size"`
	Style          string `json:"style"`
}

func (o OpenAI) TextToImage(payload TextToImagePayload) (string, error) {
	resp, err := o.client.CreateImage(context.Background(), oai.ImageRequest{
		Prompt:         payload.Prompt,
		Model:          payload.Model,
		N:              payload.NumOfImages,
		Size:           payload.Size,
		Style:          payload.Style,
		ResponseFormat: payload.ResponseFormat,
	})

	if err != nil {
		return "", errors.New("cannot call openai api")
	}

	fileNames := make([]string, 0)

	for _, item := range resp.Data {
		// random seed
		randSource := rand.New(rand.NewSource(time.Now().Unix()))
		seed := randSource.Intn(4294967294)

		fileName := fmt.Sprintf("%d-%d.jpeg", seed, time.Now().Unix())
		filePath := fmt.Sprintf("generated/%s", fileName)
		err = decodeAndSaveImage(item.B64JSON, filePath)
		if err != nil {
			return "", fmt.Errorf("failed to decode and save image: %v", err)
		}

		fileNames = append(fileNames, fileName)
	}

	return util.GetImageURL(fileNames[0]), nil
}

func decodeAndSaveImage(base64Image, filePath string) error {
	rawDecodedText, _ := base64.StdEncoding.DecodeString(base64Image)
	err := os.WriteFile(filePath, rawDecodedText, 0644) // Use appropriate file permissions
	if err != nil {
		return err
	}

	return nil
}

var mapProductSize = map[string]string{
	"t-shirt":    "1792x1024",
	"tumbler":    "1024x1024",
	"phone-case": "1024x1792",
	"hoodie":     "1792x1024",
	"mug":        "1024x1024",
	"tote-bag":   "1792x1024",
	"pillow":     "1024x1024",
	"poster":     "1792x1024",
	"notebook":   "1792x1024",
	"sticker":    "1024x1024",
}

func GetSize(product string) string {
	size, exists := mapProductSize[product]
	if !exists {
		return "1024x1024"
	}
	return size
}

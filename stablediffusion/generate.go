package stablediffusion

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"time"

	"github.com/namhq1989/demo-ai/util"
)

type TextToImagePayload struct {
	Prompt       string `json:"prompt" form:"prompt"`
	Mode         string `json:"mode" form:"mode"`
	Model        string `json:"model" form:"model"`
	AspectRatio  string `json:"aspect_ratio" form:"aspect_ratio"`
	Seed         int    `json:"seed" form:"seed"`
	OutputFormat string `json:"output_format" form:"output_format"`
}

type ImageToImagePayload struct {
	Prompt       string  `json:"prompt" form:"prompt"`
	Mode         string  `json:"mode" form:"mode"`
	Model        string  `json:"model" form:"model"`
	Image        string  `json:"image" form:"image"`
	Seed         int     `json:"seed" form:"seed"`
	OutputFormat string  `json:"output_format" form:"output_format"`
	Strength     float64 `json:"strength" form:"strength"`
}

type GenerateResponse struct {
	Image string `json:"image"`
}

type apiGenerateImageResponse struct {
	Image        string `json:"image"`
	FinishReason string `json:"finish_reason"`
	Seed         int    `json:"seed"`
}

func (sd StableDiffusion) TextToImage(payload TextToImagePayload) (string, error) {
	if payload.Prompt == "" || payload.AspectRatio == "" {
		return "", errors.New("invalid payload")
	}

	// random seed
	randSource := rand.New(rand.NewSource(time.Now().Unix()))
	payload.Seed = randSource.Intn(4294967294)

	if payload.OutputFormat == "" {
		payload.OutputFormat = "jpeg"
	}

	payload.Mode = "text-to-image"

	// Convert struct to form data
	b, contentType, err := structToFormData(payload)
	if err != nil {
		return "", fmt.Errorf("failed to convert struct to form data: %v", err)
	}

	// request
	req, err := http.NewRequest("POST", "https://api.stability.ai/v2beta/stable-image/generate/sd3", b)
	if err != nil {
		return "", err
	}

	// set headers
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("authorization", sd.apiKey)
	req.Header.Set("accept", "application/json")

	// execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to execute request: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	// read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("non-200 status code: %d, body: %s", resp.StatusCode, string(body))
	}

	// parse the response
	var responsePayload apiGenerateImageResponse
	err = json.Unmarshal(body, &responsePayload)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %v", err)
	}

	fileName := fmt.Sprintf("%d-%d.jpeg", payload.Seed, time.Now().Unix())
	filePath := fmt.Sprintf("generated/%s", fileName)
	err = decodeAndSaveImage(responsePayload.Image, filePath)
	if err != nil {
		return "", fmt.Errorf("failed to decode and save image: %v", err)
	}

	return util.GetImageURL(fileName), nil
}

func (sd StableDiffusion) ImageToImage(payload ImageToImagePayload) (*GenerateResponse, error) {
	if payload.Prompt == "" {
		return nil, errors.New("invalid payload")
	}

	// random seed
	randSource := rand.New(rand.NewSource(time.Now().Unix()))
	payload.Seed = randSource.Intn(4294967294)

	if payload.OutputFormat == "" {
		payload.OutputFormat = "jpeg"
	}

	payload.Mode = "image-to-image"
	payload.Strength = 0.5
	payload.Image = "templates/chibi.jpg"

	// Convert struct to form data
	b, contentType, err := structToFormData(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to convert struct to form data: %v", err)
	}

	// request
	req, err := http.NewRequest("POST", "https://api.stability.ai/v2beta/stable-image/generate/sd3", b)
	if err != nil {
		return nil, err
	}

	// set headers
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("authorization", sd.apiKey)
	req.Header.Set("accept", "application/json")

	// execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	// read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non-200 status code: %d, body: %s", resp.StatusCode, string(body))
	}

	// parse the response
	var responsePayload apiGenerateImageResponse
	err = json.Unmarshal(body, &responsePayload)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	filePath := fmt.Sprintf("%d-%d.jpeg", payload.Seed, time.Now().Unix())
	err = decodeAndSaveImage(responsePayload.Image, filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to decode and save image: %v", err)
	}

	return &GenerateResponse{
		Image: filePath,
	}, nil
}

func decodeAndSaveImage(base64Image, filePath string) error {
	rawDecodedText, _ := base64.StdEncoding.DecodeString(base64Image)
	err := os.WriteFile(filePath, rawDecodedText, 0644) // Use appropriate file permissions
	if err != nil {
		return err
	}

	return nil
}

func structToFormData[T any](payload T) (*bytes.Buffer, string, error) {
	var b bytes.Buffer
	writer := multipart.NewWriter(&b)

	v := reflect.ValueOf(payload)
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		tag := t.Field(i).Tag.Get("form")
		if tag == "" {
			tag = t.Field(i).Name
		}

		if tag == "image" {
			file, err := os.Open(field.String())
			if err != nil {
				return nil, "", fmt.Errorf("failed to open file: %v", err)
			}
			defer func() { _ = file.Close() }()

			// Create a form file field
			part, err := writer.CreateFormFile("image", filepath.Base(field.String()))
			if err != nil {
				return nil, "", fmt.Errorf("failed to create form file: %v", err)
			}

			// Copy the file data into the form file field
			_, err = io.Copy(part, file)
			if err != nil {
				return nil, "", fmt.Errorf("failed to copy file data: %v", err)
			}
		} else {
			var value string
			switch field.Kind() {
			case reflect.String:
				value = field.String()
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				value = strconv.FormatInt(field.Int(), 10)
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				value = strconv.FormatUint(field.Uint(), 10)
			case reflect.Float32, reflect.Float64:
				value = strconv.FormatFloat(field.Float(), 'f', -1, 64)
			case reflect.Bool:
				value = strconv.FormatBool(field.Bool())
			default:
				return nil, "", fmt.Errorf("unsupported field type: %s", field.Kind())
			}

			err := writer.WriteField(tag, value)
			if err != nil {
				return nil, "", fmt.Errorf("failed to write field %s: %v", tag, err)
			}
		}
	}

	err := writer.Close()
	if err != nil {
		return nil, "", fmt.Errorf("failed to close writer: %v", err)
	}

	return &b, writer.FormDataContentType(), nil
}

func (sd StableDiffusion) GetAspectRatioFromProduct(product string) string {
	if product == "t-shirt" {
		return "4:5"
	}
	if product == "tumbler" {
		return "1:1"
	}
	if product == "phone-case" {
		return "9:16"
	}
	if product == "hoodie" {
		return "4:5"
	}
	if product == "mug" {
		return "1:1"
	}
	if product == "tote-bag" {
		return "4:5"
	}
	if product == "pillow" {
		return "1:1"
	}
	if product == "poster" {
		return "2:3"
	}
	if product == "notebook" {
		return "3:2"
	}
	if product == "sticker" {
		return "1:1"
	}

	return "1:1"
}

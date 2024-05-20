package stablediffusion

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"mime/multipart"
	"net/http"
	"time"
)

func (sd StableDiffusion) EditImage(imgBase64, prompt string) (string, error) {
	requestBody, contentType, err := sd.generateEditRequestData(imgBase64, prompt, "jpeg")
	if err != nil {
		fmt.Println("Error generating request data:", err)
		return "", err
	}

	name, err := sd.callEditImageAPIAndProcessResponse(requestBody, contentType)
	if err != nil {
		fmt.Println("Error calling API and processing response:", err)
	}

	return name, err
}

func (sd StableDiffusion) generateEditRequestData(imageBase64 string, prompt string, outputFormat string) (*bytes.Buffer, string, error) {
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	// Decode the base64 image data
	imageData, err := base64.StdEncoding.DecodeString(imageBase64)
	if err != nil {
		return nil, "", fmt.Errorf("error decoding base64 image data: %v", err)
	}

	// Add the image file to the form
	fileName := fmt.Sprintf("%d", time.Now().Unix())
	imagePart, err := writer.CreateFormFile("image", fileName)
	if err != nil {
		return nil, "", fmt.Errorf("error creating form file for image: %v", err)
	}
	_, err = io.Copy(imagePart, bytes.NewReader(imageData))
	if err != nil {
		return nil, "", fmt.Errorf("error copying image data: %v", err)
	}

	// Add the other fields to the form
	err = writer.WriteField("prompt", prompt)
	if err != nil {
		return nil, "", fmt.Errorf("error writing prompt field: %v", err)
	}
	err = writer.WriteField("output_format", outputFormat)
	if err != nil {
		return nil, "", fmt.Errorf("error writing output_format field: %v", err)
	}

	// Close the multipart writer to set the terminating boundary
	err = writer.Close()
	if err != nil {
		return nil, "", fmt.Errorf("error closing writer: %v", err)
	}

	return &requestBody, writer.FormDataContentType(), nil
}

func (sd StableDiffusion) callEditImageAPIAndProcessResponse(requestBody *bytes.Buffer, contentType string) (string, error) {
	req, err := http.NewRequest("POST", "https://api.stability.ai/v2beta/stable-image/edit/inpaint", requestBody)
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}

	// Set the headers
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Authorization", "Bearer "+sd.apiKey)
	req.Header.Set("accept", "application/json")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

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

	// random seed
	randSource := rand.New(rand.NewSource(time.Now().Unix()))
	seed := randSource.Intn(4294967294)
	fileName := fmt.Sprintf("%d-%d.jpeg", seed, time.Now().Unix())
	filePath := fmt.Sprintf("generated/%s", fileName)
	err = decodeAndSaveImage(responsePayload.Image, filePath)
	if err != nil {
		return "", fmt.Errorf("failed to decode and save image: %v", err)
	}

	return fmt.Sprintf("http://localhost:5000/img/%s", fileName), nil
}

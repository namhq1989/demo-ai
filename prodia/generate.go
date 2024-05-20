package prodia

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/namhq1989/demo-ai/util"
)

type TextToImagePayload struct {
	Model    string `json:"model"`
	Prompt   string `json:"prompt"`
	Steps    int    `json:"steps"`
	CFGScale int    `json:"cfg_scale"`
	Sampler  string `json:"sampler"`
	Seed     int    `json:"seed"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
}

type apiTextToImageResponse struct {
	Job    string `json:"job"`
	Status string `json:"status"`
}

func (p Prodia) TextToImage(payload TextToImagePayload) (string, error) {
	// random seed
	randSource := rand.New(rand.NewSource(time.Now().Unix()))
	payload.Seed = randSource.Intn(4294967294)

	url := "https://api.prodia.com/v1/sdxl/generate"

	b, _ := json.Marshal(payload)
	params := strings.NewReader(string(b))

	req, _ := http.NewRequest("POST", url, params)

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("X-Prodia-Key", p.apiKey)

	res, _ := http.DefaultClient.Do(req)

	defer func() { _ = res.Body.Close() }()
	body, _ := io.ReadAll(res.Body)

	// map the body into apiTextToImageResponse
	var response apiTextToImageResponse
	_ = json.Unmarshal(body, &response)

	if response.Status != "queued" {
		return "", fmt.Errorf("failed to generate Prodia image: %s", response.Status)
	}

	image, err := p.fetchJobData(response.Job, 0)
	if err != nil {
		return "", err
	}

	return util.GetImageURL(image), nil
}

type apiFetchJobDataResponse struct {
	Job      string `json:"job"`
	Status   string `json:"status"`
	ImageUrl string `json:"imageUrl"`
}

func (p Prodia) fetchJobData(jobID string, attempts int) (string, error) {
	if attempts > 10 {
		return "", errors.New("cannot fetch Prodia job data")
	}

	url := fmt.Sprintf("https://api.prodia.com/v1/job/%s", jobID)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("accept", "application/json")
	req.Header.Add("X-Prodia-Key", p.apiKey)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	// map the body into apiTextToImageResponse
	var response apiFetchJobDataResponse
	_ = json.Unmarshal(body, &response)

	if response.Status != "succeeded" {
		// util.PrettyPrint(response)

		// fmt.Printf("job %s still running: %s \n", jobID, response.Status)
		time.Sleep(5 * time.Second)
		return p.fetchJobData(jobID, attempts+1)
	}

	name, err := p.downloadImage(response.ImageUrl)
	if err != nil {
		return "", err
	}

	return name, nil
}

func (p Prodia) downloadImage(url string) (string, error) {
	// random seed
	randSource := rand.New(rand.NewSource(time.Now().Unix()))
	seed := randSource.Intn(4294967294)

	fileName := fmt.Sprintf("%d-%d.jpeg", seed, time.Now().Unix())
	filePath := fmt.Sprintf("generated/%s", fileName)

	// Create the output file
	out, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	// Download the image
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Copy the response body to the file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", err
	}

	return fileName, nil
}

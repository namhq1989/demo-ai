package prodia

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

type EditImagePayload struct {
	MaskBlur            int    `json:"mask_blur"`
	InpaintingFullRes   bool   `json:"inpainting_full_res"`
	InpaitingFill       int    `json:"inpainting_fill"`
	InpantingMaskInvert int    `json:"inpainting_mask_invert"`
	ImageData           string `json:"imageData"`
	Model               string `json:"model"`
	Prompt              string `json:"prompt"`
	Steps               int    `json:"steps"`
	CFGScale            int    `json:"cfg_scale"`
	Sampler             string `json:"sampler"`
	Seed                int    `json:"seed"`
}

func (p Prodia) EditImage(payload EditImagePayload) (string, error) {
	// random seed
	randSource := rand.New(rand.NewSource(time.Now().Unix()))
	payload.Seed = randSource.Intn(4294967294)

	// url := "https://api.prodia.com/v1/sdxl/inpainting"
	url := "https://api.prodia.com/v1/sdxl/transform"

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
	if err := json.Unmarshal(body, &response); err != nil {
		fmt.Println("Error unmarshalling Prodia response:", err)
		fmt.Println("body", string(body))
		return "", err
	}

	if response.Status != "queued" {
		return "", fmt.Errorf("failed to generate Prodia image: %s", response.Status)
	}

	image, err := p.fetchJobData(response.Job, 0)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("http://localhost:5000/img/%s", image), nil
}

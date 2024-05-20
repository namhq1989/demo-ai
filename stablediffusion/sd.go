package stablediffusion

type StableDiffusion struct {
	apiKey string
}

func NewStableDiffusion(apiKey string) StableDiffusion {
	return StableDiffusion{
		apiKey: apiKey,
	}
}

const (
	ModelSD3      = "sd3"
	ModelSD3Turbo = "sd3-turbo"
)

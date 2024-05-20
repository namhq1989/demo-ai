package prodia

type Prodia struct {
	apiKey string
}

func NewProdia(apiKey string) Prodia {
	return Prodia{
		apiKey: apiKey,
	}
}

package main

import (
	"errors"
	"os"
	"strconv"
)

type (
	Server struct {
		OpenAIToken           string
		StableDiffusionAPIKey string
		ProdiaAPIKey          string

		// MongoDB
		MongoURL    string
		MongoDBName string
	}
)

func initConfig() Server {
	cfg := Server{
		MongoURL:    getEnvStr("MONGO_URL"),
		MongoDBName: getEnvStr("MONGO_DB_NAME"),

		OpenAIToken:           getEnvStr("OPENAI_TOKEN"),
		StableDiffusionAPIKey: getEnvStr("STABLE_DIFFUSION_API_KEY"),
		ProdiaAPIKey:          getEnvStr("PRODIA_API_KEY"),
	}

	// validation
	if cfg.OpenAIToken == "" {
		panic(errors.New("missing OPENAI_TOKEN"))
	}

	if cfg.StableDiffusionAPIKey == "" {
		panic(errors.New("missing STABLE_DIFFUSION_API_KEY"))
	}

	if cfg.ProdiaAPIKey == "" {
		panic(errors.New("missing ProdiaAPIKey"))
	}

	return cfg
}

func getEnvStr(key string) string {
	v := os.Getenv(key)
	return v
}

// "true" and "false"
func getEnvBool(key string) bool {
	v := os.Getenv(key)
	return v == "true"
}

func getEnvInt(key string) int {
	s := getEnvStr(key)
	v, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return v
}

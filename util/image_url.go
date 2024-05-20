package util

import (
	"fmt"
	"os"
)

func GetImageURL(name string) string {
	return fmt.Sprintf("%s/img/%s", os.Getenv("API_HOST"), name)
}

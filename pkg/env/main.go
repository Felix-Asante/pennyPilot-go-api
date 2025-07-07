package env

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func GetEnv(key string) string {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error while loading env %v", err)
	}
	return os.Getenv(key)
}

func SetEnv(key string, value string) {
	os.Setenv(key, value)
}

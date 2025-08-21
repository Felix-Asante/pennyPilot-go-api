package utils

import "github.com/Felix-Asante/pennyPilot-go-api/pkg/env"

func GetFrontendUrl() string {
	return env.GetEnv("FRONTEND_URL")
}

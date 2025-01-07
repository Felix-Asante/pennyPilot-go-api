package authServices

import (
	"fmt"

	"github.com/felix-Asante/pennyPilot-go-api/src/api/repositories"
)

type AuthServices struct {
	authRepository *repositories.AuthRepository
}

func NewAuthServices(userRepository *repositories.AuthRepository) *AuthServices {
	return &AuthServices{userRepository}
}

func (s *AuthServices) Login(email string, password string) {
	fmt.Println("Login")
}

func (s *AuthServices) Register(email string, password string) {
	fmt.Println("Register")
}

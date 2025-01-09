package authServices

import (
	"errors"
	"fmt"

	"github.com/felix-Asante/pennyPilot-go-api/src/api/repositories"
	customErrors "github.com/felix-Asante/pennyPilot-go-api/src/utils/errors"
)

type AuthServices struct {
	usersRepository *repositories.UsersRepository
}

func NewAuthServices(userRepository *repositories.UsersRepository) *AuthServices {
	return &AuthServices{userRepository}
}

func (s *AuthServices) Login(email string, password string) {
	fmt.Println("Login")
}

func (s *AuthServices) Register(body repositories.CreateUserRequest) (*repositories.NewUserResponse, error) {
	user, err := s.usersRepository.FindUserByEmail(body.Email)

	if err != nil {
		return nil, errors.New(customErrors.BadRequest)
	}

	if user.Email != "" {
		return nil, errors.New(customErrors.UserAlreadyExistsWithEmail)
	}

	newUser, error := s.usersRepository.CreateUser(body)

	if error != nil {
		return nil, errors.New(customErrors.BadRequest)
	}

	return newUser, nil
}

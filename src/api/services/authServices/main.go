package authServices

import (
	"errors"

	"github.com/felix-Asante/pennyPilot-go-api/src/api/repositories"
	"github.com/felix-Asante/pennyPilot-go-api/src/pkgs/security"
	customErrors "github.com/felix-Asante/pennyPilot-go-api/src/utils/errors"
)

type AuthServices struct {
	usersRepository *repositories.UsersRepository
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}
type ResetPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

func NewAuthServices(userRepository *repositories.UsersRepository) *AuthServices {
	return &AuthServices{userRepository}
}

func (s *AuthServices) Login(email string, password string) (*repositories.NewUserResponse, error) {
	user, err := s.usersRepository.FindUserByEmail(email)

	if err != nil {
		return nil, errors.New(customErrors.BadRequest)
	}

	if user.Email == "" {
		return nil, errors.New(customErrors.UserDoesNotExist)
	}

	if !security.IsPasswordValid(user.Password, password) {
		return nil, errors.New(customErrors.IncorrectPassword)
	}

	userResponse := &repositories.NewUserResponse{
		ID:        user.ID.String(),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
	return userResponse, nil
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

func (s *AuthServices) ResetPasswordRequest(email string) (string, error) {
	user, err := s.usersRepository.FindUserByEmail(email)

	if err != nil {
		return "", errors.New(customErrors.BadRequest)
	}

	if user.Email == "" {
		return "", errors.New(customErrors.UserDoesNotExist)
	}

	return "", nil
}

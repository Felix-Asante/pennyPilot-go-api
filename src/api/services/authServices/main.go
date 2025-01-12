package authServices

import (
	"errors"
	"time"

	"github.com/felix-Asante/pennyPilot-go-api/src/api/repositories"
	"github.com/felix-Asante/pennyPilot-go-api/src/pkgs/security"
	"github.com/felix-Asante/pennyPilot-go-api/src/utils"
	customErrors "github.com/felix-Asante/pennyPilot-go-api/src/utils/errors"
)

type AuthServices struct {
	usersRepository *repositories.UsersRepository
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}
type ForgetPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}
type ResetPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type ResetPassword struct {
	NewPassword string `json:"new_password" validate:"required,min=8"`
	Token       string `json:"token" validate:"required,min=4,max=4"`
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

	if user.ResetToken != "" && !hasTokenExpired(user.ResetTokenExpiry) {
		return "", errors.New(customErrors.AlreadyRequestedResetCodeError)
	}

	tokenCode := utils.GenerateRandomString(4)

	user.ResetToken = tokenCode
	user.ResetTokenCreatedAt = time.Now()
	user.ResetTokenExpiry = time.Now().Add(time.Minute * 15)

	_, error := s.usersRepository.Save(user)

	if error != nil {
		return "", errors.New(customErrors.BadRequest)
	}

	return tokenCode, nil
}

func (s *AuthServices) ResetPassword(body ResetPassword) error {
	user, err := s.usersRepository.FindUserByResetToken(body.Token)

	if err != nil {
		return errors.New(customErrors.InternalServerError)
	}

	if user.ResetToken == "" {
		return errors.New(customErrors.ResetTokenNotFound)
	}

	if hasTokenExpired(user.ResetTokenExpiry) {
		return errors.New(customErrors.ResetTokenExpiredError)
	}

	hashedPassword, error := security.GetHashedPassword(body.NewPassword)

	if error != nil {
		return errors.New(customErrors.InternalServerError)
	}

	user.Password = hashedPassword
	user.ResetToken = ""
	user.ResetTokenExpiry = time.Time{}
	user.ResetTokenCreatedAt = time.Time{}

	_, error = s.usersRepository.Save(user)

	if error != nil {
		return errors.New(customErrors.InternalServerError)
	}

	return nil
}

func hasTokenExpired(tokenExpiry time.Time) bool {
	return time.Now().After(tokenExpiry)
}

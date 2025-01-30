package customErrors

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type ValidationError struct {
	Message string
	Err     error
}

type StandardError struct {
	Status  int         `json:"status"`
	Error   string      `json:"error"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

var validate = validator.New()

const (
	InvalidRequestError            = "invalid request payload"
	ValidationFailedError          = "validation failed"
	UserDoesNotExist               = "user does not exist"
	UserAlreadyExists              = "user already exists"
	UserAlreadyExistsWithEmail     = "email taken"
	BadRequest                     = "bad request"
	SomethingWentWrong             = "something went wrong"
	InternalServerError            = "internal server error"
	PasswordNotMatch               = "password does not match"
	IncorrectPassword              = "incorrect password"
	ResetTokenNotExpired           = "reset token not expired"
	AlreadyRequestedResetCodeError = "already requested reset code"
	ResetTokenExpiredError         = "reset token expired"
	ResetTokenNotFound             = "reset token not found"
	AccountAlreadyExistsError      = "account already exists"
	ConflictError                  = "conflict"
	NotFoundError                  = "not found"
	InvalidAccountIDError          = "invalid account id"
	ForbiddenError                 = "forbidden error"
)

var StatusCodes = map[int]string{
	http.StatusBadRequest:          BadRequest,
	http.StatusInternalServerError: InternalServerError,
	http.StatusConflict:            ConflictError,
	http.StatusNotFound:            NotFoundError,
}

func DecodeAndValidate(r *http.Request, dst interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		fmt.Printf("\nValidation error %v\n", err)
		return &ValidationError{Message: InvalidRequestError}
	}

	if err := validate.Struct(dst); err != nil {
		fmt.Printf("\nValidation error %v\n", err)
		return &ValidationError{Message: ValidationFailedError, Err: err}
	}

	return nil
}

func (v *ValidationError) Error() string {
	if v.Err != nil {
		return v.Message + ": " + v.Err.Error()
	}
	return v.Message
}

func RespondWithError(w http.ResponseWriter, status int, errCode, message string, details interface{}) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(StandardError{
		Status:  status,
		Error:   errCode,
		Message: message,
		Details: details,
	})
}

func ValidateUUIDs(uuids []string) error {
	for _, id := range uuids {
		if _, err := uuid.Parse(id); err != nil {
			return fmt.Errorf("invalid UUID: %s", id)
		}
	}
	return nil
}

type AppError struct {
	StatusCode int
	Message    string
}

func (e *AppError) Error() string {
	return e.Message
}

func NewAppError(statusCode int, message string) *AppError {
	return &AppError{StatusCode: statusCode, Message: message}
}

package customErrors

import (
	"encoding/json"
	"strings"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

type MapError struct {
    Errors []struct {
        Field   string `json:"field"`
        Message string `json:"message"`
    }
}

func (e *MapError) Error() string {
    data, _ := json.Marshal(e.Errors)
    return string(data)
}

func NewMapErrorFromValidation(errs validator.ValidationErrors, trans ut.Translator) *MapError {
    var errors []struct {
        Field   string `json:"field"`
        Message string `json:"message"`
    }
    for _, err := range errs {
        errors = append(errors, struct {
            Field   string `json:"field"`
            Message string `json:"message"`
        }{Field: strings.ToLower(err.Field()), Message: err.Translate(trans)})
    }
    return &MapError{Errors: errors}
}
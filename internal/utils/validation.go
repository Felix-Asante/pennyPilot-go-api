package utils

import (
	"log/slog"
	"net/http"

	customErrors "github.com/Felix-Asante/pennyPilot-go-api/internal/errors"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

var (
	Validate *validator.Validate
	trans    ut.Translator
	uni      *ut.UniversalTranslator
)

func InitializeValidator() *validator.Validate {
	Validate = validator.New(validator.WithRequiredStructEnabled())

	en := en.New()
	uni = ut.New(en, en)

	trans, _ = uni.GetTranslator("en")
	err := en_translations.RegisterDefaultTranslations(Validate, trans)

	if err != nil {
		slog.Error("Failed to register default translations", "error", err.Error())
		return nil
	}

	return Validate
}

func ReadAndValidateJSON(w http.ResponseWriter, r *http.Request, data any) error {
	if err := ReadJSON(w, r, data); err != nil {
		return err
	}

	if err := Validate.Struct(data); err != nil {
		errs := err.(validator.ValidationErrors)

		return customErrors.NewMapErrorFromValidation(errs, trans)
	}

	return nil
}

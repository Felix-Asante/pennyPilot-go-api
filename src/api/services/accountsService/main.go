package accountsServices

import (
	"errors"

	"github.com/felix-Asante/pennyPilot-go-api/src/api/repositories"
	customErrors "github.com/felix-Asante/pennyPilot-go-api/src/utils/errors"
)

type AccountsServices struct {
	accountsRepository *repositories.AccountsRepository
}

func NewAccountServices(accountsRepository *repositories.AccountsRepository) *AccountsServices {
	return &AccountsServices{accountsRepository}
}

func (s *AccountsServices) Create(data repositories.CreateAccountDto) (*repositories.NewAccountResponse, error) {
	account, err := s.accountsRepository.FindByNameAndUserID(data.Name, data.UserID)

	if err != nil {
		return nil, errors.New(customErrors.InternalServerError)
	}

	if account.Name != "" {
		return nil, errors.New(customErrors.AccountAlreadyExistsError)
	}

	return s.accountsRepository.Create(data)
}

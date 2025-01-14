package accountsServices

import (
	"errors"
	"net/http"

	"github.com/felix-Asante/pennyPilot-go-api/src/api/repositories"
	customErrors "github.com/felix-Asante/pennyPilot-go-api/src/utils/errors"
)

type AccountsServices struct {
	accountsRepository *repositories.AccountsRepository
}

func NewAccountServices(accountsRepository *repositories.AccountsRepository) *AccountsServices {
	return &AccountsServices{accountsRepository}
}

func (s *AccountsServices) Create(data repositories.CreateAccountDto) (*repositories.NewAccountResponse, int, error) {
	account, err := s.accountsRepository.FindByNameAndUserID(data.Name, data.UserID)

	if err != nil {
		return nil, http.StatusInternalServerError, errors.New(customErrors.InternalServerError)
	}

	if account.Name != "" {
		return nil, http.StatusConflict, errors.New(customErrors.AccountAlreadyExistsError)
	}

	newAccount, error := s.accountsRepository.Create(data)
	return newAccount, http.StatusCreated, error
}

func (s *AccountsServices) Find(accountId string, userId string) (*repositories.NewAccountResponse, int, error) {

	account, err := s.accountsRepository.FindByIDAndUserID(accountId, userId)

	if err != nil {
		return nil, http.StatusInternalServerError, errors.New(customErrors.InternalServerError)
	}

	if account.Name == "" {
		return nil, http.StatusNotFound, errors.New(customErrors.NotFoundError)
	}

	return account, http.StatusOK, nil
}

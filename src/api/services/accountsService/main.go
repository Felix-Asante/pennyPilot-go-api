package accountsServices

import "github.com/felix-Asante/pennyPilot-go-api/src/api/repositories"

type AccountsServices struct {
	accountsRepository *repositories.AccountsRepository
}

func NewAuthServices(accountsRepository *repositories.AccountsRepository) *AccountsServices {
	return &AccountsServices{accountsRepository}
}

func (s *AccountsServices) Create() {}

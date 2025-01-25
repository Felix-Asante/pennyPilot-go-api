package accountsServices

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/felix-Asante/pennyPilot-go-api/src/api/repositories"
	"github.com/felix-Asante/pennyPilot-go-api/src/api/services/usersServices"
	customErrors "github.com/felix-Asante/pennyPilot-go-api/src/utils/errors"
	"gorm.io/gorm"
)

type AccountsServices struct {
	accountsRepository *repositories.AccountsRepository
	DB                 *gorm.DB
}

func NewAccountServices(db *gorm.DB) *AccountsServices {
	return &AccountsServices{accountsRepository: repositories.NewAccountsRepository(db), DB: db}
}

func (s *AccountsServices) Create(data repositories.CreateAccountDto) (*repositories.NewAccountResponse, int, error) {

	tx := s.DB.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	account, err := s.accountsRepository.FindByNameAndUserID(data.Name, data.UserID)

	if err != nil {
		return nil, http.StatusInternalServerError, errors.New(customErrors.InternalServerError)
	}

	if account.Name != "" {
		return nil, http.StatusConflict, errors.New(customErrors.AccountAlreadyExistsError)
	}

	usersService := usersServices.NewUsersServices(s.DB)
	user, _ := usersService.FindUserById(data.UserID)

	if user.Email == "" {
		return nil, http.StatusFound, errors.New(fmt.Sprintf("user %s", customErrors.NotFoundError))
	}

	if user.TotalAllocation+data.AllocationPoint > 100 {
		return nil, http.StatusBadRequest, errors.New("you have exceeded the maximum allocation")
	}

	user.TotalAllocation += data.AllocationPoint

	newAccount, error := s.accountsRepository.Create(data)
	_, err = usersService.SaveUser(user)

	if error != nil || err != nil {
		tx.Rollback()
	}
	if error := tx.Commit().Error; error != nil {
		return nil, http.StatusInternalServerError, errors.New(customErrors.InternalServerError)
	}

	return newAccount, http.StatusCreated, error
}

func (s *AccountsServices) Find(accountId string, userId string) (*repositories.Accounts, int, error) {

	account, err := s.accountsRepository.FindByIDAndUserID(accountId, userId)

	if err != nil {
		return nil, http.StatusInternalServerError, errors.New(customErrors.InternalServerError)
	}

	if account.Name == "" {
		message := fmt.Sprintf("Account %s", customErrors.NotFoundError)
		return nil, http.StatusNotFound, errors.New(message)
	}

	return account, http.StatusOK, nil
}

func (s *AccountsServices) UpdateBalance(accountId string, amount float64, user string) (*repositories.Accounts, int, error) {
	// Account belong to user
	account, statusCode, error := s.Find(accountId, user)

	if error != nil {
		return account, statusCode, error
	}

	newCurrentBalance := account.CurrentBalance + amount

	if newCurrentBalance > account.TargetBalance {
		return nil, http.StatusBadRequest, errors.New("you have reached the limit of this account")
	}

	account.CurrentBalance = newCurrentBalance

	account, error = s.accountsRepository.Save(account)

	return account, statusCode, error
}

func (s *AccountsServices) Remove(accountId string, user string) (int, error) {
	_, statusCode, error := s.Find(accountId, user)

	if error != nil {
		return statusCode, error
	}

	_, err := s.accountsRepository.Remove(accountId)

	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func (s *AccountsServices) FindUserAccounts(userId string) (*[]repositories.Accounts, error) {
	accounts, err := s.accountsRepository.FindAllByUserID(userId)

	return accounts, err
}

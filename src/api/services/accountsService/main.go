package accountsServices

import (
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/felix-Asante/pennyPilot-go-api/src/api/repositories"
	goalsservice "github.com/felix-Asante/pennyPilot-go-api/src/api/services/goalsService"
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
func (s *AccountsServices) SaveAccounts(account *repositories.Accounts) (*repositories.Accounts, error) {
	return s.accountsRepository.Save(account)
}

func (s *AccountsServices) AllocateToGoals(accountId string, userId string, goals []string) (int, error) {
	if len(goals) == 0 {
		return http.StatusBadRequest, errors.New("choose goals to be allocated")
	}

	tx := s.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var wg sync.WaitGroup
	errCh := make(chan error, len(goals))

	for _, goal := range goals {
		wg.Add(1)
		go func(goal string) {
			defer wg.Done()
			if err := s.AllocateToGoal(tx, accountId, userId, goal); err != nil {
				errCh <- err
			}
		}(goal)
	}

	wg.Wait()
	close(errCh)

	for err := range errCh {
		tx.Rollback()
		if appErr, ok := err.(*customErrors.AppError); ok {
			return appErr.StatusCode, appErr
		}
		return http.StatusInternalServerError, err
	}

	if err := tx.Commit().Error; err != nil {
		return http.StatusInternalServerError, errors.New(customErrors.InternalServerError)
	}

	return http.StatusOK, nil
}

func (s *AccountsServices) AllocateToGoal(tx *gorm.DB, accountId string, userId string, goalId string) error {
	account, status, err := s.Find(accountId, userId)
	if err != nil {
		return customErrors.NewAppError(status, customErrors.InternalServerError)
	}

	if account.Name == "" {
		return customErrors.NewAppError(http.StatusNotFound, fmt.Sprintf("account %s", customErrors.NotFoundError))
	}

	goalsService := goalsservice.NewGoalsService(tx)
	goal, status, err := goalsService.Get(goalId, userId)
	if err != nil {
		return customErrors.NewAppError(status, customErrors.InternalServerError)
	}

	if goal.Name == "" {
		return customErrors.NewAppError(http.StatusNotFound, fmt.Sprintf("goal %s", customErrors.NotFoundError))
	}

	if account.CurrentBalance <= 0 {
		return customErrors.NewAppError(http.StatusBadRequest, "insufficient income")
	}

	amountToAllocate := (goal.AllocationPoint / 100) * account.CurrentBalance
	goal.CurrentBalance += amountToAllocate
	account.CurrentBalance -= amountToAllocate
	account.TotalAllocation += amountToAllocate

	if err := tx.Save(account).Error; err != nil {
		return customErrors.NewAppError(http.StatusInternalServerError, customErrors.InternalServerError)
	}
	if err := tx.Save(goal).Error; err != nil {
		return customErrors.NewAppError(http.StatusInternalServerError, customErrors.InternalServerError)
	}

	return nil
}

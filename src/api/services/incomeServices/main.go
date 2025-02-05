package incomeservices

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"errors"

	"github.com/felix-Asante/pennyPilot-go-api/src/api/repositories"
	accountsServices "github.com/felix-Asante/pennyPilot-go-api/src/api/services/accountsService"
	"github.com/felix-Asante/pennyPilot-go-api/src/api/services/transactionsService"
	"github.com/felix-Asante/pennyPilot-go-api/src/api/services/usersServices"
	customErrors "github.com/felix-Asante/pennyPilot-go-api/src/utils/errors"
	"gorm.io/gorm"
)

type IncomeServices struct {
	incomesRepository *repositories.IncomeRepository
	DB                *gorm.DB
}

func NewIncomeServices(db *gorm.DB) *IncomeServices {
	return &IncomeServices{incomesRepository: repositories.NewIncomeRepository(db), DB: db}
}

func (s *IncomeServices) Create(data repositories.CreateIncomeDto) (*repositories.Incomes, int, error) {
	tx := s.DB.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	income, err := s.incomesRepository.Create(data)

	if err != nil {
		// log.Fatalf("Failed to create Income %v", err)
		tx.Rollback()
		return nil, http.StatusInternalServerError, errors.New(customErrors.InternalServerError)
	}

	userRepo := repositories.NewUsersRepository(s.DB)
	user, error := userRepo.FindUserById(data.User)

	if error != nil {
		tx.Rollback()
		return nil, http.StatusNotFound, errors.New(customErrors.UserDoesNotExist)
	}

	totalIncome := user.TotalIncome + data.Amount
	user.TotalIncome = totalIncome

	if _, error = userRepo.Save(user); error != nil {
		tx.Rollback()
		return nil, http.StatusInternalServerError, errors.New(customErrors.InternalServerError)
	}

	currentTime := time.Now()
	transactionsService := transactionsService.NewTransactionsService(tx)
	newTransaction := repositories.CreateTransactionDto{
		User:        data.User,
		Description: "Income allocation",
		Amount:      data.Amount,
		Type:        repositories.Deposit,
		Date:        &currentTime,
	}
	if _, err := transactionsService.CreateTransaction(newTransaction); err != nil {
		tx.Rollback()
		return nil, http.StatusInternalServerError, errors.New(customErrors.InternalServerError)
	}

	if error := tx.Commit().Error; error != nil {
		return nil, http.StatusInternalServerError, errors.New(customErrors.InternalServerError)
	}

	return income, http.StatusOK, nil
}

func (s *IncomeServices) Update(incomeId string, userId string, data repositories.UpdateIncomeDto) (*repositories.Incomes, int, error) {
	tx := s.DB.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	income, status, err := s.Get(incomeId, userId)
	if err != nil {
		return nil, status, err
	}

	userRepo := repositories.NewUsersRepository(s.DB)
	user, error := userRepo.FindUserById(userId)

	if error != nil {
		return nil, http.StatusNotFound, errors.New(customErrors.UserDoesNotExist)
	}

	if income.UserId != userId {
		return nil, http.StatusForbidden, errors.New(customErrors.ForbiddenError)
	}

	newAmount, newType, newFrequency, newDate := data.Amount, data.Type, data.Frequency, data.DateReceived

	if newAmount != nil {
		income.Amount = *newAmount
		if user.TotalIncome > 0 {
			user.TotalIncome -= income.Amount
		}
		user.TotalIncome += *newAmount
	}

	if newType != nil {
		income.Type = repositories.IncomeType(*newType)
	}

	if newFrequency != nil {
		income.Frequency = repositories.IncomeFrequency(*newFrequency)
	}

	if newDate != nil {
		income.DateReceived = newDate
	}

	newIncome, IncomeErr := s.incomesRepository.Save(income)

	if IncomeErr != nil {
		tx.Rollback()
		return nil, http.StatusInternalServerError, errors.New(customErrors.InternalServerError)
	}

	totalIncome := user.TotalIncome + *data.Amount
	user.TotalIncome = totalIncome

	if _, error = userRepo.Save(user); error != nil {
		tx.Rollback()
		return nil, http.StatusInternalServerError, errors.New(customErrors.InternalServerError)
	}

	if error := tx.Commit().Error; error != nil {
		return nil, http.StatusInternalServerError, errors.New(customErrors.InternalServerError)
	}

	return newIncome, http.StatusOK, nil
}

func (s *IncomeServices) Get(incomeId string, userId string) (*repositories.Incomes, int, error) {

	income, err := s.incomesRepository.FindByIDAndUserID(incomeId, userId)

	if err != nil {
		return nil, http.StatusInternalServerError, errors.New(customErrors.InternalServerError)
	}

	if income.UserId == "" {
		return nil, http.StatusNotFound, errors.New(fmt.Sprintf(`income %v`, customErrors.NotFoundError))
	}

	return income, http.StatusOK, nil
}

func (s *IncomeServices) AllocateIncomeToAccounts(userId string, accounts []string) (int, error) {
	if len(accounts) == 0 {
		return http.StatusBadRequest, errors.New("choose accounts to be allocated")
	}

	tx := s.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var wg sync.WaitGroup
	errCh := make(chan error, len(accounts))

	for _, account := range accounts {
		wg.Add(1)
		go func(account string) {
			defer wg.Done()
			if err := s.AllocateIncome(tx, userId, account); err != nil {
				errCh <- err
			}
		}(account)
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

func (s *IncomeServices) AllocateIncome(tx *gorm.DB, userId string, accountId string) error {
	accountService := accountsServices.NewAccountServices(tx)
	account, status, err := accountService.Find(accountId, userId)
	if err != nil {
		return customErrors.NewAppError(status, customErrors.InternalServerError)
	}

	if account == nil {
		return customErrors.NewAppError(http.StatusNotFound, "account not found")
	}

	usersService := usersServices.NewUsersServices(tx)
	user, err := usersService.FindUserById(userId)
	if err != nil {
		return customErrors.NewAppError(http.StatusInternalServerError, customErrors.InternalServerError)
	}

	if user.TotalIncome <= 0 {
		return customErrors.NewAppError(http.StatusBadRequest, "insufficient income")
	}

	amountToAllocate := (account.AllocationPoint / 100) * user.TotalIncome
	account.CurrentBalance += amountToAllocate
	user.TotalIncome -= amountToAllocate
	user.TotalAllocation += amountToAllocate

	if err := tx.Save(account).Error; err != nil {
		return customErrors.NewAppError(http.StatusInternalServerError, customErrors.InternalServerError)
	}
	if err := tx.Save(user).Error; err != nil {
		return customErrors.NewAppError(http.StatusInternalServerError, customErrors.InternalServerError)
	}

	currentTime := time.Now()
	transactionsService := transactionsService.NewTransactionsService(tx)
	newTransaction := repositories.CreateTransactionDto{
		User:        userId,
		Description: "Income allocation",
		Account:     &accountId,
		Amount:      amountToAllocate,
		Type:        repositories.Allocation,
		Date:        &currentTime,
	}
	if _, err := transactionsService.CreateTransaction(newTransaction); err != nil {
		return customErrors.NewAppError(http.StatusInternalServerError, customErrors.InternalServerError)
	}

	return nil
}

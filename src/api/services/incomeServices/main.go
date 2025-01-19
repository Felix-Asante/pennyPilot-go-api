package incomeservices

import (
	"fmt"
	"net/http"

	"errors"

	"github.com/felix-Asante/pennyPilot-go-api/src/api/repositories"
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
	fmt.Printf("Payload %v", data)
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

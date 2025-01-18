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

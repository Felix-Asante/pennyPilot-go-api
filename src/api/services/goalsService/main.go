package goalsservice

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/felix-Asante/pennyPilot-go-api/src/api/repositories"
	customErrors "github.com/felix-Asante/pennyPilot-go-api/src/utils/errors"
	"gorm.io/gorm"
)

type GoalsService struct {
	GoalsRepository repositories.GoalsRepository
	DB              *gorm.DB
}

type Goal *repositories.Goals

func NewGoalsService(db *gorm.DB) *GoalsService {
	return &GoalsService{GoalsRepository: *repositories.NewGoalsRepository(db), DB: db}
}

func (gs *GoalsService) Create(userId string, data repositories.CreateGoalDto) (Goal, int, error) {
	accountRepo := repositories.NewAccountsRepository(gs.DB)
	existingAccount, err := accountRepo.FindByIDAndUserID(data.Account, userId)
	total, _ := gs.GoalsRepository.FindAccountTotalAllocation(data.Account)

	fmt.Print(total)

	if err != nil {
		return nil, http.StatusInternalServerError, errors.New("an error occured while fetching account")
	}

	if existingAccount.Name == "" {
		return nil, http.StatusNotFound, errors.New(fmt.Sprintf("account %s", customErrors.NotFoundError))
	}

	newGoal := repositories.CreateNewGoalData{
		Name:            data.Name,
		Account:         data.Account,
		TargetAmount:    data.TargetAmount,
		CurrentBalance:  0,
		AllocationPoint: data.AllocationPoint,
		DueDate:         data.DueDate,
	}
	createdGoal, error := gs.GoalsRepository.Create(newGoal)

	if error != nil {
		return nil, http.StatusInternalServerError, errors.New("failed to create goal")
	}

	return createdGoal, http.StatusCreated, nil
}

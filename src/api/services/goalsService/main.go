package goalsservice

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/felix-Asante/pennyPilot-go-api/src/api/repositories"
	customErrors "github.com/felix-Asante/pennyPilot-go-api/src/utils/errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GoalsService struct {
	goalsRepository repositories.GoalsRepository
	accountRepo     repositories.AccountsRepository
	DB              *gorm.DB
}

type Goal *repositories.Goals

func NewGoalsService(db *gorm.DB) *GoalsService {
	return &GoalsService{goalsRepository: *repositories.NewGoalsRepository(db), DB: db, accountRepo: *repositories.NewAccountsRepository(db)}
}

func (gs *GoalsService) Create(userId string, data repositories.CreateGoalDto) (Goal, int, error) {

	existingAccount, err := gs.accountRepo.FindByIDAndUserID(data.Account, userId)

	if err != nil {
		return nil, http.StatusInternalServerError, errors.New("an error occured while fetching account")
	}

	if existingAccount.Name == "" {
		return nil, http.StatusNotFound, errors.New(fmt.Sprintf("account %s", customErrors.NotFoundError))
	}

	fmt.Println(data)

	newGoal := repositories.CreateNewGoalData{
		Name:            data.Name,
		Account:         data.Account,
		TargetAmount:    data.TargetAmount,
		CurrentBalance:  0,
		AllocationPoint: data.AllocationPoint,
		DueDate:         data.DueDate,
	}
	createdGoal, error := gs.goalsRepository.Create(newGoal)

	if error != nil {
		return nil, http.StatusInternalServerError, errors.New("failed to create goal")
	}

	return createdGoal, http.StatusCreated, nil
}

func (gs *GoalsService) Update(goalId string, userId string, data repositories.UpdateGoalDto) (Goal, int, error) {
	tx := gs.DB.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	existingGoal, err := gs.goalsRepository.FindByID(goalId)

	if err != nil {
		tx.Rollback()
		return nil, http.StatusInternalServerError, errors.New("failed to fetch goal")
	}

	if existingGoal.Name == "" {
		return nil, http.StatusNotFound, errors.New("goal not found")
	}

	existingAccount, err := gs.accountRepo.FindByIDAndUserID(existingGoal.AccountId, userId)

	if err != nil {
		tx.Rollback()
		return nil, http.StatusInternalServerError, errors.New(customErrors.InternalServerError)
	}

	if existingAccount.Name == "" {
		tx.Rollback()
		return nil, http.StatusForbidden, errors.New("you are not allowed to update this goal")
	}

	account, allocationPoint, dueDate, name, targetAmount := data.Account, data.AllocationPoint, data.DueDate, data.Name, data.TargetAmount

	if existingGoal.AccountId != account && account != "" {
		if err := uuid.Validate(account); err != nil {
			return nil, http.StatusBadRequest, errors.New("invalid account id")
		}

		existingAccount, err := gs.accountRepo.FindByIDAndUserID(account, userId)

		if err != nil {
			tx.Rollback()
			return nil, http.StatusInternalServerError, errors.New("failed to fetch account")
		}

		if existingAccount.Name == "" {
			tx.Rollback()
			return nil, http.StatusForbidden, errors.New("you are not allowed to update this goal")
		}

		existingAccount.CurrentBalance += existingGoal.CurrentBalance
		existingGoal.CurrentBalance = 0
		existingGoal.AccountId = account
	}

	updateField(allocationPoint != 0, func() {
		existingGoal.AllocationPoint = allocationPoint
	})

	updateField(dueDate != nil, func() {
		existingGoal.DueDate = dueDate
	})

	updateField(name != "", func() {
		existingGoal.Name = name
	})

	updateField(targetAmount != 0, func() {
		existingGoal.TargetAmount = targetAmount
	})

	if _, error := gs.goalsRepository.Save(existingGoal); error != nil {
		tx.Rollback()
		return nil, http.StatusInternalServerError, errors.New(customErrors.InternalServerError)
	}
	if _, error := gs.accountRepo.Save(existingAccount); error != nil {
		tx.Rollback()
		return nil, http.StatusInternalServerError, errors.New(customErrors.InternalServerError)
	}

	if error := tx.Commit().Error; error != nil {
		return nil, http.StatusInternalServerError, errors.New(customErrors.InternalServerError)
	}

	return existingGoal, http.StatusOK, nil
}

func (gs *GoalsService) Delete(goalId string, userId string) (int, error) {
	existingGoal, err := gs.goalsRepository.FindByID(goalId)
	if err != nil {
		return http.StatusInternalServerError, errors.New(customErrors.InternalServerError)
	}

	if existingGoal.Name == "" {
		return http.StatusNotFound, errors.New(customErrors.NotFoundError)
	}

	account, err := gs.accountRepo.FindByIDAndUserID(existingGoal.AccountId, userId)
	if err != nil {
		return http.StatusInternalServerError, errors.New(customErrors.InternalServerError)
	}
	if account.Name == "" {
		return http.StatusForbidden, errors.New("you are not allowed to delete this goal")
	}

	if isDeleted, _ := gs.goalsRepository.Remove(goalId); !isDeleted {
		return http.StatusInternalServerError, errors.New("failed to delete goal")
	}

	return http.StatusOK, nil
}

func (gs *GoalsService) Get(goalId string, userId string) (Goal, int, error) {
	existingGoal, err := gs.goalsRepository.FindByID(goalId)
	if err != nil {
		return nil, http.StatusInternalServerError, errors.New(customErrors.InternalServerError)
	}

	if existingGoal.Name == "" {
		return nil, http.StatusNotFound, errors.New(customErrors.NotFoundError)
	}

	account, err := gs.accountRepo.FindByIDAndUserID(existingGoal.AccountId, userId)
	if err != nil {
		return nil, http.StatusInternalServerError, errors.New(customErrors.InternalServerError)
	}
	if account.Name == "" {
		return nil, http.StatusForbidden, errors.New(fmt.Sprintf("goal %v", customErrors.ForbiddenError))
	}

	return existingGoal, http.StatusOK, nil
}

func updateField(condition bool, updateFunc func()) {
	if condition {
		updateFunc()
	}
}

func (gs *GoalsService) SaveGoal(goal *repositories.Goals) (*repositories.Goals, error) {
	return gs.goalsRepository.Save(goal)
}

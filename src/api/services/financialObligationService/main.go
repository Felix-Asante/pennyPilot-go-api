package financialObligationService

import (
	"errors"

	"github.com/felix-Asante/pennyPilot-go-api/src/api/repositories"
	customErrors "github.com/felix-Asante/pennyPilot-go-api/src/utils/errors"
	"gorm.io/gorm"
)

type FinancialObligationService struct {
	obligationsRepo *repositories.FinancialObligationsRepository
	usersRepo       *repositories.UsersRepository
	DB              *gorm.DB
}

type FinancialObligations *repositories.FinancialObligations

func NewFinancialObligationService(db *gorm.DB) *FinancialObligationService {
	return &FinancialObligationService{obligationsRepo: repositories.NewFinancialObligationsRepository(db), usersRepo: repositories.NewUsersRepository(db), DB: db}
}

func (fos *FinancialObligationService) Create(userId string, dto repositories.CreateFinancialObligationDto) (FinancialObligations, error) {

	user, err := fos.usersRepo.FindUserById(userId)
	if err != nil {
		return nil, err
	}

	if user.Email == "" {
		return nil, errors.New(customErrors.UserDoesNotExist)
	}

	newObligation := repositories.CreateFinancialObligation{
		Type:             repositories.FinancialObligationType(dto.Type),
		TotalAmount:      dto.TotalAmount,
		CounterpartyName: dto.CounterpartyName,
		RemainingAmount:  dto.RemainingAmount,
		RepaymentType:    repositories.FinancialObligationRepaymentType(dto.RepaymentType),
		NextDueDate:      dto.NextDueDate,
		InterestRate:     dto.InterestRate,
		UserId:           userId,
	}
	return fos.obligationsRepo.Create(newObligation)
}

func (fos *FinancialObligationService) FindAllByUserId(userId string) ([]repositories.FinancialObligations, error) {

	return fos.obligationsRepo.FindByUserID(userId)
}

func (fos *FinancialObligationService) FindById(userId string, id string) (FinancialObligations, error) {

	return fos.obligationsRepo.FindByID(id)
}

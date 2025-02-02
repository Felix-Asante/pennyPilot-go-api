package transactionsService

import (
	"github.com/felix-Asante/pennyPilot-go-api/src/api/repositories"
	"gorm.io/gorm"
)

type Transaction *repositories.Transaction

type TransactionsService struct {
	transactionsRepository *repositories.TransactionsRepository
	DB                     *gorm.DB
}

func NewTransactionsService(db *gorm.DB) *TransactionsService {
	return &TransactionsService{
		transactionsRepository: repositories.NewTransactionsRepository(db),
		DB:                     db,
	}
}

func (s *TransactionsService) CreateTransaction(transaction repositories.CreateTransactionDto) (Transaction, error) {
	return s.transactionsRepository.Create(transaction)
}

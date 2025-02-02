package repositories

import (
	"math"

	"gorm.io/gorm"
)

func SetUpRepositories(db *gorm.DB) {
	NewUsersRepository(db)
	NewAccountsRepository(db)
}

type PaginationResult struct {
	Items interface{}    `json:"items"`
	Meta  PaginationMeta `json:"meta"`
}

type PaginationMeta struct {
	CurrentPage int `json:"current_page"`
	TotalItems  int `json:"total_items"`
	TotalPages  int `json:"total_pages"`
}

func Paginate(db *gorm.DB, page, pageSize int, result interface{}) (PaginationResult, error) {
	var paginationResult PaginationResult

	query := db.Session(&gorm.Session{})

	offset := (page - 1) * pageSize

	err := query.Offset(offset).Limit(pageSize).Find(result).Error
	if err != nil {
		return paginationResult, err
	}

	var totalItems int64
	err = query.Model(result).Count(&totalItems).Error
	if err != nil {
		return paginationResult, err
	}

	totalPages := int(math.Ceil(float64(totalItems) / float64(pageSize)))

	paginationResult.Items = result
	paginationResult.Meta = PaginationMeta{
		CurrentPage: page,
		TotalItems:  int(totalItems),
		TotalPages:  totalPages,
	}

	return paginationResult, nil
}

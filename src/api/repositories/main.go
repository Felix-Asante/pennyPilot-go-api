package repositories

import "gorm.io/gorm"

func SetUpRepositories(db *gorm.DB) {
	NewUsersRepository(db)
	NewTokensRepository(db)
}

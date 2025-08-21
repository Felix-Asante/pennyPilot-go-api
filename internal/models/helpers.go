package models

import "gorm.io/gorm"

func getTxDB(db *gorm.DB, tx *gorm.DB) *gorm.DB {
	if tx != nil {
		return tx
	}
	return db
}

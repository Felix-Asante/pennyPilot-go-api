package db

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DbConfig struct {
	DbUser     string
	DbHost     string
	DbPassword string
	DbName     string
	DbPort     string
}

func NewPgDB(config DbConfig) *DbConfig {
	return &DbConfig{
		DbUser:     config.DbUser,
		DbHost:     config.DbHost,
		DbPassword: config.DbPassword,
		DbName:     config.DbName,
		DbPort:     config.DbPort,
	}
}

func (config *DbConfig) Init() (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s dbname=%s port=%s sslmode=disable password=%s", config.DbHost, config.DbUser, config.DbName, config.DbPort, config.DbPassword)

	return gorm.Open(postgres.Open(dsn), &gorm.Config{FullSaveAssociations: true})
}

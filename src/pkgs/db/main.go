package db

import (
	"fmt"
	"log"

	"github.com/felix-Asante/pennyPilot-go-api/src/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConnectToDB() *gorm.DB {
	db_host := utils.GetEnv("DB_HOST")
	db_port := utils.GetEnv("DB_PORT")
	db_user := utils.GetEnv("DB_USER")
	db_name := utils.GetEnv("DB_NAME")
	db_password := utils.GetEnv("DB_PASSWORD")

	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=disable", db_host, db_user, db_password, db_name, db_port)
	fmt.Println(dsn)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("failed to connect database%v", err)
		panic("failed to connect database")
	}

	fmt.Println("Connected to database")
	db.Logger.LogMode(logger.Info)
	return db
}

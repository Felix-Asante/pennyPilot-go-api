package main

import (
	"fmt"
	"log"

	"github.com/felix-Asante/pennyPilot-go-api/src/api/repositories"
	"github.com/felix-Asante/pennyPilot-go-api/src/pkgs/db"
	"gorm.io/gorm"
)

func main() {
	db := db.ConnectToDB()
	createDBTypes(db)
	error := db.AutoMigrate(&repositories.Users{}, &repositories.Accounts{}, &repositories.Incomes{}, &repositories.Goals{}, &repositories.FinancialObligations{})

	if error != nil {
		panic(error)
	}
	fmt.Println("Migration successful")
}

func createDBTypes(db *gorm.DB) {
	error := db.Exec(`
	DO $$
	BEGIN
		IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'income_type') THEN
			CREATE TYPE income_type AS ENUM ('salary', 'investment', 'freelance', 'others');
		END IF;
		IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'income_frequency') THEN
			CREATE TYPE income_frequency AS ENUM ('one-time', 'weekly', 'bi-weekly', 'monthly', 'yearly');
		END IF;
		IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'financial_obligation_type') THEN
			CREATE TYPE financial_obligation_type AS ENUM ('debt', 'lend');
		END IF;
		IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'financial_obligation_repayment_type') THEN
			CREATE TYPE financial_obligation_repayment_type AS ENUM ('weekly', 'daily', 'monthly', 'yearly');
		END IF;
	END
	$$;
`).Error

	if error != nil {
		log.Fatalf("Failed to create types: %v", error)
		panic(error)
	}

	fmt.Println("Created types")
}

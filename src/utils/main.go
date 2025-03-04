package utils

import (
	"math/rand"
)

func GenerateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

type IncomeType string
type IncomeFrequency string

const (
	Salary     IncomeType = "salary"
	Investment IncomeType = "investment"
	Freelance  IncomeType = "freelance"
	Others     IncomeType = "others"
)
const (
	Weekly   IncomeFrequency = "weekly"
	Monthly  IncomeFrequency = "monthly"
	Yearly   IncomeFrequency = "yearly"
	BiWeekly IncomeFrequency = "bi-weekly"
	OneTime  IncomeFrequency = "one-time"
)

func SetDefaultValues(condition bool, updateFunc func()) {
	if condition {
		updateFunc()
	}
}

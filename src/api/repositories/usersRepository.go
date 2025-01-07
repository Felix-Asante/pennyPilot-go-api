package repositories

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Users struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4();column:id"`
	FirstName string    `gorm:"column:first_name"`
	LastName  string    `gorm:"column:last_name"`
	Email     string    `gorm:"column:email;unique"`
	Password  string    `gorm:"column:password"`
	gorm.Model
}
type UsersRepository struct {
	db *gorm.DB
}

func (u *Users) TableName() string {
	return "users"
}

func NewUsersRepository(db *gorm.DB) *UsersRepository {

	return &UsersRepository{db}
}

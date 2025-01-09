package repositories

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Users struct {
	ID        uuid.UUID  `gorm:"type:uuid;primaryKey;default:uuid_generate_v4();column:id"`
	FirstName string     `gorm:"column:first_name"`
	LastName  string     `gorm:"column:last_name"`
	Email     string     `gorm:"column:email;unique"`
	Password  string     `gorm:"column:password"`
	CreatedAt *time.Time `gorm:"column:created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at"`
}

type CreateUserRequest struct {
	FirstName string `json:"first_name" validate:"required,min=2"`
	LastName  string `json:"last_name" validate:"required,min=2"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8"`
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

func (u *UsersRepository) FindUserByEmail(email string) (*Users, error) {
	var existingUser Users

	error := u.db.Where("email = ?", email).Find(&existingUser).Error

	return &existingUser, error
}

func (u *UsersRepository) CreateUser(data CreateUserRequest) (*Users, error) {
	user := Users{
		FirstName: data.FirstName,
		LastName:  data.LastName,
		Email:     data.Email,
		Password:  data.Password,
	}

	error := u.db.Create(&user).Error

	return &user, error
}

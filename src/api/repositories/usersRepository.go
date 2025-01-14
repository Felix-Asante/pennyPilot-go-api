package repositories

import (
	"time"

	"github.com/felix-Asante/pennyPilot-go-api/src/pkgs/security"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Users struct {
	ID                  uuid.UUID  `gorm:"type:uuid;primaryKey;default:uuid_generate_v4();column:id"`
	FirstName           string     `gorm:"column:first_name"`
	LastName            string     `gorm:"column:last_name"`
	Email               string     `gorm:"column:email;unique"`
	Password            string     `gorm:"column:password"`
	ResetToken          string     `gorm:"column:reset_token"`
	ResetTokenExpiry    time.Time  `gorm:"column:reset_token_expiry"`
	ResetTokenCreatedAt time.Time  `gorm:"column:reset_token_created_at"`
	CreatedAt           *time.Time `gorm:"column:created_at"`
	UpdatedAt           *time.Time `gorm:"column:updated_at"`
	Accounts            []Accounts `gorm:"columns:accounts;foreignKey:UserID;references:ID"`
}

type CreateUserRequest struct {
	FirstName string `json:"first_name" validate:"required,min=2"`
	LastName  string `json:"last_name" validate:"required,min=2"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8"`
}

type NewUserResponse struct {
	ID        string     `json:"id"`
	FirstName string     `json:"first_name"`
	LastName  string     `json:"last_name"`
	Email     string     `json:"email"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
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

func (u *UsersRepository) CreateUser(data CreateUserRequest) (*NewUserResponse, error) {
	user := Users{
		FirstName: data.FirstName,
		LastName:  data.LastName,
		Email:     data.Email,
		Password:  data.Password,
	}

	error := u.db.Create(&user).Error

	newUser := NewUserResponse{
		ID:        user.ID.String(),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return &newUser, error
}

func (u *UsersRepository) FindUserById(id string) (*Users, error) {
	var existingUser Users

	error := u.db.Where("id = ?", id).Find(&existingUser).Error

	return &existingUser, error
}

func (u *UsersRepository) FindUserByResetToken(token string) (*Users, error) {
	var existingUser Users

	error := u.db.Where("reset_token = ?", token).Find(&existingUser).Error

	return &existingUser, error
}

func (u *UsersRepository) Save(user *Users) (*Users, error) {

	error := u.db.Save(user).Error

	return user, error
}

func (u *Users) BeforeCreate(tx *gorm.DB) error {

	password, error := security.GetHashedPassword(u.Password)

	if error != nil {
		return error
	}
	u.Password = password
	return nil
}

package repositories

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Tokens struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4();column:id"`
	Token     string    `gorm:"column:token;unique;not null"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index;column:user;not null"`
	Expiry    time.Time `gorm:"not null;column:expiry"`
	CreatedAt time.Time `gorm:"autoCreateTime;column:created_at;not null"`
	User      Users     `gorm:"foreignKey:UserID"`
}

type TokensRepository struct {
	db *gorm.DB
}

func (u *Tokens) TableName() string {
	return "tokens"
}

func NewTokensRepository(db *gorm.DB) *TokensRepository {

	return &TokensRepository{db}
}

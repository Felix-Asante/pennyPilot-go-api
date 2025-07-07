package models

type User struct {
	ID           string `gorm:"type:uuid;primaryKey;column:id;index"`
	Email        string `gorm:"column:email;unique;index"`
	FullName     string `gorm:"column:full_name"`
	PasswordHash string `gorm:"column:password_hash"`
	Currency     string `gorm:"column:currency"`
	CreatedAt    int64  `gorm:"autoCreateTime:milli"`
	UpdatedAt    int64  `gorm:"autoUpdateTime:milli"`
}

package models

import "time"

type UserSerializer struct {
	ID           string `json:"id"`
	Email        string `json:"email"`
	FullName     string `json:"full_name"`
	Currency     string `json:"currency"`
	CreatedAt    string  `json:"created_at"`
	UpdatedAt    string  `json:"updated_at"`
}

func SerializeUser(user *User) *UserSerializer {
	return &UserSerializer{
		ID:           user.ID,
		Email:        user.Email,
		FullName:     user.FullName,
		Currency:     user.Currency,
		CreatedAt:    time.Unix(user.CreatedAt, 0).Format("2006-01-02 15:04:05"),
		UpdatedAt:    time.Unix(user.UpdatedAt, 0).Format("2006-01-02 15:04:05"),
	}
}
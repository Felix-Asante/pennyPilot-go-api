package usersServices

import (
	"fmt"

	"github.com/felix-Asante/pennyPilot-go-api/src/api/repositories"
	"gorm.io/gorm"
)

type UsersServices struct {
	userRepository *repositories.UsersRepository
}

func NewUsersServices(db *gorm.DB) *UsersServices {
	return &UsersServices{userRepository: repositories.NewUsersRepository(db)}
}

func (s *UsersServices) FindAccounts(id string) {}

func (s *UsersServices) CreateNewUser(id string) {
	fmt.Println("Find by id")
}

func (s *UsersServices) FindUserById(id string) (*repositories.Users, error) {
	return s.userRepository.FindUserById(id)
}
func (s *UsersServices) SaveUser(user *repositories.Users) (*repositories.Users, error) {
	return s.userRepository.Save(user)
}

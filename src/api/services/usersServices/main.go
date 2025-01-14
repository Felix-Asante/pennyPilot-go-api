package usersServices

import (
	"fmt"

	"github.com/felix-Asante/pennyPilot-go-api/src/api/repositories"
)

type UsersServices struct {
	userRepository *repositories.UsersRepository
}

func NewUsersServices(userRepository *repositories.UsersRepository) *UsersServices {
	return &UsersServices{userRepository}
}

func (s *UsersServices) FindAccounts(id string) {}

func (s *UsersServices) CreateNewUser(id string) {
	fmt.Println("Find by id")
}

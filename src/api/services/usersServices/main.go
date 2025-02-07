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
func (s *UsersServices) Me(id string) (*repositories.NewUserResponse, error) {
	user, error := s.userRepository.FindUserById(id)

	userResponse := &repositories.NewUserResponse{
		ID:              user.ID.String(),
		FirstName:       user.FirstName,
		LastName:        user.LastName,
		Email:           user.Email,
		CreatedAt:       user.CreatedAt,
		UpdatedAt:       user.UpdatedAt,
		MembershipId:    user.MembershipId,
		TotalIncome:     user.TotalIncome,
		TotalAllocation: user.TotalAllocation,
		AllocatedAmount: user.AllocatedAmount,
	}
	return userResponse, error
}

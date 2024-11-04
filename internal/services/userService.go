package services

import (
	"errors"

	"github.com/Ansalps/UserEcommerceClean/internal/models"
	"github.com/Ansalps/UserEcommerceClean/internal/repository"
)

type IUserService interface {
	UserSignUp(user *models.User) error
	UserLogin(user *models.UserLogin) (*models.User, error)
	ComparePassword(providedUser models.UserLogin, user models.User) bool
	GetProfile(userID string) (*models.User, error)
	UpdateProfile(userID uint, user models.UserUpdate) error
}
type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}
func (c *UserService) UserSignUp(user *models.User) error {
	existingUser, _ := c.userRepo.GetUserByEmail(user.Email)
	if existingUser != nil {
		return errors.New(models.UserAlreadyExists)
	}
	err := c.userRepo.UserSignUp(*user)
	if err != nil {
		return err
	}
	return nil
}
func (c *UserService) UserLogin(user *models.UserLogin) (*models.User, error) {
	User, err := c.userRepo.GetUserByEmail(user.Email)
	if err != nil {
		return nil, err
	}
	return User, nil
}
func (c *UserService) ComparePassword(providedUser models.UserLogin, user models.User) bool {
	//c.userRepo.ComparePassword(providedUser.Password, user.Password)
	check := false
	if providedUser.Password == user.Password {
		check = true
	}
	return check
}
func (c *UserService) GetProfile(userID string) (*models.User, error) {
	user, err := c.userRepo.GetProfile(userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}
func (c *UserService) UpdateProfile(userID uint, user models.UserUpdate) error {
	User, err := c.userRepo.GetUserById(userID)
	if err != nil {
		return err
	}
	User.FirstName = user.FirstName
	User.LastName = user.LastName
	User.Phone = user.Phone
	err = c.userRepo.UpdateProfile(User)
	if err != nil {
		return err
	}
	return nil
}

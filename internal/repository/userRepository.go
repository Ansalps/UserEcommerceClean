package repository

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/Ansalps/UserEcommerceClean/internal/models"
	"gorm.io/gorm"
)

type IUserRepository interface {
	GetUserByEmail(email string) (*models.User, error)
	GetUserById(userID uint) (*models.User, error)
	GetProfile(userId string) (*models.User, error)
	UpdateProfile(user *models.User) error
}
type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}
func (c *UserRepository) UserSignUp(user models.User) error {
	err := c.db.Create(&user).Error
	if err != nil {
		return err
	}
	return nil
}

// func (c *UserRepository) GetUser(field string, value interface{}) {

// }
func (c *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := c.db.Where("email=?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		} else {
			return nil, err
		}
	}
	return &user, nil
}
func (c *UserRepository) GetUserById(userID uint) (*models.User, error) {
	var user models.User
	err := c.db.Where("id=?", userID).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		} else {
			return nil, err
		}
	}
	return &user, nil
}
func (c *UserRepository) GetProfile(userId string) (*models.User, error) {
	var user models.User

	num, err := strconv.ParseUint(userId, 10, 32)
	if err != nil {
		fmt.Println("error in parsing string")
		return nil, err
	}
	err = c.db.Where("id = ?", num).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
func (c *UserRepository) UpdateProfile(user *models.User) error {

	err := c.db.Save(user).Error
	if err != nil {
		return errors.New("error in update of databse query")
	}
	return nil
}

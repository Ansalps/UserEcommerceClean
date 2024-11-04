package database

import (
	"fmt"

	"github.com/Ansalps/UserEcommerceClean/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Initialize() {
	var err error
	dsn := "postgres://postgres:123@localhost:5432/userecommerce"
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("connection failed due to ", err)
	}
}
func AutoMigrate() {
	DB.AutoMigrate(&models.User{})
}

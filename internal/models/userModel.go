package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	//ID        uint   `gorm:"primary key" json:"id"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Email     string `json:"email"  validate:"required"`
	Password  string `json:"password" validate:"required"`
	Phone     string `json:"phone" validate:"required,numeric,len=10"`
	Status    string `gorm:"type:varchar(10); check(status IN ('Active', 'Blocked', 'Deleted')) ;default:'Active'" json:"status"`
}
type UserLogin struct {
	Email    string `gorm:"unique" validate:"required,email" json:"email"`
	Password string `validate:"required" json:"password"`
}
type UserUpdate struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Phone     string `json:"phone" validate:"required,numeric,len=10"`
}

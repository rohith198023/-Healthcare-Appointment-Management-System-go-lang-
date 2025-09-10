package models

import (
	"time"
	"gorm.io/gorm"
)

type Patient struct{
	gorm.Model
	Name string `json:"name"`
	Role string `json:"role"`
	Password string `json:"password"`
	Phone string `gorm:"uniqueindex" json:"phone"`
	Email string `gorm:"uniqueindex" json:"email" `
	OTP string	`json:"-"`
	OTPExpiry time.Time	`json:"-"`
}

package models

import (
	"time"

	"gorm.io/gorm"
)

type Admin struct {
    gorm.Model
    Name     string `json:"name"`
    Email    string `gorm:"unique" json:"email"`
    Password string `json:"password"`
    Role string `json:"admin"`
    OTP string `json:"otp"`
    OTPExpiry *time.Time `json:"otp_expiry"`
}

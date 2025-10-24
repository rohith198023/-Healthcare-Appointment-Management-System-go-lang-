package models

import (
	"time"
	"gorm.io/gorm"
)

type Doctor struct{
	gorm.Model
	Name string `json:"name"`
	Role string `json:"role"`
	Email string `json:"Email"`
	Password string `json:"Password"`
	Specialization string `gorm:"index" json:"specialization"`
	AvgconsultMin int `json:"avg_consult_min"`
	OTP string
	OTPExpiry time.Time
}


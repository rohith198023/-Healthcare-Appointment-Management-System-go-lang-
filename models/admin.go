package models

import (
	"gorm.io/gorm"
)

type Admin struct {                             
    gorm.Model
    Name     string `json:"name"`
    Email    string `gorm:"unique" json:"email"`
    Password string `json:"password"`
    Role string `json:"role"`
}













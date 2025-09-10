

package models

import (
	"time"

	"gorm.io/gorm"
)

type Notification struct {
	gorm.Model
	PatientID uint      `json:"patient_id"`
	DoctorID  uint      `json:"doctor_id"`
	Title     string    `json:"title"`
	Message   string    `json:"message"`
	Read      bool      `json:"read" gorm:"default:false"`
	CreatedAt time.Time `json:"created_at"`
}














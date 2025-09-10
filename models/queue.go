package models

import (
	"gorm.io/gorm"
	"time"
)

type Queue struct{
	gorm.Model
	DoctorId int `gorm:"index" json:"doctor_id"`
	PatientId int `gorm:"index" json:"patient_id"`
	AppointmentId int `gorm:"index" json:"appointment_id"`
	Position int `gorm:"index" json:"position"`
	EstimatedStartUtc time.Time `json:"estimated_start_time"`
	Status string `gorm:"index" json:"status"`
} 












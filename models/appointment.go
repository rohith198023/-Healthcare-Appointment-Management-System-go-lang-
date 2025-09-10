package models

import (
	"time"
	"gorm.io/gorm"
)

type Appointmentstatus string


const(
	Statusactive Appointmentstatus="Active"
	Statusinactive Appointmentstatus="Inactive"
)

type Appointment struct{
	gorm.Model
	PatientId uint `gorm:"index" json:"patientId"`
	DoctorId uint `gorm:"index" json:"doctorID"`
	ScheduleID   uint     `json:"schedule_id"`
	ScheduledAt time.Time `json:"scheduled_time"`
	Status Appointmentstatus `json:"Status"`
	EstimatedMinutes int     `json:"estimated_duration"`
 
	Patient    Patient   `gorm:"foreignKey:PatientId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
    Doctor  Doctor   `gorm:"foreignKey:DoctorId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Schedule Schedule `gorm:"foreignKey:ScheduleID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	 Queue     Queue     `gorm:"foreignKey:AppointmentId"`
}



























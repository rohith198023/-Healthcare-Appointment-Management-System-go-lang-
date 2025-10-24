package models

import(
	"time"
)

type Schedule struct{
	ID        uint      `gorm:"primaryKey"`
	Date      time.Time `json:"date"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	SlotSize  int       `json:"slot_size"`
	Status    string    `json:"status"`

	// Foreign Key to Doctor
	DoctorID uint   `json:"doctor_id"`
	Doctor   Doctor `gorm:"foreignKey:DoctorID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	// One-to-many: A schedule can have multiple appointments
	Appointments []Appointment `gorm:"foreignKey:ScheduleID"`
}

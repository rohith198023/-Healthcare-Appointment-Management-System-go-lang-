package handlers

import(
	"project/database"
	"project/models"
	"time"
)

func Sendnotification(patientID uint, doctorID uint, title string, message string)error{
	notifications:=models.Notification{
		PatientID: patientID,
		DoctorID:  doctorID,
		Title:     title,
		Message:   message,
		Read:      false,
		CreatedAt: time.Now(),
	}

	return database.DB.Create(&notifications).Error
}

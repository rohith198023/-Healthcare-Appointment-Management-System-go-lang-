package handlers

import (

	"time"

	"project/database"
	"project/models"
)

// Recalculate queue for a given doctor after appointment is completed/cancelled
func RecalculateQueue(doctorID uint) error {
	 queues:= []models.Queue{}
	if err := database.DB.
		Joins("JOIN appointments ON appointments.id = queues.appointment_id").
		Where("appointments.doctor_id = ? AND queues.status != ?", doctorID, "completed").
		Order("queues.position ASC").
		Find(&queues).Error; err != nil {
		return err
	} 
	for i, q := range queues {
		newPos := i + 1
		if q.Position != newPos {
			q.Position = newPos
			q.EstimatedStartUtc = time.Now().Add(time.Duration(i*10) * time.Minute)
			database.DB.Save(&q)
		}
		
		if newPos == 1 {
			Sendnotification(uint(q.PatientId), uint(q.DoctorId),
				"It's your turn",
				"Please be ready to meet the doctor.",
			)
		} else if newPos == 2 {
			Sendnotification(uint(q.PatientId), uint(q.DoctorId),
				"Be Ready",
				"Next is your turn. Please be prepared.",
			)
		}
	}

	return nil
}


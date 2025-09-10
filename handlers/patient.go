package handlers

import (
	"fmt"
	"math/rand"
	"project/database"
	"project/models"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"project/utils"
)


func PatientsignUp(c *fiber.Ctx)error{
		var User models.Patient
		if err:=c.BodyParser(&User); err!=nil{
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})	
		}

		if User.Name==""||User.Email==""||User.Password==""{
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "name, email and password are required",
			})
		}

		var existing models.Patient
		if err:=database.DB.Where("email=?",User.Email).First(&existing).Error
		 err==nil{
			c.Status(404).JSON(&fiber.Map{
				"message":"email already exists",
			})
			return err
		}

		err:=database.DB.Create(&User).Error
		if err!=nil{
			c.Status(404).JSON(&fiber.Map{
				"message":"failes to register user",
				"details":err.Error(),
			})
			return nil
		}

		c.Status(201).JSON(&fiber.Map{
			"message":"Patient created successfully",
			"Patient":User,
		})
		return nil	
}

func PatientLogin(c *fiber.Ctx)error{
	var credentials struct{
		Email string `json:"email"`
		Password string `json:"password"`
	}

	if err:=c.BodyParser(&credentials);err!=nil{
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":"Invalid request Body",
			"details":err.Error(),
		})
	}

	if credentials.Email ==""||credentials.Password==""{
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":"email and password required",
		})
	}

	var user models.Patient

	credentials.Email=strings.TrimSpace(strings.ToLower(credentials.Email))
	credentials.Password=strings.TrimSpace(credentials.Password)


	if err:=database.DB.Where("email = ? AND password=?", credentials.Email, credentials.Password).First(&user).Error; err!=nil{
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":"Invalid email and password",
		})
	}


	otp:=fmt.Sprintf("%05d",rand.Intn(1000))
	expiry:=time.Now().Add(5 *time.Minute)
	user.OTP=otp
	user.OTPExpiry=expiry
	database.DB.Save(&user)

	if err:=utils.SendOtp(user.Email,otp);err!=nil{
		return c.Status(500).JSON(fiber.Map{
			"error":"failed to send the otp",
		})
	}

	return c.JSON(fiber.Map{
		"message":"Otp is sented to your Email",
	})  
}

func VerifyOTPPatient(c *fiber.Ctx) error{

	var payload struct{
		Email string `json:"email"`
		OTP string `json:"otp"`
	}     	

	if err:=c.BodyParser(&payload);err!=nil{
		return c.Status(400).JSON(fiber.Map{
			"error":"Invalid request",
		})
	}

	var user models.Patient
	if err:=database.DB.Where("email=?",payload.Email).First(&user).Error;err!=nil{
		return c.Status(401).JSON(fiber.Map{
			"error":"Invalid email",
		})
	}

	if user.OTP!=payload.OTP||time.Now().After(user.OTPExpiry){
		return c.Status(401).JSON(fiber.Map{
			"error":"Invalid or expired otp",
		})
	}

	user.OTP=""
	user.OTPExpiry=time.Time{}
	database.DB.Save(&user)

	token, err:=utils.GenrateToken(user.ID,user.Role)
	if err!=nil{
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":"failed to generate token",
		})
	}

	return c.JSON(fiber.Map{
		"message":"login successful",
		"token":token,
	})

	
}

func GetMyProfile(c *fiber.Ctx)error{
	id:=c.Locals("user_id")
	patientByid:=&models.Patient{}
	if id ==""{
		c.Status(404).JSON(fiber.Map{
			"Message":"id cannot be empty",
		})
		return nil
	}
	fmt.Println("the Id is",id)

	err:=database.DB.Where("id=?",id).First(patientByid).Error
	if err!=nil{         
		c.Status(404).JSON(fiber.Map{
			"message":"could not get the Patients",
		})
		return err
	}

	c.Status(200).JSON(fiber.Map{
		"message":"Patient id fetched successfuly",
		"Name":patientByid.Name,
		"Phone":patientByid.Phone,
		"Email":patientByid.Email,
	})
	return nil
}

func UpdateMyProfile(c *fiber.Ctx)error{
	id:=c.Locals("user_id")
	var patient models.Patient

	if err:=database.DB.First(&patient,id).Error;err!=nil{
		c.Status(400).JSON(fiber.Map{
			"Message":"Patient not found",
		})
		return err
	}

	if err:=c.BodyParser(&patient);err!=nil{
		c.Status(400).JSON(fiber.Map{
			"error":"Invalid request body",
		})
		return err
	}

	if err:=database.DB.Save(&patient).Error;err!=nil{
		c.Status(500).JSON(fiber.Map{
			"Status":"Could not update",
		})
		return err
	}

	c.Status(200).JSON(fiber.Map{
		 "message": "patient updated successfully",
        "data":    patient,
	})
	return nil
}

func DeleteMyAccount(c *fiber.Ctx)error{
	id:=c.Locals("user_id")

	if id==""{
		c.Status(404).JSON(fiber.Map{
			"Message":"id can't be empty",
		})
		return nil
	}

	result:=database.DB.Delete(&models.Patient{},id)
	if result.Error!=nil{
		c.Status(404).JSON(fiber.Map{
			"error":"Patient could not be deleted",
		})
		return result.Error
	}

	c.Status(200).JSON(fiber.Map{
		"message":"Patient has been deleted successfully",
	})

	return nil
} 

func GetAllDoctorBYpatient(c *fiber.Ctx)error{
	doctors:=&[]models.Doctor{}

	err:=database.DB.Find(doctors).Error
	if err!=nil{
		c.Status(404).JSON(fiber.Map{
			"Message":"Can't retrive the doctors",
		})
		return nil
	}

	c.Status(200).JSON(fiber.Map{
 		"Message":"Doctors fetched successfully",
		"Doctors":doctors,
	})
	return nil
}

func GetDoctorByIdbypatient(c *fiber.Ctx)error{
	id:=c.Params("doctorId")
	doctorbyid:=&models.Doctor{}
	if id==""{
		c.Status(404).JSON(fiber.Map{
			"Message":"Id can't be empty",
		})
		return nil
	}

	err:=database.DB.Where("id=?",id).Find(&doctorbyid).Error
	if err!=nil{
		c.Status(400).JSON(fiber.Map{
			"Message":"Could not get the Doctor",
		})
		return err
	}

	c.Status(200).JSON(fiber.Map{
		"Message":"Doctor id fetched successfully",
		"Doctor":doctorbyid,
	})
	return nil
}

func GetDoctorSchedules(c *fiber.Ctx)error{
	id:=c.Params("doctorId")
	if id==""{
		c.Status(400).JSON(fiber.Map{
			"Error":"Id can't be empty",
		})
		return nil
	}

	doctor:=models.Doctor{}

	err:=database.DB.Find(&doctor,id).Error
	if err!=nil{
		return c.Status(400).JSON(fiber.Map{
			"Message":"Could not find doctor",
		})
		
	}

	schedule:=[]models.Schedule{}

	if err:=database.DB.Where("id=?",id).Find(&schedule).Error
	err!=nil{
		return c.Status(404).JSON(fiber.Map{
			"Message":"Failed to fetch the schedule",
		})
	}

	return c.JSON(fiber.Map{
		"doctor":doctor.Name,
		"Schedules":schedule,
	})


}

func CreateAppointment(c *fiber.Ctx)error{
	appointment:=models.Appointment{}

	if err:=c.BodyParser(&appointment);
	err!=nil{
		return c.Status(400).JSON(fiber.Map{
			"Error":"Ivalid request body",
		})
	}

	if err:=database.DB.Create(&appointment).Error;
	err!=nil{
		return c.Status(500).JSON(fiber.Map{
			"error":"Failed to create appointment",
			"details":err.Error(),
		})
	}

	lastQueue:=models.Queue{}
	position:=1

	if err := database.DB.
        Joins("JOIN appointments ON appointments.id = queues.appointment_id").
        Where("appointments.doctor_id = ?", appointment.DoctorId).
        Order("queues.position DESC").
        First(&lastQueue).Error; err == nil {
         position = lastQueue.Position + 1
    }

	 if position == 1 {
        appointment.Status = models.Statusactive
    } else {
        appointment.Status = models.Statusinactive
    }

	if err:=database.DB.Create(&appointment).Error;
	err!=nil{
		c.Status(400).JSON(fiber.Map{
			"Error":"failed to Update the appointment",
		})
	}

	idStr := fmt.Sprintf("%v", c.Locals("user_id"))  
	id, err := strconv.Atoi(idStr)                   
	if err != nil {
    return c.Status(400).JSON(fiber.Map{
        "error": "Invalid user_id in token",
    })
}
	doctorId:=appointment.DoctorId
	queue:=models.Queue{
		AppointmentId:     int(appointment.ID),
        Position:          position,
        Status:          string(appointment.Status), 
		DoctorId:	int(doctorId),
		PatientId:			id,

        EstimatedStartUtc: time.Now().Add(time.Minute * 10),
	}

	if err:=database.DB.Create(&queue).Error;
	err!=nil{
		return c.Status(500).JSON(fiber.Map{
			"Error":"Failed to add the queue",
		})
	}

	 Sendnotification(uint(id), appointment.DoctorId, 
		"Appointment Booked",
		fmt.Sprintf("Your appointment is booked. Queue number: %d", position),
	)

	 return c.JSON(fiber.Map{
		"message":     "Appointment created & added to queue",
		"appointment": appointment,
		"queue":       queue,
	})	 
}


func GetMyAppointments(c *fiber.Ctx)error{
	idStr := fmt.Sprintf("%v", c.Locals("user_id"))  
	id, err := strconv.Atoi(idStr)                   
	if err != nil {
    return c.Status(400).JSON(fiber.Map{
        "error": "Invalid user_id in token",
    })
}
	appointments:=[]models.Appointment{}
	if err:=database.DB.Preload("Doctor").Preload("Queue").Where("id=?",id).Find(&appointments).Error
	err!=nil{
		return c.Status(500).JSON(fiber.Map{
			"Error":"failed to featch appointment",
			"details":err.Error(),
		})
	}
	return c.JSON(appointments)
}

func CancelMyAppointment(c *fiber.Ctx) error {
    appointmentId := c.Params("appointmentId")

    appointment := models.Appointment{}
    if err := database.DB.First(&appointment, appointmentId).Error; err != nil {
        return c.Status(404).JSON(fiber.Map{
            "Error": "Appointment not found",
        })
    } 

    queue := models.Queue{}
    if err := database.DB.Where("appointment_id = ?", appointment.ID).First(&queue).Error; err != nil {
        return c.Status(404).JSON(fiber.Map{
            "error": "Queue not found",
        })
    }

    if err := database.DB.Model(&models.Queue{}).
        Where("position > ?", queue.Position).
        Update("position", gorm.Expr("position - 1")).Error; err != nil {
        return c.Status(500).JSON(fiber.Map{
            "Error":   "Failed to reorder the position",
            "details": err.Error(),
        })
    }

    if err := database.DB.Model(&appointment).Update("status", "cancelled").Error; err != nil {
        return c.Status(500).JSON(fiber.Map{
            "error": "Failed to cancel the appointment",
        })
    }

    return c.JSON(fiber.Map{
        "Message": "Appointment cancelled and queue updated successfully",
    })
}

func GetMyQueueStatus(c *fiber.Ctx)error{
	patientId:=c.Locals("user_id")
	queue:=models.Queue{}
	if err:=database.DB.Where("patient_id=?",patientId).First(&queue).Error;
	err!=nil{
		return c.Status(404).JSON(fiber.Map{
			"Message":"Your not in the queue currently",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"queue_id": queue.ID,
		"position": queue.Position,
	})
}

// func CreateAppointment(c *fiber.Ctx)error{
// 	appointment:=models.Appointment{}

// 	err:=c.BodyParser(&appointment)
// 	if err!=nil{
// 		c.Status(400).JSON(fiber.Map{
// 			"Error":"Invalid Body Request",
// 		})
// 		return err
// 	}

// 	err=database.Create(&appointment).Error
// 	if err!=nil{
// 		c.Status(500).JSON(fiber.Map{
// 			"Message":"Could not able to create the appointment",
// 		})
// 		return err
// 	}

// 	c.Status(200).JSON(fiber.Map{
// 		"Message":"Appointment has been created successfully",
// 		"Appointment":appointment,
// 	})
// 	return nil

// }

// func GetPatientAppointment(c *fiber.Ctx)error{
// 	id=c.ParamsInt(":id")
// 	Appointments:=&[]models.Appointment{}

// 	if id==""{
// 		c.Status(400).JSON(fiber.Map{
// 			"Message":"Id can't be empty"
// 		})
// 		return nil
// 	}

// 	err:=database.Where("id=?",id).Find(Appointments).Error
// 	if err!=nil{
// 		c.Status(400).JSON(fiber.Map{
// 			"Message":"can't able to get the Appointments",
// 		})
// 		return err
// 	}

// 	c.Status(200).JSON(fiber.Map{
// 		"Message":"Patients appointments",
// 		"Appointments":Appointments,
// 	})
// 	return nil
// }

// func Getnotification(c *fiber.Ctx)error{
// 	id=Params(":id")

// 	Notifications:=&[]models.Notification{}

// 	if id==""{
// 		c.Status(400).JSON(fiber.Map{
// 			"Message":"Id can't be empty"
// 		})
// 		return nil
// 	}

// 	err:=database.Where("id=?",id).Find(Notifications).Error
// 	if err!=nil{
// 		c.Status(400).JSON(fiber.Map{
// 			"Message":"can't able to get the Notifications",
// 		})
// 		return err
// 	}

// 	c.Status(200).JSON(fiber.Map{
// 		"Message":"Patients Notifications",
// 		"Notifications":Notifications,
// 	})
// 	return nil
// }



// func GetAllDoctor(c *fiber.Ctx)error{
// 	Doctores:=&[]models.Doctor{}

// 	err:=database.Find(Doctores).Error
// 	if err!=nil{
// 		c.Status(400).JSON(fiber.Map{
// 			"Message":"Can't retrive the books",
// 		})
// 		return err
// 	}

// 	c.Status(200).JSON(fiber.Map{
// 		"Message":"Doctors fetched successfully",
// 		"Doctors":Doctores,
// 	})
// 	return nil
// }


// func GetDoctorById(c *fiber.Ctx)error{
// 	id:=Params(":id")
// 	doctorByid:=&model.GetDoctorById{}

// 	if id==""{
// 		c.Status(404).JSON(fiber.Map{
// 			"Message":"id can't be empty",
// 		})
// 		return nil
// 	}
// 	fmt.Println("the doctor id",id)

// 	err:=database.Where("id=?",id).Find(database).Error
// 	if err!=nil{
// 		c.Status(404).JSON(fiber.Map{
// 			"Message":"Could not get the doctors",
// 		})
// 		return err
// 	}

// 	c.Status(200).JSON(fiber.Map{
// 		"Message":"Doctor id is fetched successfully",
// 		"Doctor":doctorByid,
// 	})
// 	return nil
// }

// func GetDoctorSchedules(c *fiber.Ctx)error{
// 	id:=Params("id")
// 	Schedules:=&[]models.Schedule{}

// 	if id==""{
// 		c.Status(400).JSON(fiber.Map{
// 			"Message":"Id can't be empty"
// 		})
// 		return nil
// 	}

// 	err:=database.Where("id=?",id).Find(Schedules).Error
// 	if err!=nil{
// 		c.Status(400).JSON(fiber.Map{
// 			"Message":"can't able to get the Schedule",
// 		})
// 		return err
// 	}

// 	c.Status(200).JSON(fiber.Map{
// 		"Message":"Doctors Schedule",
// 		"Schedules":Schedules,
// 	})
// 	return nil
// }






package handlers

import (
	"project/database"
	"project/models"
	"project/utils"
	"strings"
	"github.com/gofiber/fiber/v2"
	"fmt"
	"math/rand"
	"time"
)

func DoctorsignUp(c *fiber.Ctx)error{
		var User models.Doctor
		if err:=c.BodyParser(&User); err!=nil{
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})	
		}
		 
	
		if User.Name==""||User.Email==""||User.Password==""{
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "name, email and password are required",
			})
		}

		var existing models.Doctor
		if err:=database.DB.Where("email=?",User.Email).First(&existing)
		 err.Error == nil{
			c.Status(404).JSON(&fiber.Map{
				"message":"email already exists",
			})
			return err.Error
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
			"message":"Doctor created successfully",
			"Doctor":User,
		})
		return nil
}

func DoctorLogin(c *fiber.Ctx)error{   
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

	var user models.Doctor

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

func VerifyOTPDoctor(c *fiber.Ctx) error{

	var payload struct{
		Email string `json:"email"`
		OTP string `json:"otp"`
	}

	if err:=c.BodyParser(&payload);err!=nil{
		return c.Status(400).JSON(fiber.Map{
			"error":"Invalid request",
		})
	}

	var user models.Doctor
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

func GetMyProfileDoctor(c *fiber.Ctx)error{
	id:=c.Locals("user_id")
	DoctorByid:=&models.Doctor{}
	if id ==""{
		c.Status(404).JSON(fiber.Map{
			"Message":"id cannot be empty",
		})
		return nil
	}
	fmt.Println("the Id is",id)

	err:=database.DB.Where("id=?",id).First(DoctorByid).Error
	if err!=nil{         
		c.Status(404).JSON(fiber.Map{
			"message":"could not get the Patients",
		})
		return err
	}

	c.Status(200).JSON(fiber.Map{
		"message":"Patient id fetched successfuly",
		"Name":DoctorByid.Name,
		"Email":DoctorByid.Email,
		"Specialization":DoctorByid.Specialization,
		"AvgconsultMin":DoctorByid.AvgconsultMin,
	})
	return nil
}

func UpdateMyProfileDoctor(c *fiber.Ctx)error{
	id:=c.Locals("user_id")
	var doctor models.Doctor

	if err:=database.DB.First(&doctor,id).Error;err!=nil{
		c.Status(400).JSON(fiber.Map{
			"Message":"Doctor not found",
		})
		return err
	}

	if err:=c.BodyParser(&doctor);err!=nil{
		c.Status(400).JSON(fiber.Map{
			"error":"Invalid request body",
		})
		return err
	}

	if err:=database.DB.Save(&doctor).Error;err!=nil{
		c.Status(500).JSON(fiber.Map{
			"Status":"Could not update",
		})
		return err
	}

	c.Status(200).JSON(fiber.Map{
		 "message": "doctor updated successfully",
        "data":  doctor,
	})
	return nil
}

func DeleteMyAccountDoctor(c *fiber.Ctx)error{
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

func GetMyAppointmentsbydoctor(c *fiber.Ctx)error{
	id:=c.Locals("user_id")

	if id==""{
		c.Status(400).JSON(fiber.Map{
			"error":"id can't be empty",
		})
	}
	appointments:=models.Appointment{}
	if err:=database.DB.Where("id=?",id).Find(&appointments).Error;
	err!=nil{
		c.Status(404).JSON(fiber.Map{
			"Message":"can't fetch appointments",
			"datails":err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"appointments":appointments,
	})
} 

func CompleteAppointment(c *fiber.Ctx)error{
	appointMentId:=c.Params("appointmentId")
	appointments:=models.Appointment{}

	if err:=database.DB.First(&appointments,appointMentId).Error;
	err!=nil{
		return c.Status(404).JSON(fiber.Map{
			"Error":"Appointment not found",
		})
	}
	const StatusCompleted = "completed"
	appointments.Status = StatusCompleted // Ensure the constant is defined in the models package
	if err:=database.DB.Save(&appointments).Error;
	err!=nil{
		c.Status(404).JSON(fiber.Map{
			"Error":"Failed to update appointment",
		})
	}

		if err := database.DB.Model(&models.Queue{}).
		Where("appointment_id = ?", appointments.ID).
		Update("status", "completed").Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to update queue",
		})
	}

	if err := RecalculateQueue(appointments.DoctorId); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to recalculate queue",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Appointment marked as completed and queue updated",
	})
	
}

// func DeleteDoctor(c *fiber.Ctx)error{
// 	id:=Params("id")
// 	if id==""{
// 		c.Status(404).JSON(fiber.Map{
// 			"message":"id can't be empty",
// 		})
// 		return nil
// 	}

// 	result:=database.Delete(&model.Doctor{},id)

// 	if result.Error!=nil{
// 		c.Status(404).JSON(fiber.Map{
// 			"Message":"Doctor can not be deleted",
// 		})
// 		return result.Error
// 	}

// 	c.Status(200).JSON(fiber.Map{
// 		"message":"Doctor has been deleted successfully",
// 	})
// }

// func GetDoctorAppointment(c *fiber.Ctx)error{
// 	id:=Params("id")
// 	appointments:=&[]models.Appointment{}
// 	if id==""{
// 		c.Status(400).JSON(fiber.Map{
// 			"Message":"id can't empty",
// 		})
// 		return nil
// 	}

// 	err:=database.Where("id=?",id).Find(appointments).Error
// 	if err!=nil{
// 		c.Status(404).JSON(fiber.Map{
// 			"Message":"appontments not able to fetch"
// 		})
// 		return err
// 	}

// 	c.Status(200).JSON(fiber.Map{
// 		"Message":"Doctor appointment",
// 		"DoctorAppointment":appointments
// 	})
// }


// func GetDoctorQueue(c *fiber.Ctx)error{
// 	id=Params(":id")

// 	Queues:=&[]models.Queue{}

// 	if id==""{
// 		c.Status(400).JSON(fiber.Map{
// 			"Message":"Id can't be empty"
// 		})
// 		return nil
// 	}

// 	err:=database.Where("id=?",id).Find(Queues).Error
// 	if err!=nil{
// 		c.Status(400).JSON(fiber.Map{
// 			"Message":"can't able to get the Queue",
// 		})
// 		return err
// 	}

// 	c.Status(200).JSON(fiber.Map{
// 		"Message":"Doctors Queue",
// 		"Queue":Queues,
// 	})
// 	return nil
// }

// func UpdateDoctor(c *fiber.Ctx)error{
// 	id:=Params("id")
// 	var doctor models.Doctor

// 	if err:=database.Find(&doctor,id).Error
// 	err!=nil{
// 		c.Status(404).JSON(fiber.Map{
// 			"Message":"Patient not Updated",
// 		})
// 		return err
// 	}

// 	if err:=c.BodyParser(&doctor);err!=nil{
// 		c.Status(500).JSON(fiber.Map{
// 			"Error":"Invalid body request",
// 		})
// 		return err
// 	}

// 	if err:=database.Save(&doctor).Error;
// 	err!=nil{
// 		c.Status(400).JSON(fiber.Map{
// 			"Status":"Could not update",
// 		})
// 		return err
// 	}

// 	c.Status(200).JSON(fiber.Map{
// 		 "message": "Doctor updated successfully",
//         "data":    doctor,
// 	})
// 	return nil
// }

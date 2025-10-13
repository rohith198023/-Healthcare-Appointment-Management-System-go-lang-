package handlers

import (
	"project/database"
	"project/models"
	"strings"
	"github.com/gofiber/fiber/v2"
	"project/utils"
)

func AdminLogin(c *fiber.Ctx)error{
	var credentials struct{
		Email string `json:"email"`
		Password string `json:"Password"`
	}	

	if err:=c.BodyParser(&credentials);err!=nil{
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":"Invalid request Body",
			"details":err.Error(),
		}) 
	}

	if credentials.Email==""||credentials.Password==""{
		return c.Status(404).JSON(fiber.Map{
			"error":"required email or password",
		})
	}

	var user models.Admin
	credentials.Email=strings.TrimSpace(strings.ToLower(credentials.Email))
	credentials.Password=strings.TrimSpace(credentials.Password)

	if err:=database.DB.Where("Email=? and Password=?",credentials.Email,credentials.Password ).First(&user).Error;err!=nil{
		return c.Status(404).JSON(fiber.Map{
			"Message":"Invalid email or password",
		}) 
	}

	token, err := utils.GenrateToken(user.ID, "admin")
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Failed to generate token",
        })
    }

    return c.JSON(fiber.Map{
        "message": "Login successful",
        "token":   token,
    }) 
}

func AdminGetAllPatients(c *fiber.Ctx)error{
		Patients:=&[]models.Patient{}

		err:=database.DB.Find(Patients).Error
		if err!=nil{
			c.Status(400).JSON(fiber.Map{
				"Message":"Can't retrive the Patients",
			})
			return err
		}

		c.Status(200).JSON(fiber.Map{
			"Message":"Patient fetched Successfully",
			"Patients":Patients,
		})
		return nil
}


func GetPatientById(c *fiber.Ctx)error{
	id:=c.Params("patientId")
	patientbyid:=&models.Patient{}

	if id==""{
		c.Status(404).JSON(fiber.Map{
			"Message":"id can't be empty",
		})
		return nil
	}
	

	err:=database.DB.Where("id=?",id).Find(&patientbyid).Error
	if err!=nil{
		c.Status(404).JSON(fiber.Map{
			"Message":"Could not get the patients",
		})
		return err
	}

	c.Status(200).JSON(fiber.Map{
		"Message":"Patient id is fetched successfully",
		"Patient":patientbyid,
	})
	return nil
}

func DeletePatentid(c *fiber.Ctx)error{
	id:=c.Params("patientId")

	if id==""{
		c.Status(404).JSON(fiber.Map{
			"Error":"ID can't empty",
		})
		return nil
	}

	result:=database.DB.Delete(&models.Patient{},id)
	if result.Error!=nil{
		c.Status(400).JSON(fiber.Map{
			"error":"Patient could not be deleted",
		})
		return result.Error
	}

	c.Status(200).JSON(fiber.Map{
		"Message":"Patient has been deletedSuccessfully",
	})

	return nil
}

func GetAllDoctor(c *fiber.Ctx)error{
	doctors:=&[]models.Doctor{}
	err:=database.DB.Find(doctors).Error
	if err!=nil{
		c.Status(400).JSON(fiber.Map{
			"Message":"Can't retrive the doctors",
		})
		return err
	}

	c.Status(200).JSON(fiber.Map{
			"Message":"doctors fetched Successfully",
			"Patients":doctors,
		})
		return nil
}

func GetDoctorById(c *fiber.Ctx)error{
	id:=c.Params("doctorId")
	doctorByid:=&models.Doctor{}

	if id==""{
		c.Status(404).JSON(fiber.Map{
			"Message":"id can't be empty",
		})
		return nil
	}
	

	err:=database.DB.Where("id=?",id).Find(&doctorByid).Error
	if err!=nil{
		c.Status(404).JSON(fiber.Map{
			"Message":"Could not get the doctors",
		})
		return err
	}

	c.Status(200).JSON(fiber.Map{
		"Message":"doctors id is fetched successfully",
		"Doctors":doctorByid,
	})
	return nil
}

func DeleteDoctor(c *fiber.Ctx)error{
	id:=c.Params("doctorId")

	if id==""{
		c.Status(404).JSON(fiber.Map{
			"Error":"ID can't empty",
		})
		return nil
	}

	result:=database.DB.Delete(&models.Doctor{},id)
	if result.Error!=nil{
		c.Status(400).JSON(fiber.Map{
			"error":"Doctor could not be deleted",
		})
		return result.Error
	}

	c.Status(200).JSON(fiber.Map{
		"Message":"=doctor has been deletedSuccessfully",
	})

	return nil
}

func CreateSchedule(c *fiber.Ctx) error {
    doctorId := c.Params("doctorId")

    var doctor models.Doctor
    if err := database.DB.First(&doctor, doctorId).Error; err != nil {
        return c.Status(404).JSON(fiber.Map{
            "error": "Doctor not found",
        })
    }

    var schedule models.Schedule
    if err := c.BodyParser(&schedule); err != nil {
        return c.Status(400).JSON(fiber.Map{
            "error":   "Invalid body request",
            "details": err.Error(),
        })
    }

    // Attach doctor to schedule
    schedule.DoctorID = doctor.ID

    if err := database.DB.Create(&schedule).Error; err != nil {
        return c.Status(500).JSON(fiber.Map{
            "error":   "Unable to create schedule",
            "details": err.Error(),
        })
    }

    return c.Status(201).JSON(fiber.Map{
        "message":  "Schedule created successfully",
        "schedule": schedule,
    })
}


func UpdateSchedule(c *fiber.Ctx)error{
	id:=c.Params("scheduleId")

	if id==""{
		c.Status(400).JSON(fiber.Map{
			"Message":"id can't be empty",
		})
	}

	updatedSchedule:=&models.Schedule{}
	if err:=c.BodyParser(&updatedSchedule);
	err!=nil{
		return c.Status(400).JSON(fiber.Map{
			"Error":"failed body request",
		})
	}

	result:=database.DB.Model(&models.Schedule{}).Where("id=?",id).Updates(updatedSchedule)

	if result.Error!=nil{
		c.Status(400).JSON(fiber.Map{
			"Error":"can't be Updated",
		})
	}

	if result.RowsAffected==0{
		return c.Status(400).JSON(fiber.Map{
			"Message":" not found",
		})
	}


	return c.Status(400).JSON(fiber.Map{
		"Message":"Schedule Updated successfully",
		"Schedule":updatedSchedule,
	})
}


func DeleteSchedule(c *fiber.Ctx)error{
	id:=c.Params("scheduleId")

	if id==""{
		c.Status(404).JSON(fiber.Map{
			"Error":"ID can't empty",
		})
		return nil
	}

	result:=database.DB.Delete(&models.Schedule{},id)
	if result.Error!=nil{
		c.Status(400).JSON(fiber.Map{
			"error":"Schedule could not be deleted",
		})
		return result.Error
	}

	c.Status(200).JSON(fiber.Map{
		"Message":"Schedule has been deletedSuccessfully",
	})

	return nil
}

func GetAllappointments(c *fiber.Ctx)error{
	appointments:=[]models.Appointment{}

	 if err := database.DB.Preload("Doctor").Preload("Patient").Find(&appointments).Error; err != nil {
        return c.Status(500).JSON(fiber.Map{
            "Message": "appointment can't be fetched",
            "details": err.Error(),
        })
    }

	if len(appointments)==0{
		return c.Status(400).JSON(fiber.Map{
			"Message":"No appointment found",
		})
	}
	return c.Status(200).JSON(fiber.Map{
		"Appointments":appointments,
	})
}

func GetAllQueue(c *fiber.Ctx)error{
	queue:=&[]models.Queue{}
	err:=database.DB.Find(queue).Error
	if err!=nil{
		c.Status(400).JSON(fiber.Map{
			"Message":"can't retrive queue",
		})
		return nil
	}

	c.Status(200).JSON(fiber.Map{
		"Message":"queue fetched",
		"Queue":queue,
	})
	return nil
}
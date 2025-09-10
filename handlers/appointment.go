package handlers

// import (
// 	"fmt"
// 	"project/database"
// 	"project/models"
// 	"github.com/gofiber/fiber"
// )



// func GetappointmentById(c *fiber.Ctx)error{
// 	id:=c.Params("id")
// 	appointmentById:=&models.Appointment{}

// 	if id==""{
// 		c.Status(404).JSON(fiber.Map{
// 			"Message":"id can't be empty",
// 		})
// 		return nil
// 	}

// 	fmt.Println("The id is",id)

// 	err:=database.Where("id=?",id).Find(appointmentById).Error
// 	if err!=nil{
// 		c.Status(404).JSON(fiber.Map{
// 			"Message":"Could nit get  the appointment",
// 		})
// 		return err
// 	}

// 	c.Status(200).JSON(fiber.Map{
// 		"Message":"appointment id is fetched successfully",
// 		"appointment":appointmentById,
// 	})
// 	return nil
// }

// func Deleteappointment(c *fiber.Ctx)error{
// 	id:=c.Params("id")

// 	if id==""{
// 		c.Status(404).JSON(fiber.Map{
// 			"message":"id can't be empty",
// 		})
// 		return nil
// 	}
// 	result:=database.Delete(&models.Appointment{},id)
// 	if result.Error!=nil{
// 		c.Status(404).JSON(fiber.Map{
// 			"Message":"Appointment can;t be cancled",
// 		})
// 		return result.Error
// 	}

// 	c.Status(200).JSON(fiber.Map{
// 		"message":"Appointment is cancled successfully",
// 	})

// 	return nil
// }

// func GetAllappointments(c *fiber.Ctx)error{
// 	Appointments:=&[]models.Appointment{}

// 	err:=database.Find(Appointments).Error
// 	if err!=nil{
// 		c.Status(404).JSON(fiber.Map{
// 			"Message":"Can't fetch the appointment",
// 		})
// 		return err
// 	}

// 	c.Status(200).JSON(fiber.Map{
// 		"Message":"Appointment is fetched successfullty",
// 		"Appointments":Appointments,
// 	})
// 	return nil
// }


// func Updateappointment(c *fiber.Ctx)error{
// 	id:=c.Params("id")
// 	var appointment models.Appointment

// 	if err:=database.Find(&appointment,id).Error
// 	err!=nil{
// 		c.Status(400).JSON(fiber.Map{
// 			"Message":"appointment not found",
// 		})
// 		return err
// 	}

// 	if err:=c.BodyParse(&appointment);
// 	err!=nil{
// 		c.Status(404).JSON(fiber.Map{
// 			"error":"Invalid body request",
// 		})
// 		return err
// 	}

// 	if err:=database.Save(&appointment).Error;
// 	err!=nil{
// 		c.Status(404).JSON(fiber.Map{
// 			"Status":"Could not Update",
// 		})
// 		return err
// 	}

// 	c.Status(200).JSON(fiber.Map{
// 		"Message":"appointment updated successfully",
// 		"Appointment":appointment,
// 	})
// 	return nil
// }



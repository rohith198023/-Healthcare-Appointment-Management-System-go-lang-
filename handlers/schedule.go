package handlers

// import (
// 	"fmt"
// 	"project/database"
// 	"project/models"

// 	"github.com/gofiber/fiber"
// )


// func GetScheduleById(c *fiber.Ctx)error{
// 	id:=c.Params("id")
// 	SchedulerById:=&models.Schedule{}
// 	if id==""{
// 		c.Status(404).JSON(fiber.Map{
// 			"Message":"Id can't be empty",
// 		})
// 		return nil
// 	}
// 	fmt.Println("The id is",id)

// 	err:=database.Where("id=?",id).Find(&SchedulerById).Error
// 	if err!=nil{
// 		c.Status(404).JSON(fiber.Map{
// 			"Message":"Could not get the patients",
// 		})
// 		return err
// 	}

// 	c.Status(200).JSON(fiber.Map{
// 		"Message":"Scheduler id fetched successfully",
// 		"Scheduler":SchedulerById
// 	})
// }

// func DeleteSchedule(c *fiber.Ctx)error{
// 	id:=Params("id")

// 	if id==""{
// 		c.Status(404).JSON(fiber.Map{
// 			"message":"id can't be empty",
// 		})
// 		return nil
// 	}

// 	result:=database.Delete(&model.Schedule{},id)
// 	if result.Error!=nil{
// 		c.Status(404).JSON(fiber.Map{
// 			"message":"Schedule can't be deleted",
// 		})
// 		return nil
// 	}

// 	c.Status(200).JSON(fiber.Map{
// 		"message":"Schedule has been deleted successfully",
// 	})
// 	return nil

// }

// func CreateSchedule(c *fiber.Ctx)error{
// 	Schedule:=models.Schedule{}

// 	err:=c.BodyParser(&Schedule)
// 	if err!=nil{
// 		c.Status(400).JSON(fiber.Map{
// 			"Error":"Invalid body request",
// 		})
// 		return err
// 	}

// 	err=database.Create(&Schedule).Error
// 	if err!=nil{
// 		c.Status(500).JSON(fiber.Map{
// 			"Message":"Could not Create Schedule",
// 		})
// 		return err
// 	}

// 	c.Status(200).JSON(fiber.Map{
// 		"Message":"Schedule hs been created",
// 		"Schedule":Schedule
// 	})
// 	return nil
// }

// func UpdateSchedule(c *fiber.Ctx)error{
// 	id:=Params("id")
// 	var schedule models.Schedule

// 	if err:=database.First(&schedule,id).Error;
// 	err!=nil{
// 		c.Status(400).JSON(fiber.Map{
// 			"Message":"Patient not found",
// 		})
// 		return err
// 	}

// 	if err:=c.BodyParse(&schedule);err!=nil{
// 		c.Status(400).JSON(fiber.Map{
// 			"error":"Invalid Body request",
// 		})
// 		return err
// 	}

// 	if err:=database.Save(&schedule).Error
// 	if err!=nil{
// 		c.Status(404).JSON(fiber.Map{
// 			"Status":"could not able to update",
// 		})
// 		return err
// 	}

// 	c.Status(200).JSON(fiber.Map{
// 		"message":"Schedule has been Updated Successfully",
// 		"Schedule":schedule,
// 	})

// 	return nil
// }

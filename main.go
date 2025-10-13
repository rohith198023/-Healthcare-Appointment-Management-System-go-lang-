package main

import (
	// "os"
	"project/database"
	"project/handlers"
	"github.com/gofiber/fiber/v2"
	"project/middleware"
)

func main(){
	database.DBconnect()
    app :=fiber.New() 

	// ================== AUTH ==================
	app.Post("/auth/patient/signup", handlers.PatientsignUp)   
	app.Post("/auth/patient/login", handlers.PatientLogin)
	app.Post("/auth/verify-otp/patient", handlers.VerifyOTPPatient)


	app.Post("/auth/doctor/signup", handlers.DoctorsignUp)     // doctor register
	app.Post("/auth/doctor/login", handlers.DoctorLogin)  
	app.Post("/auth/verify-otp/doctor", handlers.VerifyOTPDoctor)
	//admin login
	app.Post("/auth/admin/login", handlers.AdminLogin)

	//patient api's 
	PatientApi := app.Group("/api/patient", middleware.VerifyToken)
	PatientApi.Get("/profile", handlers.GetMyProfile)
	PatientApi.Put("/profile/update", handlers.UpdateMyProfile) 
	PatientApi.Delete("/profile/delete", handlers.DeleteMyAccount)

	
	PatientApi.Get("/doctors", handlers.GetAllDoctorBYpatient)   
	PatientApi.Get("/doctors/:doctorId", handlers.GetDoctorByIdbypatient)
	PatientApi.Get("/doctors/:doctorId/schedules", handlers.GetDoctorSchedules)


	PatientApi.Post("/appointments/post", handlers.CreateAppointment)    
	PatientApi.Get("/appointments/get", handlers.GetMyAppointments)    
	PatientApi.Delete("/appointments/:appointmentId", handlers.CancelMyAppointment) 

	PatientApi.Get("/queue", handlers.GetMyQueueStatus)




	
	//Doctor api's
	DoctorApi := app.Group("/api/doctor", middleware.VerifyToken)
	DoctorApi.Get("/profile", handlers.GetMyProfileDoctor)
	DoctorApi.Put("/profile/update", handlers.UpdateMyProfileDoctor)       
	DoctorApi.Delete("/profile/delete", handlers.DeleteMyAccountDoctor)
	DoctorApi.Get("/appointments/get", handlers.GetMyAppointmentsbydoctor) 

	DoctorApi.Put("/appointments/:appointmentId/complete", handlers.CompleteAppointment)

	//admin api's
	AdminApi := app.Group("/api/admin", middleware.VerifyToken)

	//manage patients
	AdminApi.Get("/patients", handlers.AdminGetAllPatients)  
	AdminApi.Get("/patients/get/:patientId", handlers.GetPatientById)  
	AdminApi.Delete("/patients/delete/:patientId", handlers.DeletePatentid) 
	
	// Manage doctors
	AdminApi.Get("/doctors", handlers.GetAllDoctor) 
	AdminApi.Get("/doctors/get/:doctorId", handlers.GetDoctorById) 
	AdminApi.Delete("/doctors/delete/:doctorId", handlers.DeleteDoctor)  

	AdminApi.Post("/doctors/:doctorId/schedules", handlers.CreateSchedule) 
	AdminApi.Put("/schedules/:scheduleId", handlers.UpdateSchedule)   
	AdminApi.Delete("/schedules/:scheduleId", handlers.DeleteSchedule)

	AdminApi.Get("/appointments", handlers.GetAllappointments)

	AdminApi.Get("/queues", handlers.GetAllQueue) 
	// port:=os.Getenv("port")	

	// AdminApi.Get("/queues", handlers.GetAllQueue)   
	app.Listen(":5678")
	
}

















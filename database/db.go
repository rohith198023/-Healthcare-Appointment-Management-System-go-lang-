package database

import (
	"log"
	"os"
	"project/models"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func DBconnect() {
	var err error
	err=godotenv.Load()           
	if err!=nil{
		log.Fatal("Failed to load env")
	}
	databaseurl:=os.Getenv("dns")
    dsn :=databaseurl						
    db, err := gorm.Open(postgres.New(postgres.Config{
    DSN: dsn,                                                     
    PreferSimpleProtocol: true, // disables implicit prepared statement usage
}), &gorm.Config{})

    if err != nil {
        log.Fatalf("Unable to connect to database: %v\n", err)
    }
	DB=db 
    // Auto migrate the schema
    // DB.AutoMigrate(&models.Appointment{},&models.Doctor{},&models.Notification{},&models.Patient{},&models.Queue{},&models.Schedule{})

	err = DB.AutoMigrate(
    &models.Appointment{},
    &models.Doctor{},
    &models.Notification{},
    &models.Patient{},
    &models.Queue{},
    &models.Schedule{},
	&models.Admin{},
	)
	if err != nil {
    log.Fatalf("Migration failed: %v", err)
	}

	DB.Create(&models.Admin{
		Name:"Pavan",
		Email:"chavapavankumar1234@gmail.com",
		Password:"pavango",
		Role:"admin",
	})
}








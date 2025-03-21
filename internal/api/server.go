package api

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/reminders/config"
	remindercontroller "github.com/reminders/controllers/reminderController"
	usercontrollers "github.com/reminders/controllers/userControllers"
	"github.com/reminders/middlewares"
	"github.com/reminders/models"
	sendReminder "github.com/reminders/send-reminder"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func StartServer(config config.AppConfig) {
	app := gin.Default()

	db, err := gorm.Open(postgres.Open(config.Dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error in connecting DB: %v", err.Error())
	}

	SetupRoutes(app, db, config)

	err = db.AutoMigrate(&models.User{})
	models.MigrateDB(db)
	models.MigrateReminder(db)

	if err != nil {
		log.Fatalf("Error in migrating: %v", err)
	}

	go sendReminder.RunCron(db)

	app.Run(fmt.Sprintf(":%v", config.ServerPort))
}

func SetupRoutes(app *gin.Engine, db *gorm.DB, config config.AppConfig) {
	userRoutes := app.Group("/users")
	usercontrollers.SetUpUserRoutes(userRoutes, db, config.Secret)

	middleware := middlewares.SetUpMiddleware(db, config.Secret)
	reminderRoutes := app.Group("/reminder", middleware.UserAuth)
	remindercontroller.SetupReminderRoutes(reminderRoutes, db, config.Secret)
}

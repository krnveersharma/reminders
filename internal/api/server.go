package api

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-contrib/cors"
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
	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	SetupRoutes(app, db, config)

	err = db.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatalf("Error in migrating: %v", err)
	}
	models.MigrateDB(db)
	models.MigrateReminder(db)
	err = db.AutoMigrate(&models.Draft{})

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

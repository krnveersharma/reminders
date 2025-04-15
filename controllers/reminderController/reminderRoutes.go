package remindercontroller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/reminders/internal/dto"
	"github.com/reminders/models"
	"gorm.io/gorm"
)

type ReminderRoutes struct {
	DB     *gorm.DB
	Secret string
}

func SetupReminderRoutes(router *gin.RouterGroup, db *gorm.DB, secret string) {
	routes := &ReminderRoutes{
		DB:     db,
		Secret: secret,
	}
	router.POST("/add-reminder", routes.addReminder)
	router.POST("/add-draft", routes.saveDraft)
	router.GET("/get-drafts", routes.getDrafts)
}

func (r *ReminderRoutes) addReminder(ctx *gin.Context) {
	userVal, found := ctx.Get("user")
	if !found {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Please login first",
		})
		return
	}

	user, ok := userVal.(models.User)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid user session",
		})
		return
	}

	fmt.Printf("user id is:%d\n", user.ID)

	var reminderData dto.Reminder
	if err := ctx.ShouldBindBodyWithJSON(&reminderData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Please enter correct information",
		})
		return
	}

	date := reminderData.Date + "T" + reminderData.Time + ":00"
	query := "INSERT INTO reminders(user_id,reciever_info,priority,data,data_type,reminder_type,date) VALUES(?,?,?,?,?,?,?)"
	result := r.DB.Exec(query, user.ID, reminderData.RecieverInfo, reminderData.Priority, reminderData.Data, reminderData.DataType, reminderData.ReminderType, date)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": result.Error,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Added reminder Successfuly",
	})
	return
}

func (r *ReminderRoutes) saveDraft(ctx *gin.Context) {
	userVal, found := ctx.Get("user")
	if !found {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Please login first",
		})
		return
	}

	user, ok := userVal.(models.User)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid user session",
		})
		return
	}
	var draftData dto.Draft
	if err := ctx.ShouldBindBodyWithJSON(&draftData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Please enter correct information",
		})
		return
	}

	query := "INSERT INTO drafts(user_id,data,data_type,reminder_type) VALUES(?,?,?,?)"
	result := r.DB.Exec(query, user.ID, draftData.Data, draftData.DataType, draftData.ReminderType)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": result.Error,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Added reminder Successfuly",
	})
	return
}

func (r *ReminderRoutes) getDrafts(ctx *gin.Context) {
	var data []dto.Draft

	userVal, found := ctx.Get("user")
	if !found {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "please login first",
		})
	}

	user, ok := userVal.(models.User)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid user session",
		})
		return
	}

	if err := r.DB.Where("user_id = ?", user.ID).Find(&data).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"drafts": data,
	})
}

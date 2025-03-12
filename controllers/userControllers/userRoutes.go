package usercontrollers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/reminders/internal/domain"
	"gorm.io/gorm"
)

type userRoutes struct {
	db *gorm.DB
}

func SetUpUserRoutes(router *gin.RouterGroup, db *gorm.DB) {
	routes := &userRoutes{db: db}
	router.POST("/register", routes.RegisterUser)
}

func (r *userRoutes) RegisterUser(ctx *gin.Context) {
	var user domain.User
	if err := ctx.ShouldBindBodyWithJSON(&user); err != nil {
		fmt.Errorf("give correct user data")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user data"})
		return
	}
	err := Register(user, r.db)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "User registered successfully", "user": user})
}

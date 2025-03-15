package usercontrollers

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/reminders/internal/dto"
	"github.com/reminders/middlewares"
	"github.com/reminders/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type userRoutes struct {
	db     *gorm.DB
	Secret string
}

func SetUpUserRoutes(router *gin.RouterGroup, db *gorm.DB, secret string) {
	routes := &userRoutes{
		db:     db,
		Secret: secret,
	}

	setupMiddleware := middlewares.SetUpMiddleware(db, secret)

	privateRoutes := router.Group("/verify", setupMiddleware.UserAuth)
	router.POST("/register", routes.RegisterUser)
	router.POST("/login", routes.login)
	privateRoutes.PUT("/edit-user", routes.editUser)
}

func (r *userRoutes) RegisterUser(ctx *gin.Context) {
	var user models.User
	if err := ctx.ShouldBindBodyWithJSON(&user); err != nil {
		fmt.Errorf("give correct user data")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user data"})
		return
	}
	err := Register(user, r.db)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "User registered successfully", "user": user})
}

func (r *userRoutes) login(ctx *gin.Context) {
	var user dto.UserLogin
	if err := ctx.ShouldBindBodyWithJSON(&user); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": "Invalid user data"})
		return
	}

	getUser, err := getUser(user.Email, r.db)
	fmt.Printf("user is:", getUser)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(getUser.Password), []byte(user.Password)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	jwtToken, err := createToken(getUser, r.Secret)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": jwtToken})
}

func (r *userRoutes) editUser(ctx *gin.Context) {
	var user models.User
	if err := ctx.ShouldBindBodyWithJSON(&user); err != nil {
		fmt.Errorf("user information wrong: %s", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "please provide correct information",
		})
		return
	}

	err := editUserInfo(user, r.db)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusAccepted, gin.H{
		"mesage": "Edited Successfully",
	})
	return
}

func createToken(user models.User, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"id":    user.ID,
			"email": user.Email,
			"exp":   time.Now().Add(time.Hour * 24).Unix(),
		})

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		fmt.Errorf("Error in jwt token parsing: %v", err)
		return "", errors.New("some error occured")
	}

	return tokenString, nil
}

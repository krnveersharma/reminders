package usercontrollers

import (
	"errors"
	"fmt"
	"time"

	"github.com/reminders/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Register(user models.User, db *gorm.DB) error {
	existingUser, err := getUser(user.Email, db)

	if existingUser.Email != "" {
		fmt.Printf("Register:user already exist\n")
		return errors.New("user already exists")
	}
	if err != nil {
		return err
	}

	user.Password, err = getHashedPassword(user.Password)
	if err != nil {
		return err
	}

	insertQuery := "INSERT INTO users (first_name, last_name, email, password) VALUES (?, ?, ?, ?)"
	result := db.Exec(insertQuery, user.FirstName, user.LastName, user.Email, user.Password)

	if result.Error != nil {
		fmt.Printf("Error inserting user: %v\n", result.Error)
		return errors.New("unable to register user")
	}

	fmt.Println("User registered successfully:", user.Email)
	return nil
}
func getHashedPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		fmt.Printf("Error hashing password: %s\n", err.Error())
		return "", errors.New("unable to generate hash of password")
	}
	return string(hashedPassword), nil
}

func getUser(email string, db *gorm.DB) (models.User, error) {
	var existingUser models.User

	query := "SELECT * FROM users WHERE email = ?"
	result := db.Raw(query, email).Scan(&existingUser)

	if result.Error != nil {
		fmt.Errorf("Error in executing sql query:", result.Error)
		return models.User{}, errors.New("Some error occured")
	}

	return existingUser, nil
}

func EditUserInfo(user models.User, db *gorm.DB) error {
	existingUser, err := getUser(user.Email, db)
	if err != nil {
		fmt.Errorf("error in finding user: %s", err.Error())
		return errors.New("No user found")
	}

	if user.PlanID != existingUser.PlanID {
		existingUser.UpdatedAt = time.Now()
	}
	existingUser.FirstName = user.FirstName
	existingUser.LastName = user.LastName
	existingUser.Email = user.Email
	existingUser.Phone = user.Phone
	existingUser.AppPassword = user.AppPassword
	existingUser.WhatsappNumber = user.WhatsappNumber
	existingUser.PlanID = user.PlanID
	existingUser.UpdatedAt = time.Now()
	existingUser.LowRemindersUtilized = user.LowRemindersUtilized
	existingUser.MediumRemindersUtilized = user.MediumRemindersUtilized
	existingUser.HighRemindersUtilized = user.HighRemindersUtilized

	result := db.Save(&existingUser)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

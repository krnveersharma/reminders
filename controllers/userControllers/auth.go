package usercontrollers

import (
	"errors"
	"fmt"

	"github.com/reminders/internal/domain"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Register(user domain.User, db *gorm.DB) error {
	existingUser, err := getUser(user.Email, db)

	if existingUser.Email != "" {
		fmt.Printf("Register:user already exist\n")
		return errors.New("user already exists")
	}
	if err != nil {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil {
		fmt.Printf("Error hashing password: %s\n", err.Error())
		return errors.New("unable to generate hash of password")
	}
	user.Password = string(hashedPassword)

	insertQuery := "INSERT INTO users (first_name, last_name, email, password) VALUES (?, ?, ?, ?)"
	result := db.Exec(insertQuery, user.FirstName, user.LastName, user.Email, user.Password)

	if result.Error != nil {
		fmt.Printf("Error inserting user: %v\n", result.Error)
		return errors.New("unable to register user")
	}

	fmt.Println("User registered successfully:", user.Email)
	return nil
}

func getUser(email string, db *gorm.DB) (domain.User, error) {
	var existingUser domain.User

	query := "SELECT * FROM users WHERE email = ?"
	result := db.Raw(query, email).Scan(&existingUser)

	if result.Error != nil {
		fmt.Errorf("Error in executing sql query:", result.Error)
		return domain.User{}, errors.New("Some error occured")
	}

	return existingUser, nil
}

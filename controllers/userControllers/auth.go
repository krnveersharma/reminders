package usercontrollers

import (
	"errors"
	"fmt"

	"github.com/reminders/internal/domain"
	"gorm.io/gorm"
)

func Register(user domain.User, db *gorm.DB) error {
	registeredUser := db.Create(&user)
	if registeredUser.Error != nil {
		fmt.Errorf("Error in creating user: %v", registeredUser.Error)
		return errors.New("error in registering user")
	}

	return nil
}

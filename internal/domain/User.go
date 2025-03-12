package domain

import "time"

type User struct {
	ID             uint      `json:"id" gorm:"PrimaryKey"`
	FirstName      string    `json:"first_name" binding:"required"`
	LastName       string    `json:"last_name"`
	Email          string    `json:"email"`
	Phone          string    `json:"phone"`
	AppPassword    string    `json:"app_password"`
	WhatsappNumber string    `json:"whatsapp_number"`
	Password       string    `json:"password"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

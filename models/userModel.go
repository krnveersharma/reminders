package models

import "time"

type User struct {
	ID             uint      `json:"id" gorm:"PrimaryKey"`
	FirstName      string    `json:"first_name" binding:"required"`
	LastName       string    `json:"last_name"`
	Email          string    `json:"email" binding:"required" gorm:"unique"`
	Phone          string    `json:"phone"`
	AppPassword    string    `json:"app_password"`
	WhatsappNumber string    `json:"whatsapp_number"`
	Password       string    `json:"password" binding:"required"`
	PlanID         uint      `json:"plan_id" gorm:"not null;default:1"`
	Plan           Plan      `json:"plan" gorm:"foreignKey:PlanID;constraint:OnUpdate:CASCADE,OnDelete:SET DEFAULT;"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

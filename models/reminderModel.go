package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type ReminderType string

const (
	Whatsapp ReminderType = "whatsapp"
	Email    ReminderType = "email"
	Sms      ReminderType = "sms"
)

type Reminder struct {
	ID           int          `gorm:"primaryKey"`
	UserId       *int         `json:"user_id" gorm:"default:null"`
	User         User         `json:"user" gorm:"foreignKey:UserId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	RecieverInfo string       `json:"reciever_info" gorm:"not null"`
	Priority     string       `json:"priority" gorm:"not null"`
	Data         string       `json:"data" gorm:"not null"`
	DataType     string       `json:"data_type" gorm:"not null"`
	ReminderType ReminderType `json:"reminder_type" gorm:"type:reminder_type;not null"`
	Date         time.Time    `json:"date" gorm:"type:time;not null"`
}

func MigrateReminder(db *gorm.DB) error {
	err := db.Exec("DO $$ BEGIN IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'reminder_type') THEN CREATE TYPE reminder_type AS ENUM ('whatsapp', 'email', 'sms'); END IF; END $$;").Error
	if err != nil {
		fmt.Println("ENUM reminder_type already exists or couldn't be created:", err)
	}

	return db.AutoMigrate(&Reminder{})
}

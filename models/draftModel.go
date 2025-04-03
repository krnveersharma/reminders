package models

type Draft struct {
	ID           uint   `gorm:"PRimaryKey"`
	UserId       *int   `json:"user_id" gorm:"default:null"`
	User         User   `json:"user" gorm:"foreignKey:UserId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Data         string `json:"data" gorm:"not null"`
	DataType     string `json:"data_type" gorm:"not null"`
	ReminderType string `json:"reminder_type" gorm:"not null"`
}

package dto

type Draft struct {
	ID           uint   `gorm:"PRimaryKey"`
	Data         string `json:"data" gorm:"not null"`
	DataType     string `json:"data_type" gorm:"not null"`
	ReminderType string `json:"reminder_type" gorm:"not null"`
}

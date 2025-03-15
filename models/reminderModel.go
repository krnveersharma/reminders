package models

type Reminder struct {
	ID       int    `gorm:"primaryKey"`
	UserId   int    `json:"user_id" gorm:"not null"`
	User     User   `json:"user" gorm:"foreginKey:UserId;constraint:OnUpdate:CASCADE,OnDelete:SET DEFAULT;"`
	Priority string `json:"priority" gorm:"not null"`
	Data     string `json:"data" gorm:"not null"`
	DataType string `json:"data_type" gorm:"not null"`
}

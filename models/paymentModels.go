package models

import "time"

type Order struct {
	ID        uint      `gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	OrderId   uint      `json:"order_id" gorm:"uniqueIndex;autoIncrement"`
	PaymentId string    `json:"payment_id"`
	Signature string    `json:"signature"`
	Amount    int64     `json:"amount"`
	Currency  string    `json:"currency"`
	Status    string    `json:"status"`
	UserID    uint      `json:"user_id"`
	Usser     User      `gorm:"foreignKey:UserID" json:"user"`
}

package utils

import (
	"errors"
	"fmt"

	"github.com/reminders/internal/service"
	"github.com/reminders/models"
	"gorm.io/gorm"
)

func GetPlanDetails(db *gorm.DB) (*models.Plan, error) {
	var plans models.Plan

	query := "SELECT * from plans"
	if err := db.Raw(query).Scan(&plans).Error; err != nil {
		fmt.Printf("error while getting plans from db:", err)
		return nil, errors.New("unable to get plan details from db")
	}
	service.SetRedisKey("plans", plans)
	return &plans, nil
}

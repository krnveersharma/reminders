package utils

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/reminders/internal/service"
	"github.com/reminders/models"
	"gorm.io/gorm"
)

func GetPlanDetails(db *gorm.DB) (*models.Plan, error) {
	var plans models.Plan

	val, err := service.GetFromRedis("plans")
	fmt.Printf("val is:%v", val)
	if err == nil && val != "" {
		err = json.Unmarshal([]byte(val), &plans)
		if err == nil {
			fmt.Println("plans are:%v", plans)
			return &plans, nil
		}
	}

	fmt.Println("plans are:%+v\n", plans)

	query := "SELECT * from plans"
	if err := db.Raw(query).Scan(&plans).Error; err != nil {
		fmt.Printf("error while getting plans from db:", err)
		return nil, errors.New("unable to get plan details from db")
	}

	service.SetRedisKey("plans", plans)
	return &plans, nil
}

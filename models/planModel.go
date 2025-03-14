package models

import (
	"log"

	"gorm.io/gorm"
)

type PlanType string

const (
	FreePlan         PlanType = "free"
	DeltaPlan        PlanType = "delta"
	DeltaPlusPlan    PlanType = "delta_plus"
	DeltaPremiumPlan PlanType = "delta_premium"
)

func (PlanType) GormDataType() string {
	return "plan_type"
}

type Plan struct {
	ID              uint     `gorm:"primaryKey"`
	PlanType        PlanType `gorm:"type:plan_type;not null;default:'free'"` // Uses PostgreSQL ENUM
	Price           string
	LowReminders    uint8
	MediumReminders uint8
	HighReminders   uint8
}

func MigrateDB(db *gorm.DB) {
	err := db.Exec(`
		DO $$ 
		BEGIN 
			IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'plan_type') THEN 
				CREATE TYPE plan_type AS ENUM ('free', 'delta', 'delta_plus', 'delta_premium'); 
			END IF; 
		END $$;
	`).Error
	if err != nil {
		log.Fatalf("Error creating ENUM: %v", err)
	}

	// Migrate the Plan table
	err = db.AutoMigrate(&Plan{})
	if err != nil {
		log.Fatalf("Error in migrating: %v", err)
	}
}

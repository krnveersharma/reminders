package sendReminder

import (
	"fmt"
	"sort"
	"time"

	"github.com/reminders/models"
	"github.com/robfig/cron"
	"gorm.io/gorm"
)

func RunCron(db *gorm.DB) {
	c := cron.New()
	fmt.Println("cron job started")

	err := c.AddFunc("@every 1m", func() {
		fmt.Printf("function execution in cron started")
		getDataFromDB(db)
	})

	if err != nil {
		fmt.Printf("Error scheduling cron job: %v\n", err)
		return
	}

	c.Start()
}

func getDataFromDB(db *gorm.DB) {
	var reminders []models.Reminder
	nowTime := time.Now()
	newTime := nowTime.Add(time.Minute).Truncate(time.Minute)
	result := db.Preload("User").Where("date = ?", newTime).Find(&reminders)
	if result.Error != nil {
		fmt.Printf("unable to run query or parse result: %s", result.Error)
		return
	}

	sort.Slice(reminders, func(i, j int) bool {
		priorityOrder := map[string]int{"high": 1, "medium": 2, "low": 3}
		return priorityOrder[reminders[i].Priority] > priorityOrder[reminders[j].Priority]
	})

	for i := 0; i < len(reminders); i++ {
		if reminders[i].ReminderType == "email" {
			fmt.Printf("user is:%+v\n", reminders[i].User)
			sendUsingSMTP(reminders[i].User.Email, reminders[i].RecieverInfo, reminders[i].User.AppPassword, reminders[i].Data, reminders[i].DataType)
		}
	}
}

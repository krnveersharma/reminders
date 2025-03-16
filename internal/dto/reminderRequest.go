package dto

type ReminderType string

const (
	Whatsapp ReminderType = "whatsapp"
	Email    ReminderType = "email"
	Sms      ReminderType = "sms"
)

type Reminder struct {
	RecieverInfo string       `json:"reciever_info"`
	Priority     string       `json:"priority" `
	Data         string       `json:"data" `
	DataType     string       `json:"data_type" `
	ReminderType ReminderType `json:"reminder_type"`
	Date         string       `json:"date" `
	Time         string       `json:"time"`
}

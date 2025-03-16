package sendReminder

import (
	"fmt"
	"net/smtp"
)

func sendUsingSMTP(from, to, appPassword, data string) {

	fmt.Printf("from: %s to: %s appPassword: %s\n", from, to, appPassword)

	host := "smtp.gmail.com"
	port := "587"

	auth := smtp.PlainAuth("", from, appPassword, host)
	err := smtp.SendMail(host+":"+port, auth, from, []string{to}, []byte(data))

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Successfuly sent email to all users\n")
}

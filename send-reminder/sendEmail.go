package sendReminder

import (
	"fmt"
	"net/smtp"
)

func sendUsingSMTP(from, to, appPassword, data, dataType string) {
	host := "smtp.gmail.com"
	port := "587"
	var message string
	auth := smtp.PlainAuth("", from, appPassword, host)

	fmt.Printf("data type is:", dataType)

	message = "Subject: Your Plain Text Email\r\n" +
		"\r\n" + data
	if dataType == "html" {
		message = "MIME-Version: 1.0\r\n" +
			"Content-Type: text/html; charset=\"UTF-8\"\r\n" +
			"Subject: Your HTML Email\r\n" +
			"\r\n" + data
	}
	err := smtp.SendMail(host+":"+port, auth, from, []string{to}, []byte(message))

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Successfuly sent email to all users\n")
}

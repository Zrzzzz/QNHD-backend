package email

import (
	"log"
	"net/smtp"
	"qnhd/pkg/setting"
)

var Code map[string]string

func init() {
	Code = make(map[string]string)
}

func SendEmail(to string, success func(), failure func()) {
	from := setting.OfficeEmail
	pass := setting.OfficePass

	msg := "Hello world"

	err := smtp.SendMail("smtp.163.com",
		smtp.PlainAuth("", from, pass, "smtp.163.com"),
		from, []string{to}, []byte(msg))

	if err != nil {
		log.Printf("Failed to send email: %v", err)
		failure()
	}
	log.Printf("Successfully send email.")
	success()
}

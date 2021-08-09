package email

import (
	"fmt"
	"net/smtp"
	"qnhd/pkg/logging"
	"qnhd/pkg/setting"
	"time"

	"github.com/jordan-wright/email"
)

var Code map[string]int64

func init() {
	Code = make(map[string]int64)
}

func SendEmail(to string, success func(), failure func()) {
	// TODO: 连接池和计时器清除Code保留

	from := setting.OfficeEmail
	pass := setting.OfficePass
	server := setting.EmailSmtp
	port := setting.EmailPort

	t := time.Now().Unix()
	code := fmt.Sprintf("%06d", t%1000000)
	Code[to] = t
	msg := fmt.Sprintf("验证码: %v", code)

	mail := email.NewEmail()
	mail.From = from
	mail.To = []string{to}
	mail.Subject = "青年湖底用户注册"
	mail.Text = []byte(msg)
	err := mail.Send(fmt.Sprintf("%v:%v", server, port), smtp.PlainAuth("", from, pass, server))

	if err != nil {
		logging.Error("Failed to send email: %v", err)
		failure()
	} else {
		success()
	}
}

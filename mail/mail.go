package mail

import (
	"gopkg.in/gomail.v2"
	"errors"
)

type SMTP struct {
	Server string
	Port int
	User string
	Password string
}

func SendMail(subject, body string, attach, to []string, smtpServer SMTP) (bool, error) {
	if smtpServer.Server == "" || smtpServer.Port == 0 || smtpServer.User == "" || smtpServer.Password == "" {
		err := errors.New("Check the [SMTPServer|SMTPPort|SMTPUser|SMTPPassword] environment variables")
		return false, err
	}

	m := gomail.NewMessage()
	m.SetHeader("From", smtpServer.User)
	m.SetHeader("To", to...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)
	for _, v := range attach {
		m.Attach(v)
	}

	d := gomail.NewDialer(smtpServer.Server, smtpServer.Port, smtpServer.User, smtpServer.Password)

	if err := d.DialAndSend(m); err != nil {
		return false, err
	}

	return true, nil
}
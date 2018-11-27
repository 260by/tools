package mail

import (
	"errors"
	"gopkg.in/gomail.v2"
)

// Server 邮件服务器配置信息
type Server struct {
	Addr     string
	Port     int
	User     string
	Password string
}

// Send 发送普通文本邮件
func (m *Server) Send(subject, body string, to []string) error {
	err := m.checkMailConfig()
	if err != nil {
		return err
	}

	return m.sendMail(subject, body, []string{}, to)
}

// SendAttach 发送带附件的普通文本邮件
func (m *Server) SendAttach(subject, body string, attach, to []string) error {
	err := m.checkMailConfig()
	if err != nil {
		return err
	}

	return m.sendMail(subject, body, attach, to)
}

// 检查mail配置信息
func (m *Server) checkMailConfig() error {
	if m.Addr == "" || m.Port == 0 || m.User == "" || m.Password == "" {
		err := errors.New("Check mail server [Addr|Port|User|Password] configration")
		return err
	}
	return nil
}

// 构造邮件内容并发送
func (m *Server) sendMail(subject, body string, attach, to []string) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", m.User)
	msg.SetHeader("To", to...)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/plain", body)

	for _, v := range attach {
		msg.Attach(v)
	}

	d := gomail.NewDialer(m.Addr, m.Port, m.User, m.Password)

	if err := d.DialAndSend(msg); err != nil {
		return err
	}
	return nil
}

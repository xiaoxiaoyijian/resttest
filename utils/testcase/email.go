package testcase

import (
	. "github.com/xiaoxiaoyijian/logger"
	"net/smtp"
	"strings"
)

type Email struct {
	User   string
	Passwd string
	Host   string
	To     string
}

func NewEmail(settings map[string]string) *Email {
	m := &Email{}

	for k, v := range settings {
		switch k {
		case "user":
			m.User = v
		case "passwd":
			m.Passwd = v
		case "host":
			m.Host = v
		case "to":
			m.To = v
		}
	}

	return m
}

func (m *Email) IsValid() bool {
	Logger.Infof("Email: user:%s, passwd:***, host:%s, to:%s", m.User, m.Host, m.To)

	return m.User != "" && m.Passwd != "" && m.To != "" && m.Host != ""
}

func (m *Email) Send(subject, content string) error {
	hp := strings.Split(m.Host, ":")
	auth := smtp.PlainAuth("", m.User, m.Passwd, hp[0])

	content_type := "Content-Type: text/plain" + "; charset=UTF-8"
	msg := []byte("To: " + m.To + "\nFrom: " + m.User + "<" + m.User + ">\nSubject: " + subject + "\n" + content_type + "\n\n" + content)
	Logger.Infof("Sended email to %s: \n%s", m.To, msg)
	return smtp.SendMail(m.Host, auth, m.User, strings.Split(m.To, ";"), msg)
}

package util

import (
	"github.com/KSkun/tqb-backend/config"
	"github.com/KSkun/tqb-backend/util/log"
	"gopkg.in/gomail.v2"
)

var (
	dialer *gomail.Dialer
	msgChan chan *gomail.Message
)

func initMail() {
	dialer = gomail.NewDialer(config.C.Mail.SMTPHost, config.C.Mail.SMTPPort,
		config.C.Mail.Addr, config.C.Mail.Password)
	msgChan = make(chan *gomail.Message, 100)

	go doSendMail()
}

func doSendMail() {
	for {
		msg := <- msgChan
		err := dialer.DialAndSend(msg)
		if err != nil {
			log.Logger.Error(err)
		}
	}
}

func SendMail(email string, title string, content string) {
	m := gomail.NewMessage()
	m.SetHeader("From", config.C.Mail.Addr)
	m.SetHeader("To", email)
	m.SetHeader("Subject", title)
	m.SetBody("text/html", content)

	msgChan <- m
}

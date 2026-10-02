package email

import (
	"fmt"
	"log"
	"net/smtp"
	"strings"
	"time"
)

type Sender struct {
	host     string
	port     string
	username string
	password string
	from     string
}

func New(host, port, username, password, from string) *Sender {
	return &Sender{host, port, username, password, from}
}

func (s *Sender) Send(to, subject, text, html string) error {
	startedAt := time.Now()
	defer func() {
		log.Printf("service=notification smtp_send duration_ms=%.3f", float64(time.Since(startedAt).Microseconds())/1000)
	}()
	if s.host == "" || s.from == "" || to == "" {
		return fmt.Errorf("email SMTP configuration and recipient are required")
	}
	content := text
	if html != "" {
		content = html
	}
	message := []byte("From: " + s.from + "\r\nTo: " + to + "\r\nSubject: " + subject + "\r\nMIME-Version: 1.0\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n" + content)
	var auth smtp.Auth
	if s.username != "" {
		auth = smtp.PlainAuth("", s.username, s.password, s.host)
	}
	return smtp.SendMail(s.host+":"+s.port, auth, s.from, strings.Fields(to), message)
}

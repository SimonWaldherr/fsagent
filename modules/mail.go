package modules

import (
	"crypto/tls"
	"gopkg.in/gomail.v2"
)

type mailConfig struct {
	Name    string   `json:"name"`
	Subject string   `json:"subject"`
	Body    string   `json:"body"`
	From    string   `json:"from"`
	To      []string `json:"to"`
	Cc      []string `json:"cc"`
	Bcc     []string `json:"bcc"`
	User    string   `json:"user"`
	Pass    string   `json:"pass"`
	Server  string   `json:"server"`
	Port    int      `json:"port"`
}

func addRecipients(rtype string, recipients []string, mail *gomail.Message) {
	addresses := make([]string, len(recipients))
	for i, recipient := range recipients {
		addresses[i] = mail.FormatAddress(recipient, "")
	}
	mail.SetHeader(rtype, addresses...)
}

type Mail struct {
}

func (Mail) Name() string {
	return "mail"
}

func (Mail) EmptyConfig() interface{} {
	return &mailConfig{}
}

func (Mail) Perform(config interface{}, fileName string) error {
	c := config.(mailConfig)

	m := gomail.NewMessage()
	m.SetHeader("From", c.From)

	addRecipients("To", c.To, m)
	addRecipients("Cc", c.Cc, m)
	addRecipients("Bcc", c.Bcc, m)

	m.SetHeader("Subject", c.Subject)
	m.SetBody("text/plain", c.Body)
	m.Attach(fileName)

	d := gomail.NewDialer(c.Server, c.Port, c.User, c.Pass)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	return d.DialAndSend(m)
}

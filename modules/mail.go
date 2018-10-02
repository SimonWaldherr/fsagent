package modules

import (
	"crypto/tls"
	"encoding/json"
	"gopkg.in/gomail.v2"
	"simonwaldherr.de/go/golibs/file"
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

func SendMail(configName, fileName string) error {
	var c mailConfig
	str, _ := file.Read(configName)
	err := json.Unmarshal([]byte(str), &c)
	if err == nil {

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

		if err := d.DialAndSend(m); err != nil {
			return err
		}
	}
	return err
}

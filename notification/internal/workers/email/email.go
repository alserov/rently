package email

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
)

type Emailer interface {
	Registration(to string) error
}

type Params struct {
	From     string
	Password string
	SmtpHost string
	SmtpPort int
}

func NewEmailer(p Params) Emailer {
	return &emailer{
		from:     p.From,
		password: p.Password,
		smtpHost: p.SmtpHost,
		smtpPort: p.SmtpPort,
	}
}

type emailer struct {
	from string

	password string

	smtpHost string
	smtpPort int
}

func (e *emailer) Registration(to string) error {
	auth := smtp.PlainAuth("", e.from, e.password, e.smtpHost)

	t, _ := template.ParseFiles("templates/auth.html")
	var body bytes.Buffer
	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: This is a test subject \n%s\n\n", mimeHeaders)))
	t.Execute(&body, struct {
		Name string
	}{Name: "useR"})

	err := smtp.SendMail(fmt.Sprintf("%s:%d", e.smtpHost, e.smtpPort), auth, e.from, []string{to}, body.Bytes())
	if err != nil {
		return err
	}
	return nil
}

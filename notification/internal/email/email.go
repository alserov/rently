package email

import (
	"fmt"
	"github.com/jordan-wright/email"
	"net/smtp"
	"os"
)

type Emailer interface {
	Registration(to string) error
}

type Params struct {
	From     string
	Password string
	Sender   string
	SmtpHost string
	SmtpPort int
}

func NewEmailer(p Params) Emailer {
	return &emailer{
		from:     p.From,
		sender:   p.Sender,
		password: p.Password,
		smtpHost: p.SmtpHost,
		smtpPort: p.SmtpPort,
	}
}

type emailer struct {
	from     string
	password string
	sender   string

	smtpHost string
	smtpPort int
}

func (e *emailer) Registration(to string) error {
	m := email.NewEmail()
	m.From = e.from
	m.Sender = e.sender
	m.To = []string{to}
	m.Subject = "Registration"
	b, err := os.ReadFile("internal/workers/email/templates/auth.html")
	if err != nil {
		return fmt.Errorf("failed to read html template file: %w", err)
	}
	m.HTML = b

	return m.Send(fmt.Sprintf("%s:%d", e.smtpHost, e.smtpPort), smtp.PlainAuth("", e.from, e.password, e.smtpHost))
}

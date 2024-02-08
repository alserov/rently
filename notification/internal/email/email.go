package email

import (
	"errors"
	"fmt"
	"github.com/jordan-wright/email"
	"net/smtp"
	"os"
	"strconv"
)

type Emailer interface {
	Send(t MessageType, to string) error
}

type MessageType int

const (
	REGISTER MessageType = iota
	LOGIN

	ERR_UNKNOWN_MESSAGE_TYPE = "unknown message type"
)

func MessageTypeStringToInt(val string) MessageType {
	t, err := strconv.Atoi(val)
	if err != nil {
		return -1
	}

	return MessageType(t)
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

func (e *emailer) Send(t MessageType, to string) error {
	var (
		m   = email.NewEmail()
		err error
	)
	switch t {
	case REGISTER:
		err = e.register(m, to)
	case LOGIN:
		err = e.login(m, to)
	default:
		err = errors.New(ERR_UNKNOWN_MESSAGE_TYPE)
	}

	fmt.Println(m)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return m.Send(fmt.Sprintf("%s:%d", e.smtpHost, e.smtpPort), smtp.PlainAuth("", e.from, e.password, e.smtpHost))
}

func (e *emailer) login(m *email.Email, to string) error {
	m.From = e.from
	m.Sender = e.sender
	m.To = []string{to}
	m.Subject = "Login"
	b, err := os.ReadFile("internal/email/templates/auth.html")
	if err != nil {
		return fmt.Errorf("failed to read html template file: %w", err)
	}
	m.HTML = b

	return nil
}

func (e *emailer) register(m *email.Email, to string) error {
	m.From = e.from
	m.Sender = e.sender
	m.To = []string{to}
	m.Subject = "Registration"
	b, err := os.ReadFile("internal/email/templates/auth.html")
	if err != nil {
		return fmt.Errorf("failed to read html template file: %w", err)
	}
	m.HTML = b

	return nil
}

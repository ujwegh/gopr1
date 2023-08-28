package models

import "github.com/go-mail/mail/v2"

const (
	// DefaultSender is the default email address to send emails from.
	DefaultSender = "support@lenslocked.com"
)

type EmailService struct {
	DefaultSender string
	dialer        *mail.Dialer
}
type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
}

func NewEmailService(config SMTPConfig) *EmailService {
	es := EmailService{
		dialer: mail.NewDialer(config.Host, config.Port, config.Username, config.Password),
	}
	return &es
}

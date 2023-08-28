package main

import (
	"errors"
	"fmt"
	"github.com/go-mail/mail/v2"
	_ "github.com/jackc/pgx/v4/stdlib"
	"os"
	"strings"
)

func Connect() error {
	// try to connect
	// pretend we got an error
	return errors.New("connection failed")
}
func CreateUser() error {
	err := Connect()
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}
	return nil
}
func CreateOrg() error {
	err := CreateUser()
	if err != nil {
		return fmt.Errorf("create org: %w", err)
	}
	return nil
}
func Demo(numbers ...int) {
	for _, number := range numbers {
		fmt.Print(number, " ")
	}
	fmt.Println()
}
func Sum(numbers ...int) int {
	sum := 0
	for i := 0; i < len(numbers); i++ {
		sum += numbers[i]
	}
	return sum
}

type Order struct {
	ID          int
	UserID      int
	Amount      int
	Description string
}

const (
	host     = "sandbox.smtp.mailtrap.io"
	port     = 587
	username = "fill this in"
	password = "fill this in"
)

func main() {
	from := "test@ujwegh.com"
	to := "nik29200018@gmail.com"
	subject := "This is a test email"
	plaintext := "This is the body of the email"
	html := `<h1>Hello there buddy!</h1><p>This is the email</p><p>Hope you enjoy it</p>`
	msg := mail.NewMessage()
	msg.SetHeader("To", to)
	msg.SetHeader("From", from)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/plain", plaintext)
	msg.AddAlternative("text/html", html)
	msg.WriteTo(os.Stdout)
	dialer := mail.NewDialer(host, port, username, password)
	err := dialer.DialAndSend(msg)
	if err != nil {
		// TODO: Handle the error correctly
		panic(err)
	}
}
func Join(vals ...string) string {
	var sb strings.Builder
	for i, s := range vals {
		sb.WriteString(s)
		if i < len(vals)-1 {
			sb.WriteString(", ")
		}
	}
	return sb.String()
}

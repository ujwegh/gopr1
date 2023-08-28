package main

import (
	stdctx "context"
	"errors"
	"fmt"
	_ "github.com/jackc/pgx/v4/stdlib"
	"gopr/context"
	"gopr/models"
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

func main() {
	ctx := stdctx.Background()
	user := models.User{
		Email: "jon@calhoun.io",
	}
	ctx = context.WithUser(ctx, &user)
	retrievedUser := context.User(ctx)
	fmt.Println(retrievedUser.Email)
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

package main

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/jackc/pgx/v4/stdlib"
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

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	SSLMode  string
}

func (cfg PostgresConfig) String() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database, cfg.SSLMode)
}

type Order struct {
	ID          int
	UserID      int
	Amount      int
	Description string
}

func main() {
	cfg := PostgresConfig{
		Host:     "localhost",
		Port:     "5432",
		User:     "postgres",
		Password: "mysecretpassword",
		Database: "postgres",
		SSLMode:  "disable",
	}
	db, err := sql.Open("pgx", cfg.String())
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name TEXT,
		email TEXT NOT NULL
		);
		CREATE TABLE IF NOT EXISTS orders (
		id SERIAL PRIMARY KEY,
		user_id INT NOT NULL,
		amount INT,
		description TEXT
		);`)
	if err != nil {
		panic(err)
	}
	fmt.Println("Tables created.")

	//name := "New User"
	//email := "new@calhoun.io"
	//row := db.QueryRow(`
	//	INSERT INTO users (name, email)
	//	VALUES ($1, $2) RETURNING id;`, name, email)
	//var id int
	//err = row.Scan(&id)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println("User created. id =", id)
	//userID := 3 // Pick an ID that exists in your DB
	//for i := 0; i < 5; i++ {
	//	amount := i * 100
	//	desc := fmt.Sprintf("Fake order #%d", i)
	//	_, err := db.Exec(`
	//		INSERT INTO orders(user_id, amount, description)
	//		VALUES($1, $2, $3)`, userID, amount, desc)
	//	if err != nil {
	//		panic(err)
	//	}
	//}

	//id := 5
	//row := db.QueryRow(`
	//	SELECT id, name, email
	//	FROM users
	//	WHERE id=$1;`, id)
	//var name, email string
	//err = row.Scan(&id, &name, &email)
	//if errors.Is(err, sql.ErrNoRows) {
	//	fmt.Println("Error, no rows!")
	//}
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Printf("User information: id:%d name=%s, email=%s\n", id, name, email)

	var orders []Order
	userID := 3 // Use the same ID you used in the previous lesson
	rows, err := db.Query(`
		SELECT id, amount, description
		FROM orders
		WHERE user_id=$1`, userID)
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var order Order
		order.UserID = userID
		err := rows.Scan(&order.ID, &order.Amount, &order.Description)
		if err != nil {
			panic(err)
		}
		orders = append(orders, order)
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	fmt.Println("Orders:", orders)

	defer db.Close()
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

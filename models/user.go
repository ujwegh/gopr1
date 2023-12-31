package models

import (
	"database/sql"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	appErrors "gopr/errors"
	"strings"
)

type User struct {
	ID           int
	Email        string
	PasswordHash string
}

type UserService struct {
	DB *sql.DB
}

func (us *UserService) Create(email, password string) (*User, error) {
	email = strings.ToLower(email)
	passwordHash := GeneratePasswordHash(password)
	user := User{
		Email:        email,
		PasswordHash: passwordHash,
	}
	row := us.DB.QueryRow(`
		INSERT INTO users (email, password_hash)
		VALUES ($1, $2) returning id`, user.Email, user.PasswordHash)
	err := row.Scan(&user.ID)
	if err != nil {
		// See if we can use this error as a PgError
		var pgError *pgconn.PgError
		if errors.As(err, &pgError) {
			// This is a PgError, so see if it matches a unique violation.
			if pgError.Code == pgerrcode.UniqueViolation {
				// If this is true, it has to be an email violation since this is the
				// only way to trigger this type of violation with our SQL.
				return nil, appErrors.ErrEmailTaken
			}
		}
		return nil, fmt.Errorf("create user: %w", err)
	}
	return &user, nil
}

func GeneratePasswordHash(password string) string {
	hashedBytes, err := bcrypt.GenerateFromPassword(
		[]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(fmt.Errorf("generate hash error: %w", err))
	}
	return string(hashedBytes)
}

type UserCreator interface {
	Create(email, password string) (*User, error)
}

func (us *UserService) Authenticate(email, password string) (*User, error) {
	email = strings.ToLower(email)
	user := User{
		Email: email,
	}
	row := us.DB.QueryRow(`SELECT id, password_hash FROM users WHERE email=$1`, email)
	err := row.Scan(&user.ID, &user.PasswordHash)
	if err != nil {
		fmt.Printf("authenticate: %v ", err)
		return nil, appErrors.ErrNotFound
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		fmt.Printf("authenticate: %v ", err)
		return nil, appErrors.ErrPasswordCheck
	}
	return &user, nil
}

func (us *UserService) UpdatePassword(userID int, password string) error {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("update password: %w", err)
	}
	passwordHash := string(hashedBytes)
	_, err = us.DB.Exec(`UPDATE users SET password_hash = $2 WHERE id = $1;`, userID, passwordHash)
	if err != nil {
		return fmt.Errorf("update password: %w", err)
	}
	return nil
}

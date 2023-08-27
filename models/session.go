package models

import (
	"database/sql"
	"fmt"
)

type Session struct {
	ID     int
	UserID int
	// Token is only set when creating a new session. When looking up a session
	// this will be left empty, as we only store the hash of a session token
	// in our database and we cannot reverse it into a raw token.
	Token     string
	TokenHash string
}

type SessionService struct {
	DB *sql.DB
	tm TokenManager
}

// Create will create a new session for the user provided. The session token
// will be returned as the Token field on the Session type, but only the hashed
// session token is stored in the database.
func (ss *SessionService) Create(userID int) (*Session, error) {
	token, tokenHash, err := ss.tm.New()
	if err != nil {
		return nil, fmt.Errorf("token create: %w", err)
	}
	session := Session{
		UserID:    userID,
		Token:     token,
		TokenHash: tokenHash,
	}
	row := ss.DB.QueryRow(`UPDATE sessions SET token_hash = $2 WHERE user_id = $1 RETURNING id;`,
		session.UserID, session.TokenHash)
	err = row.Scan(&session.ID)
	if err == sql.ErrNoRows {
		row = ss.DB.QueryRow(`INSERT INTO sessions (user_id, token_hash) VALUES ($1, $2) RETURNING id;`,
			session.UserID, session.TokenHash)
		err = row.Scan(&session.ID)
	}
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}
	return &session, nil
}

func (ss *SessionService) Delete(token string) error {
	tokenHash := ss.tm.hash(token)
	_, err := ss.DB.Exec(`DELETE FROM sessions WHERE token_hash = $1;`, tokenHash)
	if err != nil {
		return fmt.Errorf("delete: %w", err)
	}
	return nil
}

func (ss *SessionService) User(token string) (*User, error) {
	tokenHash := ss.tm.hash(token)
	var user User
	row := ss.DB.QueryRow(`SELECT u.id, email, password_hash FROM users u
    join sessions s on s.token_hash = $1 and u.id = s.user_id`, tokenHash)
	err := row.Scan(&user.ID, &user.Email, &user.PasswordHash)
	if err != nil {
		return nil, fmt.Errorf("user: %w", err)
	}
	return &user, nil
}

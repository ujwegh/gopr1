package models

import (
	"database/sql"
	"fmt"
	"gopr/rand"
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
	DB            *sql.DB
	BytesPerToken int
}

const (
	// The minimum number of bytes to be used for each session token.
	MinBytesPerToken = 32
)

// Create will create a new session for the user provided. The session token
// will be returned as the Token field on the Session type, but only the hashed
// session token is stored in the database.
func (ss *SessionService) Create(userID int) (*Session, error) {
	bytesPerToken := ss.BytesPerToken
	if bytesPerToken < MinBytesPerToken {
		bytesPerToken = MinBytesPerToken
	}
	token, err := rand.String(bytesPerToken)
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}
	// TODO: Hash the session token
	session := Session{
		UserID: userID,
		Token:  token,
		// TODO: Set the TokenHash
	}
	// TODO: Store the session in our DB
	return &session, nil
}

func (ss *SessionService) User(token string) (*User, error) {
	// TODO: Implement SessionService.User
	return nil, nil
}

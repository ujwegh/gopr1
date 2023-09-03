package models

import (
	"crypto/sha256"
	"encoding/base64"
	"gopr/rand"
)

type TokenManager struct {
	BytesPerToken int
}

const (
	// The minimum number of bytes to be used for each session token.
	MinBytesPerToken = 32
)

func (tm *TokenManager) New() (token, tokenHash string, err error) {
	bytesPerToken := tm.BytesPerToken
	if bytesPerToken < MinBytesPerToken {
		bytesPerToken = MinBytesPerToken
	}
	token, err = rand.String(bytesPerToken)
	tokenHash = tm.hash(token)
	return token, tokenHash, err
}

func (tm *TokenManager) hash(token string) string {
	tokenHash := sha256.Sum256([]byte(token))
	// base64 encode the data into a string
	return base64.URLEncoding.EncodeToString(tokenHash[:])
}

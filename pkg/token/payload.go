package token

import (
	"errors"
	"time"

	db "github.com/Luckny/space-it/db/sqlc"
	"github.com/google/uuid"
)

// Different types of error returned by the VerifyToken function
var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token has expired")
)

// Payload contains the payload data of the token
type Payload struct {
	ID         uuid.UUID         `json:"id"`
	User       db.User           `json:"user"`
	Attributes map[string]string `json:"attributes"`
	IssuedAt   time.Time         `json:"issued_at"`
	ExpiresAt  time.Time         `json:"expired_at"`
}

// NewPayload creates a new token payload with a specific user id and duration
func NewPayload(user db.User, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	payload := &Payload{
		ID:         tokenID,
		User:       user,
		Attributes: make(map[string]string),
		IssuedAt:   time.Now(),
		ExpiresAt:  time.Now().Add(duration),
	}
	return payload, nil
}

// Valid checks if the token payload is valid or not
func (payload *Payload) Valid() bool {
	if time.Now().After(payload.ExpiresAt) {
		return false
	}
	return true
}

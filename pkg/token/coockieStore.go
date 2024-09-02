package token

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/gob"
	"time"

	db "github.com/Luckny/space-it/db/sqlc"
	"github.com/Luckny/space-it/pkg/config"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
)

func init() {
	// Register the types that will be stored in the session
	gob.Register(map[string]string{})
	gob.Register(time.Time{})
	gob.Register(uuid.UUID{})
	gob.Register(db.User{})
}

type CookieStore struct {
	Name string
	Maker
	store    *sessions.CookieStore
	secure   bool
	httpOnly bool
}

// var Store = newCookieStore()
var SessionName = "_HOST-session"

func NewCookieStore() *CookieStore {
	return &CookieStore{
		Name:     "_HOST-session",
		store:    sessions.NewCookieStore([]byte(config.Envs.CookieSecret)),
		secure:   config.Envs.CookieIsSecure,
		httpOnly: config.Envs.CookieIsHttpOnly,
	}
}

// CreateToken create a new token for a specific user and duration
func (c *CookieStore) CreateToken(ctx *gin.Context, payload *Payload) (string, error) {
	session, err := c.store.Get(ctx.Request, c.Name)
	if err != nil {
		return "", err
	}

	if !session.IsNew {
		// invalidate old session
		session.Options.MaxAge = -1
		err = session.Save(ctx.Request, ctx.Writer)
		if err != nil {
			return "", err
		}

		session = sessions.NewSession(c.store, c.Name)
	}

	maxAge := int(payload.ExpiresAt.Sub(time.Now()).Seconds())
	session.Options.MaxAge = maxAge
	session.Options.Secure = c.secure
	session.Options.HttpOnly = c.httpOnly
	session.Options.Path = "/"

	// set session values
	session.Values["user"] = payload.User // TODO: dto
	session.Values["attributes"] = payload.Attributes
	session.Values["expiresAt"] = payload.ExpiresAt
	session.Values["sessionId"] = payload.ID
	session.Values["issuedAt"] = payload.IssuedAt

	// save session
	err = session.Save(ctx.Request, ctx.Writer)
	if err != nil {
		return "", err
	}

	b := sha256.Sum256([]byte(payload.ID.String()))
	return encodeToBase64(b[:]), nil
}

// VerifyToken checks if the token is valid
func (c *CookieStore) VerifyToken(ctx *gin.Context, tokenID string) (*Payload, error) {
	// get the session
	session, err := c.store.Get(ctx.Request, c.Name)
	if err != nil {
		return nil, err
	}

	if session.IsNew {
		return nil, ErrInvalidToken
	}

	// get the session id
	sessionId := session.Values["sessionId"].(uuid.UUID)

	// compare the provided token with the computed token
	provided, err := decodeBase64String(tokenID)
	if err != nil {
		return nil, err
	}

	computed := sha256.Sum256([]byte(sessionId.String()))
	// constant time comparison to help prevent timing attacks
	if subtle.ConstantTimeCompare(provided, computed[:]) != 1 {
		return nil, ErrInvalidToken
	}

	// TODO: user dto
	user := session.Values["user"].(db.User)

	// return the token
	token := &Payload{
		ID:         uuid.MustParse(sessionId.String()),
		User:       user,
		Attributes: session.Values["attributes"].(map[string]string),
		ExpiresAt:  session.Values["expiresAt"].(time.Time),
		IssuedAt:   session.Values["issuedAt"].(time.Time),
	}

	if ok := token.Valid(); !ok {
		return nil, ErrInvalidToken
	}

	return token, nil
}

package auth

import (
	"time"

	"efaturas-xtreme/pkg/hash"

	"github.com/google/uuid"
)

type Session struct {
	UserID    string
	Value     string
	CreatedAt int64
	ExpiresAt int64
	Username  string
	Password  string
}

func newSession(duration time.Duration, uname string, pword string) *Session {
	return &Session{
		UserID:    hash.NewUserID(uname, pword),
		Value:     hash.New(uuid.NewString()),
		CreatedAt: time.Now().Unix(),
		ExpiresAt: time.Now().Add(duration).Unix(),
		Username:  uname,
		Password:  pword,
	}
}

package auth

import (
	"context"
	"sync"
	"time"

	"efaturas-xtreme/pkg/efaturas"
	"efaturas-xtreme/pkg/errors"

	cache "github.com/jellydator/ttlcache/v3"
)

type Service struct {
	l           sync.RWMutex
	duration    time.Duration
	credentials *cache.Cache[string, *Session]

	efaturas efaturas.Service
}

func (s *Service) Login(ctx context.Context, uname string, pword string) (*Session, error) {
	s.l.Lock()
	defer s.l.Unlock()

	if _, err := s.efaturas.Login(ctx, uname, pword); err != nil {
		return nil, err
	}

	session := newSession(s.duration, uname, pword)
	s.credentials.Set(session.Value, session, cache.DefaultTTL)

	return session, nil
}

func (s *Service) GetByToken(ctx context.Context, token string) (*Session, error) {
	s.l.RLock()
	defer s.l.RUnlock()

	existing := s.credentials.Get(token)
	if existing != nil && existing.IsExpired() == false {
		return existing.Value(), nil
	}

	return nil, errors.NewNotFound("failed to get by token:", token, "- expired:", existing != nil)
}

func (s *Service) GetAndExtend(ctx context.Context, token string) (*Session, error) {
	s.l.Lock()
	defer s.l.Unlock()

	existing := s.credentials.Get(token)
	if existing == nil || existing.IsExpired() {
		return nil, errors.NewNotFound("failed to get by token:", token, "- expired:", existing != nil)
	}

	creds := existing.Value()
	session := newSession(s.duration, creds.Username, creds.Password)
	session.Value = creds.Value

	s.credentials.Set(session.Value, session, cache.DefaultTTL)

	return session, nil
}

func (s *Service) Logout(ctx context.Context, token string) error {
	s.l.Lock()
	defer s.l.Unlock()

	s.credentials.Delete(token)
	return nil
}

func New(efaturas efaturas.Service) *Service {
	duration := time.Hour * 32

	credentials := cache.New[string, *Session](cache.WithTTL[string, *Session](duration))
	go credentials.Start()

	return &Service{
		l:           sync.RWMutex{},
		duration:    duration,
		credentials: credentials,
		efaturas:    efaturas,
	}
}

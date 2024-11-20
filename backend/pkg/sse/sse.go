package sse

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"efaturas-xtreme/pkg/errors"

	"github.com/r3labs/sse/v2"
)

type Publisher interface {
	Publish(event string, obj interface{}) error
}

type Subscriber interface {
	Subscribe(w http.ResponseWriter, r *http.Request, event string)
}

type Server interface {
	Publisher
	Subscriber
}

type service struct {
	sse *sse.Server
}

func (s *service) Publish(event string, obj interface{}) error {
	data, err := json.Marshal(obj)
	if err != nil {
		return errors.New("failed to marshal json:", err)
	}

	s.sse.Publish(event, &sse.Event{Data: data})
	return nil
}

func (s *service) Subscribe(w http.ResponseWriter, r *http.Request, event string) {
	r.URL.RawQuery = fmt.Sprintf("stream=%s&%s", event, r.URL.RawQuery)
	s.sse.ServeHTTP(w, r)
}

func New() Server {
	s := sse.New()
	s.AutoStream = true
	s.EventTTL = time.Second

	return &service{sse: s}
}

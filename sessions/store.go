package sessions

import (
	"github.com/nu7hatch/gouuid"
)

type Session struct {
	Id     string            `json:"id"`
	Values map[string]string `json:"values"`
}

type Store interface {
	Create(values map[string]string) (*Session, error)

	Set(string, *Session) error

	Get(string) (*Session, error)

	Delete(string) error
}

func create(values map[string]string) *Session {
	id, _ := uuid.NewV4()
	sessionId := id.String()
	return &Session{
		Id:     sessionId,
		Values: values,
	}
}

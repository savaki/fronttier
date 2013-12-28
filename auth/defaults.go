package auth

import (
	"github.com/nu7hatch/gouuid"
	"github.com/savaki/fronttier/sessions"
)

var defaultIdFactory = func() string {
	id, _ := uuid.NewV4()
	return id.String()
}

var defaultSessionStore = sessions.Memory()

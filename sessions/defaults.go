package sessions

import (
	"github.com/nu7hatch/gouuid"
)

var defaultIdFactory = func() string {
	id, _ := uuid.NewV4()
	return id.String()
}

var defaultSessionStore = Memory()

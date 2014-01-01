package config

import (
	"strings"
)

var (
	NoProxyErr = "invalid route - no proxy specified"
	NoPathsErr = "invalid route - at least one path must be specified"
	NoRouteErr = "invalid config - at least one route must be specified"
)

type Errors struct {
	Messages map[string][]string
}

func (self *Errors) Add(key, value string) {
	if self.Messages == nil {
		self.Messages = make(map[string][]string)
	}
	self.Messages[key] = append(self.Messages[key], value)
}

func (self *Errors) Error() string {
	var gripes []string
	for _, values := range self.Messages {
		gripes = append(gripes, values...)
	}

	return strings.Join(gripes, ", ")
}

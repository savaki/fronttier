package config

import (
	"strings"
)

type Config struct {
	Routes []RouteConfig `json:"routes"`
}

type RouteConfig struct {
	Paths []string `json:"paths"`
	Proxy string   `json:"proxy"`
}

func (self Config) Validate() error {
	err := &Errors{}
	if self.Routes == nil || len(self.Routes) == 0 {
		gripe := "at least one route must be specified"
		err.Add("routes", gripe)
	}

	if len(err.Messages) > 0 {
		return err
	}
	return nil
}

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

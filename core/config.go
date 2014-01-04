package core

import (
	"bytes"
	"encoding/json"
	"github.com/savaki/fronttier/config"
	"github.com/savaki/fronttier/proxy"
	"io"
	"io/ioutil"
)

// Load configuration from the provided json file.  Any errors
// that occur during the load can be fetched when you invoke
// #Build on the returned *BuilderConfig
func (self *BuilderConfig) LoadFile(filename string) *BuilderConfig {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		self.err = err
		return self
	}

	return self.Load(bytes.NewReader(data))
}

// Load a json configuration from the provided reader.  As with
// #LoadFile, Load will return any errors when you invoke #Build
// on the returned *BuilderConfig
func (self *BuilderConfig) Load(reader io.Reader) *BuilderConfig {
	cfg := &config.Config{}
	err := json.NewDecoder(reader).Decode(cfg)
	if err != nil {
		self.err = err
		return self
	}

	err = cfg.Validate()
	if err != nil {
		self.err = err
		return self
	}

	for _, route := range cfg.Routes {
		handler, _ := proxy.Builder().Url(route.Proxy).Build()
		self.Paths(route.Paths...).Handler(handler)
	}

	return self
}

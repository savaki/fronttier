package fronttier

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
func LoadFile(filename string) *BuilderConfig {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		builder := Builder()
		builder.err = err
		return builder
	}

	return Load(bytes.NewReader(data))
}

// Load a json configuration from the provided reader.  As with
// #LoadFile, Load will return any errors when you invoke #Build
// on the returned *BuilderConfig
func Load(reader io.Reader) *BuilderConfig {
	builder := Builder()

	cfg := &config.Config{}
	err := json.NewDecoder(reader).Decode(cfg)
	if err != nil {
		builder.err = err
		return builder
	}

	err = cfg.Validate()
	if err != nil {
		builder.err = err
		return builder
	}

	for _, route := range cfg.Routes {
		handler, _ := proxy.Builder().Url(route.Proxy).Build()
		builder.Paths(route.Paths...).Handler(handler)
	}

	return builder
}

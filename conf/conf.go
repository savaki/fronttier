package conf

import (
	"encoding/json"
	"io/ioutil"
)

type Route struct {
	Paths          string
	Proxy          string
	SessionFactory bool `json:"session-factory"`
}

type Sessions struct {
	CookieName string `json:"cookie-name"`
}

type Config struct {
	Routes   []Route
	Sessions Sessions
}

func LoadFile(filename string) (*Config, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return Load(data)
}

func Load(data []byte) (*Config, error) {
	config := &Config{}
	err := json.Unmarshal(data, config)
	return config, err
}

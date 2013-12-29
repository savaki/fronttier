package config

type Config struct {
	Routes []RouteConfig `json:"routes"`
}

type RouteConfig struct {
	Paths []string `json:"paths"`
	Proxy string   `json:"proxy"`
}

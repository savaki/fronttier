package config

type RouteConfig struct {
	Paths []string `json:"paths"`
	Proxy string   `json:"proxy"`
}

func (self RouteConfig) validate(err *Errors) {
	if self.Proxy == "" {
		err.Add("routes", NoProxyErr)
	}

	if self.Paths == nil || len(self.Paths) == 0 {
		err.Add("routes", NoPathsErr)
	}
}

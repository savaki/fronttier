package config

type Config struct {
	Routes []RouteConfig `json:"routes"`
	Auth   *Auth         `json:"auth"`
}

func (self Config) Validate() error {
	err := &Errors{}

	if self.Routes == nil || len(self.Routes) == 0 {
		err.Add("routes", NoRouteErr)

	} else {
		for _, route := range self.Routes {
			route.validate(err)
		}
	}

	if self.Auth != nil {
		self.Auth.validate(err)
	}

	if len(err.Messages) > 0 {
		return err
	}
	return nil
}

package config

type Auth struct {
	Reserved     []string          `json:"reserved"`
	Logout       string            `json:"logout"`
	SessionStore map[string]string `json:"sessions"`
	Cookie       map[string]string `json:"cookie"`
}

func (self Auth) validate(err *Errors) {
	if len(self.Reserved) == 0 {
		err.Add("auth", "invalid auth - at least 1 reserved header must be specified")
	}
}

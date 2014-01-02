package fronttier

import (
	"github.com/savaki/fronttier/proxy"
)

func Proxy() *proxy.BuilderConfig {
	return proxy.Builder()
}

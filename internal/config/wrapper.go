package config

import (
	"github.com/harluo/config"
)

// Wrapper 汇率配置包装器。
type Wrapper struct {
	Rate *Rate `json:"rate,omitempty" validate:"required"`
}

func newWrapper(getter config.Getter) (wrapper *Wrapper, err error) {
	wrapper = &Wrapper{
		Rate: &Rate{Endpoint: "http://rate-server.service/api/private/parities/singles"},
	}
	err = getter.Get(wrapper)

	return
}

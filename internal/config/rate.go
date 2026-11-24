package config

import (
	"time"
)

// Rate 汇率客户端配置。
type Rate struct {
	Endpoint string        `default:"http://rate-server.service/api/private/parities/singles" json:"endpoint,omitempty" validate:"required,http_url"`
	TTL      time.Duration `default:"24h" json:"ttl,omitempty" validate:"gt=0"`
}

func newRate(yaothink *Yaothink) *Rate {
	return yaothink.Rate
}

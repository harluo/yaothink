package config

import (
	"time"
)

// Rate 汇率客户端配置。
type Rate struct {
	Hostname string        `default:"rate-server.service" json:"hostname,omitempty" validate:"hostname_rfc1123|ip"`
	Port     int           `default:"80" json:"port,omitempty" validate:"gte=1,lte=65535"`
	Prefix   string        `default:"/api/private" json:"prefix,omitempty" validate:"required,startswith=/"`
	TTL      time.Duration `default:"24h" json:"ttl,omitempty" validate:"gt=0"`
}

func newRate(yaothink *Yaothink) *Rate {
	return yaothink.Rate
}

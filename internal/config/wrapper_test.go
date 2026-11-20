package config

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/harluo/config"
	"github.com/harluo/di"
)

func TestRateEndpoint(t *testing.T) {
	t.Chdir(t.TempDir())
	var getter config.Getter
	if err := di.New().Instance().Get(func(injected config.Getter) {
		getter = injected
	}).Build().Inject(); err != nil {
		t.Fatal(err)
	}

	wrapper, err := newWrapper(getter)
	const endpoint = "http://rate-server.service/api/private/parities/singles"
	if err != nil || wrapper.Rate == nil || wrapper.Rate.Endpoint != endpoint || wrapper.Rate.TTL != 24*time.Hour {
		t.Fatalf("默认汇率配置错误：%+v，%v", wrapper.Rate, err)
	}

	if mkdirErr := os.Mkdir("config", 0o700); mkdirErr != nil {
		t.Fatal(mkdirErr)
	}
	if writeErr := os.WriteFile("config/rate.jsonc", []byte(`{"rate":{"endpoint":"ftp://rate-server.service/rates"}}`),
		0o600); writeErr != nil {
		t.Fatal(writeErr)
	}
	_, validationErr := newWrapper(getter)
	if validationErr == nil || !strings.Contains(validationErr.Error(), "http_url") {
		t.Fatalf("应拒绝非 HTTP 地址：%v", validationErr)
	}
}

package config

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/harluo/config"
	"github.com/harluo/di"
)

func TestRateConfig(t *testing.T) {
	t.Chdir(t.TempDir())
	var getter config.Getter
	if err := di.New().Instance().Get(func(injected config.Getter) {
		getter = injected
	}).Build().Inject(); err != nil {
		t.Fatal(err)
	}

	wrapper, err := newWrapper(getter)
	rate := wrapper.Yaothink.Rate
	if err != nil || rate == nil || rate.Hostname != "rate-server.service" || rate.Port != 80 ||
		rate.Prefix != "/api/private" || rate.TTL != 24*time.Hour {
		t.Fatalf("默认汇率配置错误：%+v，%v", rate, err)
	}

	if mkdirErr := os.Mkdir("config", 0o700); mkdirErr != nil {
		t.Fatal(mkdirErr)
	}
	for _, invalid := range []struct {
		name, settings, field string
	}{
		{name: "主机名无效", settings: `"hostname":"bad host"`, field: "Hostname"},
		{name: "端口低于下限", settings: `"port":-1`, field: "Port"},
		{name: "端口超过上限", settings: `"port":65536`, field: "Port"},
		{name: "前缀不以斜杠开头", settings: `"prefix":"api/private"`, field: "Prefix"},
	} {
		t.Run(invalid.name, func(t *testing.T) {
			contents := `{"yaothink":{"rate":{` + invalid.settings + `}}}`
			if writeErr := os.WriteFile("config/rate.jsonc", []byte(contents), 0o600); writeErr != nil {
				t.Fatal(writeErr)
			}
			_, validationErr := newWrapper(getter)
			if validationErr == nil || !strings.Contains(validationErr.Error(), invalid.field) {
				t.Fatalf("应拒绝无效的 %s：%v", invalid.field, validationErr)
			}
		})
	}
}

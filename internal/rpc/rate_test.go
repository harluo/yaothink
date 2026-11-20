package rpc

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/harluo/di"

	"github.com/harluo/yaothink/internal/core"
	"github.com/harluo/yaothink/internal/rpc/internal"
)

func TestRateValue(t *testing.T) {
	const path = "/custom/rates"
	var requests atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		count := requests.Add(1)
		if req.Method != http.MethodGet || req.URL.Path != path {
			t.Errorf("请求路由错误：%s %s", req.Method, req.URL.Path)
		}
		query := req.URL.Query()
		pair := query.Get("from") + "/" + query.Get("to")
		if count == 1 && pair == "USD/CNY" {
			w.WriteHeader(http.StatusServiceUnavailable)
		} else if pair == "USD/CNY" || pair == "CNY/USD" {
			_, _ = fmt.Fprintf(w, `{"value":%d}`, 7_000_000+count)
		} else {
			t.Errorf("币种对错误：%s", pair)
			w.WriteHeader(http.StatusBadRequest)
		}
	}))
	defer server.Close()

	t.Chdir(t.TempDir())
	if mkdirErr := os.Mkdir("config", 0o700); mkdirErr != nil {
		t.Fatal(mkdirErr)
	}
	contents := fmt.Sprintf(`{"yaothink":{"rate":{"endpoint":%q,"ttl":"1h"}}}`, server.URL+path)
	if writeErr := os.WriteFile("config/rate.jsonc", []byte(contents), 0o600); writeErr != nil {
		t.Fatal(writeErr)
	}
	var rate *Rate
	err := di.New().Instance().Get(func(injected *Rate) {
		rate = injected
	}).Build().Inject()
	if err != nil {
		t.Fatal(err)
	}
	if rate.config.TTL != time.Hour || rate.cache == nil {
		t.Fatalf("自定义缓存时长未生效：%s", rate.config.TTL)
	}

	ctx := t.Context()
	if _, err = rate.Value(ctx, 0, core.CNY); err == nil || !strings.Contains(err.Error(), "无效的币种") {
		t.Fatalf("无效币种：%v", err)
	}
	if _, err = rate.Value(ctx, core.USD, core.CNY); err == nil || !strings.Contains(err.Error(), "汇率服务返回非成功状态") {
		t.Fatalf("服务错误：%v", err)
	}
	for _, want := range []struct {
		from, to core.Currency
		value    uint64
	}{
		{from: core.USD, to: core.CNY, value: 7_000_002},
		{from: core.USD, to: core.CNY, value: 7_000_002},
		{from: core.CNY, to: core.USD, value: 7_000_003},
	} {
		if value, valueErr := rate.Value(ctx, want.from, want.to); valueErr != nil || value != want.value {
			t.Fatalf("%s/%s：得到 %d, %v；期望 %d", want.from, want.to, value, valueErr, want.value)
		}
	}

	key := [2]core.Currency{core.USD, core.CNY}
	if remaining := time.Until(rate.cache[key].Expires); remaining <= 0 || remaining > time.Hour {
		t.Fatalf("缓存时长错误：%s", remaining)
	}
	rate.cache[key] = internal.RateCacheEntry{Value: 7_000_002, Expires: time.Now().Add(-time.Second)}
	if value, valueErr := rate.Value(ctx, core.USD, core.CNY); valueErr != nil || value != 7_000_004 {
		t.Fatalf("过期后应重新查询：得到 %d, %v", value, valueErr)
	}
}

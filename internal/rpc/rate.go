package rpc

import (
	"context"
	"net"
	"net/url"
	"path"
	"strconv"
	"sync"
	"time"

	"github.com/goexl/exception"
	"github.com/goexl/gox/field"
	"github.com/goexl/http"
	"github.com/harluo/yaothink/internal/config"
	"github.com/harluo/yaothink/internal/core"
	"github.com/harluo/yaothink/internal/rpc/internal"
)

// 查询汇率服务并按币种对缓存成功结果，可并发使用。
type Rate struct {
	config *config.Rate
	client *http.Client
	mu     sync.RWMutex
	cache  map[[2]core.Currency]internal.RateCacheEntry
}

func newRate(settings *config.Rate, client *http.Client) *Rate {
	return &Rate{
		config: settings,
		client: client,
		cache:  make(map[[2]core.Currency]internal.RateCacheEntry),
	}
}

// 返回源币种到目标币种的汇率。
func (r *Rate) Value(ctx context.Context, from, to core.Currency) (value float64, err error) {
	key := [2]core.Currency{from, to}
	if from.String() == "" || to.String() == "" {
		err = exception.New().Message("无效的币种").Build()
	} else if cached, hit := r.cachedValue(key); hit {
		value = cached
	} else if value, err = r.fetch(ctx, from, to); err == nil {
		r.cacheValue(key, value)
	}

	return
}

func (r *Rate) cachedValue(key [2]core.Currency) (float64, bool) {
	r.mu.RLock()
	cached, exists := r.cache[key]
	r.mu.RUnlock()

	return cached.Value, exists && time.Now().Before(cached.Expires)
}

func (r *Rate) cacheValue(key [2]core.Currency, value float64) {
	r.mu.Lock()
	r.cache[key] = internal.RateCacheEntry{Value: value, Expires: time.Now().Add(r.config.TTL)}
	r.mu.Unlock()
}

func (r *Rate) fetch(ctx context.Context, from, to core.Currency) (value float64, err error) {
	requestCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	endpoint := url.URL{
		Scheme: "http",
		Host:   net.JoinHostPort(r.config.Hostname, strconv.Itoa(r.config.Port)),
		Path:   path.Join(r.config.Prefix, "parities/singles"),
	}
	result := new(internal.RateResponse)
	request := r.client.R().SetContext(requestCtx).ForceContentType("application/json").SetResult(result).
		SetQueryParam("from", from.String()).SetQueryParam("to", to.String())
	if rsp, requestErr := request.Get(endpoint.String()); requestErr != nil {
		err = requestErr
	} else if rsp.StatusCode() != 200 {
		err = exception.New().Message("汇率服务返回非成功状态").Field(field.New("status", rsp.Status())).Build()
	} else if result.Value == nil {
		err = exception.New().Message("汇率服务响应缺少汇率").Build()
	} else {
		value = *result.Value
	}

	return
}

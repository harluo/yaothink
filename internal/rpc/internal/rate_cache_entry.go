package internal

import (
	"time"
)

// 汇率缓存条目。
type RateCacheEntry struct {
	Value   uint64
	Expires time.Time
}

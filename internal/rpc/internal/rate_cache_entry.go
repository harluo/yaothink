package internal

import (
	"time"
)

// 汇率缓存条目。
type RateCacheEntry struct {
	Value   float64
	Expires time.Time
}

package cachex

import "time"

type CacheXOption func(*cachex)

func WithCacheBarrier(cacheBarrier bool) CacheXOption {
	return func(c *cachex) {
		c.cacheBarrier = cacheBarrier
	}
}

func WithUseLocal(useLocal bool) CacheXOption {
	return func(c *cachex) {
		c.useLocal = useLocal
	}
}

func WithFallbackMode(fallbackMode bool) CacheXOption {
	return func(c *cachex) {
		c.fallbackMode = fallbackMode
	}
}

func WithLocalTTL(localTTL time.Duration) CacheXOption {
	return func(c *cachex) {
		c.localTTL = localTTL
	}
}

func WithRedisTTL(redisTTL time.Duration) CacheXOption {
	return func(c *cachex) {
		c.redisTTL = redisTTL
	}
}

package cachex

import (
	"context"
	"errors"
	"github.com/TremblingV5/box/gofer"
	"github.com/TremblingV5/box/logx"
	"github.com/cloudwego/hertz/pkg/common/json"
	"github.com/go-redis/redis"
	"time"
)

const (
	NotFoundPlaceHolder = "{}"
)

var (
	ErrCachePlaceHolder = errors.New("use cache place holder")
)

type cachex struct {
	cacheBarrier bool // 是否启用防缓存击穿
	useLocal     bool // 是否启用本地缓存
	fallbackMode bool // 是否启用降级模式。降级模式：优先fetch数据，fetch失败后使用本地缓存->redis
	localTTL     time.Duration
	redisTTL     time.Duration
	client       *redis.Client
}

func New(options ...CacheXOption) *cachex {
	c := cachex{}

	for _, option := range options {
		option(&c)
	}

	return &c
}

func (c *cachex) fetchSet(ctx context.Context, key string, fetch func(ctx context.Context) (any, error)) (any, error) {
	value, err := fetch(ctx)
	if err != nil {
		if c.cacheBarrier {
			_ = c.client.Set(key, NotFoundPlaceHolder, c.redisTTL).Err()
			return nil, ErrCachePlaceHolder
		}
		return nil, err
	}

	val, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}

	ttl := c.redisTTL
	err = c.client.Set(key, val, ttl).Err()
	if err != nil {
		logx.WarnCtx(ctx, "fetchSet write redis failed: %v", err)
	}

	if c.useLocal {
		if c.localTTL > 0 {
			ttl = c.localTTL
		}

		// TODO: set local cache
	}

	return val, nil
}

func (c *cachex) getCache(ctx context.Context, key string) ([]byte, error) {
	if c.useLocal {
		// TODO: 读本地缓存并返回
	}

	result, err := c.client.Get(key).Result()
	if err != nil {
		return nil, err
	}

	if result == NotFoundPlaceHolder {
		return nil, ErrCachePlaceHolder
	}

	return []byte(result), nil
}

func (c *cachex) AutoFetch(
	ctx context.Context,
	key string,
	result any,
	fetch func(ctx context.Context) (any, error),
) error {
	value, err, _ := gofer.SingleFlightDo(key, func() (interface{}, error) {
		if c.fallbackMode {
			val, err := c.fetchSet(ctx, key, fetch)
			if err == nil {
				return val, nil
			}

			return c.getCache(ctx, key)
		}

		val, err := c.getCache(ctx, key)
		if err == nil {
			return val, nil
		}

		if errors.Is(err, ErrCachePlaceHolder) {
			return nil, err
		}

		return c.fetchSet(ctx, key, fetch)
	})

	gofer.SingleFlightForget(key)
	if err != nil {
		return err
	}

	return json.Unmarshal(value.([]byte), result)
}

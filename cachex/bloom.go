package cachex

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/spaolacci/murmur3"
)

const (
	maxOffset = 10000000 // 最大偏移量
	locations = 8        // 使用8个bit标识一个value是否存在

	setScript = `
	for _, offset in ipairs(ARGV) do
		redis.call("setbit", KEYS[1], offset, 1)
	end
	`

	getScript = `
	for _, offset in ipairs(ARGV) do
		if redis.call("getbit", KEYS[1], offset) == 0 then
			return 0
		end
	end
	return 1
	`
)

type Bloom struct {
	client *redis.Client
}

func NewBloom() *Bloom {
	return &Bloom{}
}

func (b *Bloom) Add(ctx context.Context, key string, value any) error {
	location, err := getLocation(value)
	if err != nil {
		return err
	}

	_, err = b.client.Eval(ctx, setScript, []string{key}, location).Result()
	if errors.Is(err, redis.Nil) {
		return nil
	}

	return err
}

func (b *Bloom) Existed(ctx context.Context, key string, value any) bool {
	location, err := getLocation(value)
	if err != nil {
		return false
	}

	result, err := b.client.Eval(ctx, getScript, []string{key}, location).Result()
	if err != nil {
		return false
	}

	if i, ok := result.(int64); ok && i == 1 {
		return true
	}

	return false
}

func getLocation(value any) ([]string, error) {
	byteData, err := getBytes(value)
	if err != nil {
		return nil, err
	}

	var data []string
	for i := 0; i < locations; i++ {
		s := fmt.Sprintf("%d", hashNum(byteData)%maxOffset)
		data = append(data, s)
	}

	return data, nil
}

func getBytes(value any) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(value)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func hashNum(values []byte) uint64 {
	return murmur3.Sum64(values)
}

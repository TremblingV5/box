package errorx

import (
	"sync"
)

type Serializer func(e *Error) string

var (
	isCustomSerializer   bool
	serializer           Serializer
	customSerializerOnce sync.Once
)

func CustomSerializer(s Serializer) {
	customSerializerOnce.Do(func() {
		isCustomSerializer = true
		serializer = s
	})
}

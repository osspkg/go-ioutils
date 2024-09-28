package cache

import "time"

type TCache[K comparable, T interface{}] interface {
	Has(key K) bool
	Get(key K) (T, bool)
	Set(key K, value T)
	SetWithTTL(key K, value T, ttl time.Time)
	Del(key K)
	Keys() []K
	Flush()
}

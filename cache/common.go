package cache

type TCache[K comparable, T interface{}] interface {
	Has(key K) bool
	Get(key K) (T, bool)
	Set(key K, value T)
	Del(key K)
	Keys() []K
	Flush()
}

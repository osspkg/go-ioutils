package pool

func NewSlicePool[T any](l, c int) *Pool[*SlicePool[T]] {
	return New(func() *SlicePool[T] {
		return &SlicePool[T]{B: make([]T, l, c)}
	})
}

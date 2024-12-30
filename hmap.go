package lockfreemap

import (
	"maps"
	"sync/atomic"
	"unsafe"
)

type Immutable[K comparable, V any] struct {
	hashmapAddr *unsafe.Pointer
}

func NewImmutable[K comparable, V any](size ...int) *Immutable[K, V] {
	return Create(
		make(map[K]V, append(size, 0)[0]),
	)
}

func Create[K comparable, V any](hmap map[K]V) *Immutable[K, V] {
	prt := unsafe.Pointer(&hmap)

	return &Immutable[K, V]{
		hashmapAddr: &prt,
	}
}

func (im *Immutable[K, V]) Copy() *Immutable[K, V] {
	return Create(im.GetValues())
}

func (im *Immutable[K, V]) GetValues() map[K]V {
	return maps.Clone(*(*map[K]V)(*im.hashmapAddr))
}

func (im *Immutable[K, V]) Get(k K) (V, bool) {
	v, ok := (*(*map[K]V)(*im.hashmapAddr))[k]

	return v, ok
}

func (im *Immutable[K, V]) Set(k K, v V) {
	im.Action(func(m map[K]V) {
		m[k] = v
	})
}

func (im *Immutable[K, V]) Del(k K) {
	im.Action(func(m map[K]V) {
		delete(m, k)
	})
}

func (im *Immutable[K, V]) Action(fn func(map[K]V)) {

	for {
		var (
			old    = im.hashmapAddr
			newest = im.GetValues()
		)
		fn(newest)

		if atomic.CompareAndSwapPointer(im.hashmapAddr, *old, unsafe.Pointer(&newest)) {
			break
		}
	}
}

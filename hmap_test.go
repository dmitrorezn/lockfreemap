package lockfreemap

import (
	"errors"
	"sync"
	"sync/atomic"
	"testing"
)

func TestSet(t *testing.T) {
	immap := NewImmutable[int, int]()

	immap.Set(1, 1)

	v, ok := immap.Get(1)
	if !ok {
		t.Fail()
	}
	if v != 1 {
		t.Fail()
	}
}

func TestDel(t *testing.T) {
	immap := NewImmutable[int, int]()

	immap.Set(1, 1)
	immap.Del(1)
	v, ok := immap.Get(1)
	if ok {
		t.Fail()
	}
	if v != 0 {
		t.Fail()
	}
}

func TestMultithreadGet(t *testing.T) {
	immap := NewImmutable[int, int]()

	wg := sync.WaitGroup{}
	count := new(atomic.Int32)

	errc := make(chan error)
	const N = 100
	wg.Add(N)
	for i := 0; i < N; i++ {
		go func() {
			defer wg.Done()
			count.Add(1)

			k := 1
			immap.Set(k, k)
			v, ok := immap.Get(k)
			if !ok {
				errc <- errors.New("fail !ok")
				return
			}
			if v != k {
				errc <- errors.New("fail v!=k")
			}
		}()
	}
	go func() {
		wg.Wait()
		close(errc)
	}()
	for err := range errc {
		if err != nil {
			t.Log(err.Error())
			t.Fail()
		}
	}

	//for i := 0; i < 100_000; i++ {
	// wg.Add(1)
	// go func() {
	//    defer wg.Done()
	//    count.Add(1)
	//    k := 1
	//    if v, ok := immap.Get(k); !ok || v != k {
	//       t.Fail()
	//    }
	// }()
	//}

	wg.Wait()
	if count.Load() != 1*N {
		t.Log(count)
		t.Fail()
	}
}

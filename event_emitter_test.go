package eventemitter

import (
	"math/rand"
	"sync"
	"testing"
	"time"
)

func Test_EventEmitter_On(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	ee := New()

	eventX := "x only"
	eventY := "y only"
	eventBoth := "both"

	var x, y int

	ee.On(eventX, func(args ...interface{}) {
		x++
	})

	ee.On(eventY, func(args ...interface{}) {
		y++
	})

	ee.On(eventBoth, func(args ...interface{}) {
		ee.Emit(eventX).Wait()
		ee.Emit(eventY).Wait()
	})

	// Emit eventX and wait for it to complete.
	ee.Emit(eventX).Wait()

	if x != 1 {
		t.Errorf("unexpected value: %v != %v", x, 1)
	}

	if y != 0 {
		t.Errorf("unexpected value: %v != %v", y, 0)
	}

	// Emit eventY twice and wait for them to complete.
	ee.Emit(eventY).Wait()
	ee.Emit(eventY).Wait()

	if x != 1 {
		t.Errorf("unexpected value: %v != %v", x, 1)
	}

	if y != 2 {
		t.Errorf("unexpected value: %v != %v", y, 2)
	}

	// Emit eventBoth and wait for it to complete.
	ee.Emit(eventBoth).Wait()

	if x != 2 {
		t.Errorf("unexpected value: %v != %v", x, 2)
	}

	if y != 3 {
		t.Errorf("unexpected value: %v != %v", y, 3)
	}
}

func Benchmark_EventEmitter(b *testing.B) {
	rand.Seed(time.Now().UnixNano())

	ee := New()

	e1 := "testing"

	wg := &sync.WaitGroup{}

	for idx := 0; idx < b.N; idx++ {
		fn := func(args ...interface{}) {
			args[0].(*sync.WaitGroup).Done()
		}

		ee.Once(e1, fn)
	}

	ee.Emit(e1, wg).Wait()

	if len(ee.listeners) > 0 {
		b.Error("listener count should be 0")
	}
}

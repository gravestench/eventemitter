package eventemitter

import (
	"math/rand"
	"testing"
	"time"
)

func Test_EventEmitter_On(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	ee := New()

	eventX := "x only"
	eventY := "y only"

	ee.On(eventX, func(args ...any) {})
	if len(ee.listeners) != 1 {
		t.Errorf("unexpected listener count")
	}

	ee.On(eventY, func(args ...any) {})

	if len(ee.listeners) != 2 {
		t.Errorf("unexpected listener count")
	}
}

func Test_EventEmitter_Off(t *testing.T) {
	ee := New()

	eventX := "x"
	eventY := "y"

	handler := func(...any) {}

	ee.On(eventX, handler)
	ee.On(eventY, handler)

	if len(ee.listeners) != 2 {
		t.Errorf("unexpected listener count")
	}

	ee.Off(eventX, handler)

	if len(ee.listeners[eventX]) != 0 {
		t.Errorf("unexpected listener count")
	}

	if len(ee.listeners) != 1 {
		t.Errorf("unexpected listener count")
	}

	ee.Off(eventY, handler)
	if len(ee.listeners[eventY]) != 0 {
		t.Errorf("unexpected listener count")
	}

	if len(ee.listeners) != 0 {
		t.Errorf("unexpected listener count")
	}
}

func Test_EventEmitter_Emit(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	ee := New()

	eventX := "x only"
	eventY := "y only"
	eventBoth := "both"

	var x, y int

	ee.On(eventX, func(args ...any) {
		x++
	})

	ee.On(eventY, func(args ...any) {
		y++
	})

	ee.On(eventBoth, func(args ...any) {
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

	for idx := 0; idx < b.N; idx++ {
		ee.Once(e1, func(...any) {})
	}

	ee.Emit(e1).Wait()

	if len(ee.listeners) > 0 {
		b.Error("listener count should be 0")
	}
}

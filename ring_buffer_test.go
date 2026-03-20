package talive

import (
	"testing"
)

func TestRingBufferCapacity(t *testing.T) {
	capacity := 2
	rb := newRingBuffer(capacity)
	for i := 0; i < 20; i++ {
		rb.Push(float64(i))
		if cap(rb.buffer) > capacity {
			t.Errorf("buffer capacity %d reached limit %d", cap(rb.buffer), capacity)
			break
		}
	}
}

func TestRingBuffer_Sum(t *testing.T) {
	rb := newRingBuffer(5)
	for i := 1; i <= 6; i++ {
		rb.Push(float64(i))
	}
	if rb.Sum != 20 {
		t.Errorf("wrong Sum value %f, expected 20.0", rb.Sum)
	}
}

func TestRingBuffer_SumExceptLast(t *testing.T) {
	rb := newRingBuffer(5)
	for i := 1; i <= 6; i++ {
		rb.Push(float64(i))
	}
	if rb.SumExceptLast() != 18 {
		t.Errorf("wrong SumExceptLast value %f, expected 18.0", rb.SumExceptLast())
	}
}

func TestRingBuffer_Min(t *testing.T) {
	rb := newRingBuffer(2)

	rb.Push(float64(3)) // removed
	rb.Push(float64(2)) // removed
	rb.Push(float64(1))
	rb.Push(float64(2))

	if rb.Min() != 1 {
		t.Errorf("wrong Min value %f, expected 1.0", rb.Min())
	}
}

func TestRingBuffer_MinExceptLast(t *testing.T) {
	rb := newRingBuffer(2)

	rb.Push(float64(3)) // removed
	rb.Push(float64(2)) // removed
	rb.Push(float64(1)) // next to rewrite - "last"
	rb.Push(float64(2))

	if rb.MinExceptLast() != 2 {
		t.Errorf("wrong Min value %f, expected 2.0", rb.MinExceptLast())
	}
}

func TestRingBuffer_Max(t *testing.T) {
	rb := newRingBuffer(2)
	rb.Push(float64(1)) // removed
	rb.Push(float64(2)) // removed
	rb.Push(float64(3))
	rb.Push(float64(2))

	if rb.Max() != 3 {
		t.Errorf("wrong Max value %f, expected 3.0", rb.Max())
	}
}

func TestRingBuffer_MaxExceptLast(t *testing.T) {
	rb := newRingBuffer(2)
	rb.Push(float64(1)) // removed
	rb.Push(float64(2)) // removed
	rb.Push(float64(3)) // next to rewrite - "last"
	rb.Push(float64(2))

	if rb.MaxExceptLast() != 2 {
		t.Errorf("wrong MaxExceptLast value %f, expected 2.0", rb.MaxExceptLast())
	}
}

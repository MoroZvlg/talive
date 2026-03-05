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

func TestRingBufferSum(t *testing.T) {
	rb := newRingBuffer(5)
	for i := 1; i <= 6; i++ {
		rb.Push(float64(i))
	}
	if rb.Sum != 20 {
		t.Errorf("wrong Sum value %f, expected 20.0", rb.Sum)
	}
}

func TestRingBufferSumExceptLast(t *testing.T) {
	rb := newRingBuffer(5)
	for i := 1; i <= 6; i++ {
		rb.Push(float64(i))
	}
	if rb.SumExceptLast() != 18 {
		t.Errorf("wrong SumExceptLast value %f, expected 18.0", rb.SumExceptLast())
	}
}

package talive

import "fmt"

type ringBuffer struct {
	buffer   []float64
	Sum      float64
	lastElem float64
	capacity int
	writeIdx int
}

func (buf *ringBuffer) String() string {
	return fmt.Sprintf(
		"[ringBuffer size:%d capacity:%d]",
		len(buf.buffer),
		buf.capacity,
	)
}

func newRingBuffer(capacity int) *ringBuffer {
	return &ringBuffer{
		buffer:   make([]float64, capacity),
		Sum:      0.0,
		lastElem: 0.0,
		capacity: capacity,
		writeIdx: 0,
	}
}

func (buf *ringBuffer) Push(el float64) {
	tailElement := buf.buffer[buf.writeIdx]
	buf.buffer[buf.writeIdx] = el
	buf.Sum = buf.Sum - tailElement + el
	buf.incrWriteIdx()
}

func (buf *ringBuffer) SumExceptLast() float64 {
	return buf.Sum - buf.buffer[buf.writeIdx]
}

func (buf *ringBuffer) incrWriteIdx() {
	if buf.writeIdx == (buf.capacity - 1) {
		buf.writeIdx = 0
	} else {
		buf.writeIdx++
	}
}

package talive

import (
	"fmt"
	"math"
)

type ringBuffer struct {
	buffer   []float64
	Sum      float64
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
		buffer:   make([]float64, 0, capacity),
		Sum:      0.0,
		capacity: capacity,
		writeIdx: 0,
	}
}

func (buf *ringBuffer) Push(el float64) {
	if len(buf.buffer) < cap(buf.buffer) {
		buf.buffer = append(buf.buffer, el)
		buf.Sum += el
		return
	}
	tailElement := buf.buffer[buf.writeIdx]
	buf.buffer[buf.writeIdx] = el
	buf.Sum = buf.Sum - tailElement + el
	buf.incrWriteIdx()
}

func (buf *ringBuffer) SumExceptLast() float64 {
	return buf.Sum - buf.Last()
}

func (buf *ringBuffer) Last() float64 {
	if len(buf.buffer) < cap(buf.buffer) {
		return 0.0
	}
	return buf.buffer[buf.writeIdx]
}

func (buf *ringBuffer) Min() float64 {
	minV := math.Inf(1)
	for _, el := range buf.buffer {
		if el < minV {
			minV = el
		}
	}
	return minV
}

func (buf *ringBuffer) MinExceptLast() float64 {
	minV := math.Inf(1)
	for i, el := range buf.buffer {
		if i == buf.writeIdx {
			continue
		}
		if el < minV {
			minV = el
		}
	}
	return minV
}

func (buf *ringBuffer) Max() float64 {
	maxV := math.Inf(-1)
	for _, el := range buf.buffer {
		if el > maxV {
			maxV = el
		}
	}
	return maxV
}

func (buf *ringBuffer) MaxExceptLast() float64 {
	maxV := math.Inf(-1)
	for i, el := range buf.buffer {
		if i == buf.writeIdx {
			continue
		}
		if el > maxV {
			maxV = el
		}
	}
	return maxV
}

func (buf *ringBuffer) MinMax() (float64, float64) {
	minV := math.Inf(1)
	maxV := math.Inf(-1)
	for _, el := range buf.buffer {
		if el < minV {
			minV = el
		}
		if el > maxV {
			maxV = el
		}
	}
	return minV, maxV
}

func (buf *ringBuffer) MinMaxExceptLast() (float64, float64) {
	minV := math.Inf(1)
	maxV := math.Inf(-1)
	for i, el := range buf.buffer {
		if i == buf.writeIdx {
			continue
		}
		if el < minV {
			minV = el
		}
		if el > maxV {
			maxV = el
		}
	}
	return minV, maxV
}

func (buf *ringBuffer) Len() int {
	return len(buf.buffer)
}

func (buf *ringBuffer) incrWriteIdx() {
	if buf.writeIdx == (buf.capacity - 1) {
		buf.writeIdx = 0
	} else {
		buf.writeIdx++
	}
}

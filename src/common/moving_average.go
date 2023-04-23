package common

import "fmt"

type MovingAverage[N Numeric] struct {
	buf     RingBuf[N]
	currAvg N
}

func NewMovingAverage[N Numeric](order uint) *MovingAverage[N] {
	return &MovingAverage[N]{
		buf:     NewRingBuf[N]((int)(order)),
		currAvg: 0,
	}
}

func (avg *MovingAverage[N]) Add(value N) N {
	avg.currAvg += value / N(avg.buf.Len())
	avg.currAvg -= avg.buf.Add(value) / N(avg.buf.Len())
	return avg.currAvg
}

func (avg *MovingAverage[N]) Order() int {
	return avg.buf.Len()
}

func (avg *MovingAverage[N]) Resize(order uint) N {
	if order > (uint)(avg.Order()) {
		avg.Grow(order - (uint)(avg.Order()))
	} else if order < (uint)(avg.Order()) {
		avg.Shrink((uint)(avg.Order()) - order)
	}

	return avg.currAvg
}

func (avg *MovingAverage[N]) Shrink(amount uint) N {
	fmt.Println("shrink")
	avg.currAvg *= N(avg.Order())
	for _, removed := range avg.buf.Shrink(amount) {
		avg.currAvg -= removed
	}
	avg.currAvg /= N(avg.Order())

	return avg.currAvg
}

func (avg *MovingAverage[N]) Grow(amount uint) N {
	fmt.Println("grow")
	avg.currAvg *= N(avg.Order())
	for _, added := range avg.buf.Grow(amount) {
		avg.currAvg += added
	}
	avg.currAvg /= N(avg.Order())

	return avg.currAvg
}

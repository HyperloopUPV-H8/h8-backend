package common

type MovingAverage[N Numeric] struct {
	buf     RingBuf[N]
	elems   int
	currAvg N
}

func NewMovingAverage[N Numeric](order uint) *MovingAverage[N] {
	return &MovingAverage[N]{
		buf:     NewRingBuf[N]((int)(order)),
		elems:   0,
		currAvg: 0,
	}
}

func (avg *MovingAverage[N]) Add(value N) N {
	if avg.elems > avg.Order() {
		avg.elems = avg.Order()
	}

	if avg.elems == avg.Order() {
		rem := avg.buf.Add(value)
		avg.currAvg -= rem / N(avg.Order())
		avg.currAvg += value / N(avg.Order())
	} else {
		avg.addElem(value)
	}
	return avg.currAvg
}

func (avg *MovingAverage[N]) addElem(value N) {
	prevElems := avg.elems
	avg.elems += 1
	avg.buf.Add(value)
	avg.currAvg = (N(prevElems)*avg.currAvg + value) / N(avg.elems)
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
	avg.currAvg *= N(avg.Order())
	for _, removed := range avg.buf.Shrink(amount) {
		avg.currAvg -= removed
	}
	avg.currAvg /= N(avg.Order())

	return avg.currAvg
}

func (avg *MovingAverage[N]) Grow(amount uint) N {
	avg.buf.Grow(amount)

	return avg.currAvg
}

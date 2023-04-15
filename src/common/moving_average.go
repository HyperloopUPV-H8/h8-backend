package common

type MovingAverage[N Numeric] struct {
	buf     RingBuf[N]
	currAvg N
}

func NewMovingAverage[N Numeric](order uint) MovingAverage[N] {
	return MovingAverage[N]{
		buf:     NewRingBuf[N]((int)(order)),
		currAvg: 0,
	}
}

func (avg *MovingAverage[N]) Add(value N) N {
	avg.currAvg += value / N(avg.buf.Len())
	avg.currAvg -= avg.buf.Add(value) / N(avg.buf.Len())
	return avg.currAvg
}

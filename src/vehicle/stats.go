package vehicle

type Stats struct {
	recv     int
	sent     int
	sentFail int
}

func newStats() *Stats {
	return new(Stats)
}

func (vehicle *Vehicle) Stats() Stats {
	return *vehicle.stats
}

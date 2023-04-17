package vehicle

type Stats struct {
	recv     int
	sent     int
	sentFail int
}

func (vehicle *Vehicle) Stats() Stats {
	vehicle.statsMx.Lock()
	defer vehicle.statsMx.Unlock()
	return vehicle.stats
}

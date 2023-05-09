package models

type TimeAccumulation struct {
	Value     float64 `json:"value"`
	Bound     float64 `json:"bound"`
	TimeLimit float64 `json:"timelimit"`
}

package domain

type PodData struct {
	Boards       []Board  `json:"boards"`
	LastBatchIDs []uint16 `json:"lastBatchIDs"`
}

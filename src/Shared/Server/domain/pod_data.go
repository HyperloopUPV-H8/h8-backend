package domain

type PodData struct {
	Boards map[string]Board `json:"boards"`
}

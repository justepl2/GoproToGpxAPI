package request

type CreateVideo struct {
	Name     string  `json:"name"`
	Duration float64 `json:"duration"`
}

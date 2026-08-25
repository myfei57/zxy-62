package cable

type Cable struct {
	ID     string  `json:"id"`
	Cabin  string  `json:"cabin"`
	Temp   float64 `json:"temp"`
	Status string  `json:"status"`
}

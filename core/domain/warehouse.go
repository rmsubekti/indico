package domain

type Warehouse struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Address  string `json:"address"`
	Capacity uint   `json:"capacity"`
}

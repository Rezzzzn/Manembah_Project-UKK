package models

type Accommodation struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Type        string  `json:"type"`
	Description string  `json:"description"`
	Location    string  `json:"location"`
	Price       float64 `json:"price"`
}

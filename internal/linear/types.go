package linear

// State represents the state of a Linear issue
type State struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// Issue represents a Linear issue
type Issue struct {
	ID         string `json:"id"`
	Identifier string `json:"identifier"`
	Title      string `json:"title"`
	State      State  `json:"state"`
}

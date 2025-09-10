package entities

type Comp struct {
	ID        string   `json:"id"`
	Game      string   `json:"game"`
	Region    string   `json:"region"`
	League    string   `json:"league,omitempty"`
	Patch     string   `json:"patch"`
	Map       string   `json:"map,omitempty"`
	Side      string   `json:"side,omitempty"`
	SlotsJSON string   `json:"slots"` // raw JSON (permitimos estructura libre)
	PickRate  *float64 `json:"pick_rate,omitempty"`
	WinRate   *float64 `json:"win_rate,omitempty"`
	DeltaWin  *float64 `json:"delta_win,omitempty"`
}

package entities

type Patch struct {
	ID         string `json:"id"`
	Game       string `json:"game"` // "lol" | "val"
	Version    string `json:"version"`
	ReleasedAt string `json:"released_at"`
}

type Change struct {
	PatchID    string  `json:"patch_id"`
	EntityType string  `json:"entity_type"` // champion|agent|item|weapon|map
	EntityID   string  `json:"entity_id"`
	Field      string  `json:"field"`
	Old        string  `json:"old"`
	New        string  `json:"new"`
	Impact     float64 `json:"impact_score"`
}

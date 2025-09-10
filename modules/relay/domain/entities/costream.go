package entities

type CoStream struct {
	ID       string `json:"id"`
	EventID  string `json:"event_id"`
	Platform string `json:"platform"` // twitch|youtube
	URL      string `json:"url"`
	Lang     string `json:"lang"`
	Country  string `json:"country"`
	Viewers  int    `json:"viewers"`
	Verified bool   `json:"verified"`
	Live     bool   `json:"is_live"`
}

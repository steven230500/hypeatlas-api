package entities

type League struct {
	ID     string `json:"id"`
	Game   string `json:"game"`
	Region string `json:"region"`
	Name   string `json:"name"`
	Slug   string `json:"slug"`
}

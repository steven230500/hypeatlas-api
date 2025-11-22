package entities

// Champion represents a League of Legends champion
type Champion struct {
	ID    string `json:"id"`
	Key   string `json:"key"`
	Name  string `json:"name"`
	Title string `json:"title"`
}

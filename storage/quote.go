package storage

type Quote struct {
	Quote  string `json:"quote"`
	Author string `json:"author"`
}

type QuoteStore struct {
	Quote  string `json:"quote"`
	Author string `json:"author"`
	ID     int    `json:"id"`
}

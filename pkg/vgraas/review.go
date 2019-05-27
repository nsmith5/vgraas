package vgraas

// Review is a video game review
type Review struct {
	Title    string    `json:"title"`
	Body     string    `json:"body"`
	Author   string    `json:"author"`
	Comments []Comment `json:"comments"`
}

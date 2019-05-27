package vgraas

// Comment is a comment on a video game review
type Comment struct {
	Body   string `json:"body"`
	Author string `json:"author"`
}

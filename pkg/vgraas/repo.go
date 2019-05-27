package vgraas

import "errors"

// These two domain specific errors should be used when
// implementing the Repo interface.
var (
	ReviewNotFound  = errors.New("Review not found")
	CommentNotFound = errors.New("Comment not found")
)

// Repo is an interface that an storage mechanism for reviews
// should obey.
type Repo interface {
	// All Reviews
	ReadReviews() ([]Review, error)

	// Review CRUD
	CreateReview(r Review) (id int, err error)
	ReadReview(id int) (Review, error)
	UpdateReview(id int, r Review) error
	DeleteReview(id int) error

	// All Comments
	ReadComments(postID int) ([]Comment, error)

	// Comment CRUD
	CreateComment(reviewID int, r Comment) (id int, err error)
	ReadComment(reviewID, id int) (Comment, error)
	UpdateComment(reviewID, id int, c Comment) error
	DeleteComment(reviewID, id int) error
}

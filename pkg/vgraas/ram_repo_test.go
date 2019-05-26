package vgraas

import (
	"testing"
)

func TestReviewCRUD(t *testing.T) {
	rr := NewRAMRepo()
	{
		/* CRUD that should work */
		reviews, err := rr.ReadReviews()
		if err != nil {
			t.Error("Failed to read reviews on empty repo")
		}
		if len(reviews) != 0 {
			t.Error("Empty repo should have no reviews")
		}

		review := Review{Author: "author", Body: "body", Comments: nil}
		id, err := rr.CreateReview(review)
		if err != nil {
			t.Error("Failed to created new review")
		}

		read, err := rr.ReadReview(id)
		if err != nil {
			t.Error("Failed to read created review")
		}
		if read.Author != review.Author || read.Body != review.Body || read.Comments != nil {
			t.Error("Submitted review not equal to retreived review")
		}

		err = rr.UpdateReview(id, Review{Author: "author2", Body: review.Body, Comments: review.Comments})
		if err != nil {
			t.Error("Failed to update review")
		}

		err = rr.DeleteReview(id)
		if err != nil {
			t.Error("Failed to delete review")
		}
	}
	{
		/* CRUD that shoudn't Work */
		_, err := rr.ReadReview(0)
		if err != ReviewNotFound {
			t.Error("Read review that doesn't exist")
		}

		err = rr.UpdateReview(0, Review{Author: "author2"})
		if err != ReviewNotFound {
			t.Error("Updated review that doesn't exist")
		}

		err = rr.DeleteReview(0)
		if err != ReviewNotFound {
			t.Error("Deleted review that doesn't exist")
		}
	}
}

func TestCommentCRUD(t *testing.T) {
	rr := NewRAMRepo()
	{
		/* CRUD that should Work */
		id, err := rr.CreateReview(Review{})
		if err != nil {
			t.Fatal("Failed to add review")
		}

		cid, err := rr.CreateComment(id, Comment{})
		if err != nil {
			t.Error("Failed to create comment on review")
		}

		err = rr.UpdateComment(id, cid, Comment{Author: "author"})
		if err != nil {
			t.Error("Failed to update comment")
		}

		read, err := rr.ReadComment(id, cid)
		if err != nil {
			t.Error("Failed to read comment")
		}
		if read.Author != "author" {
			t.Error("Failed to update comment in repository")
		}

		err = rr.DeleteComment(id, cid)
		if err != nil {
			t.Error("Failed to delete review")
		}

	}
	{
		/* CRUD that shouldn't Work */
		_, err := rr.CreateComment(10, Comment{})
		if err != ReviewNotFound {
			t.Error("Added comment to review that doesn't exist")
		}

		id, err := rr.CreateReview(Review{})
		if err != nil {
			t.Fatal("Failed to create review")
		}

		_, err = rr.ReadComment(id, 2)
		if err != CommentNotFound {
			t.Error("Read comment that doesn't exist")
		}

		err = rr.UpdateComment(id, 2, Comment{})
		if err != CommentNotFound {
			t.Error("Updated comment that doesn't exits")
		}

		err = rr.DeleteComment(id, 2)
		if err != CommentNotFound {
			t.Error("Delete comment that doesn't exist")
		}
	}
}

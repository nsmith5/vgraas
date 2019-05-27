package vgraas

import "sync"

type ramRepo struct {
	sync.RWMutex
	Reviews []Review
}

func NewRAMRepo() Repo {
	return &ramRepo{}
}

func (rr *ramRepo) ReadReviews() ([]Review, error) {
	rr.RLock()
	defer rr.RUnlock()
	return rr.Reviews, nil
}

/* Review CRUD */

func (rr *ramRepo) CreateReview(r Review) (id int, err error) {
	rr.Lock()
	defer rr.Unlock()

	rr.Reviews = append(rr.Reviews, r)
	return len(rr.Reviews) - 1, nil
}

func (rr *ramRepo) ReadReview(id int) (Review, error) {
	rr.RLock()
	defer rr.RUnlock()

	if id > len(rr.Reviews)-1 {
		return Review{}, ReviewNotFound
	}
	return rr.Reviews[id], nil
}

func (rr *ramRepo) UpdateReview(id int, r Review) error {
	rr.Lock()
	defer rr.Unlock()

	if id > len(rr.Reviews)-1 {
		return ReviewNotFound
	}
	rr.Reviews[id] = r
	return nil
}

func (rr *ramRepo) DeleteReview(id int) error {
	rr.Lock()
	defer rr.Unlock()

	if id > len(rr.Reviews)-1 {
		return ReviewNotFound
	}

	rr.Reviews = append(rr.Reviews[:id], rr.Reviews[id+1:]...)
	return nil
}

func (rr *ramRepo) ReadComments(reviewID int) ([]Comment, error) {
	rr.RLock()
	defer rr.RUnlock()

	if reviewID > len(rr.Reviews)-1 {
		return nil, ReviewNotFound
	}
	return rr.Reviews[reviewID].Comments, nil
}

/* Comment CRUD */
func (rr *ramRepo) CreateComment(reviewID int, c Comment) (id int, err error) {
	rr.Lock()
	defer rr.Unlock()

	if reviewID > len(rr.Reviews)-1 {
		return 0, ReviewNotFound
	}

	rr.Reviews[reviewID].Comments = append(rr.Reviews[reviewID].Comments, c)
	return len(rr.Reviews[reviewID].Comments) - 1, nil
}

func (rr *ramRepo) ReadComment(reviewID, id int) (Comment, error) {
	rr.RLock()
	defer rr.RUnlock()

	if reviewID > len(rr.Reviews)-1 {
		return Comment{}, ReviewNotFound
	}

	if id > len(rr.Reviews[reviewID].Comments)-1 {
		return Comment{}, CommentNotFound
	}

	return rr.Reviews[reviewID].Comments[id], nil
}

func (rr *ramRepo) UpdateComment(reviewID, id int, c Comment) error {
	rr.Lock()
	defer rr.Unlock()

	if reviewID > len(rr.Reviews)-1 {
		return ReviewNotFound
	}

	if id > len(rr.Reviews[reviewID].Comments)-1 {
		return CommentNotFound
	}

	rr.Reviews[reviewID].Comments[id] = c
	return nil
}

func (rr *ramRepo) DeleteComment(reviewID, id int) error {
	rr.Lock()
	defer rr.Unlock()

	if reviewID > len(rr.Reviews)-1 {
		return ReviewNotFound
	}

	if id > len(rr.Reviews[reviewID].Comments)-1 {
		return CommentNotFound
	}

	rr.Reviews[reviewID].Comments =
		append(rr.Reviews[reviewID].Comments[:id], rr.Reviews[reviewID].Comments[id+1:]...)
	return nil
}

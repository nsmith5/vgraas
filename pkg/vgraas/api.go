package vgraas

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type API struct {
	Repo
	*mux.Router
}

type Route struct {
	Name        string
	Methods     string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

func NewAPI(r Repo) http.Handler {
	var a API
	a.Repo = r
	a.Router = mux.NewRouter().StrictSlash(true)

	routes := []Route{
		/* All Reviews */
		Route{"ReadReviews", "GET", "/reviews/", a.ReadReviews},

		/* Review CRUD */
		Route{"CreateReview", "POST", "/reviews/", a.CreateReview},
		Route{"ReadReview", "GET", "/reviews/{id}", a.ReadReview},
		Route{"UpdateReview", "PUT", "/reviews/{id}", a.UpdateReview},
		Route{"DeleteReview", "DELETE", "/reviews/{id}", a.DeleteReview},

		/* All Comments */
		Route{"ReadComments", "GET", "/reviews/{rid}/comments", a.ReadComments},

		/* Comments CRUD */
		Route{"CreateComment", "POST", "/reviews/{rid}/comments", a.CreateComment},
		Route{"ReadComment", "GET", "/reviews/{rid}/comments/{id}", a.ReadComment},
		Route{"UpdateComment", "PUT", "/reviews/{rid}/comments/{id}", a.UpdateComment},
		Route{"DeleteComment", "DELETE", "/reviews/{rid}/comments/{id}", a.DeleteComment},
	}

	for _, route := range routes {
		a.Router.
			Methods(route.Methods).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}

	return a
}

func HandleError(w http.ResponseWriter, r *http.Request, status int, err string) {
	w.WriteHeader(status)
	fmt.Fprintf(w, `{"err": "%s"}`, err)
}

func (a API) ReadReviews(w http.ResponseWriter, r *http.Request) {
	reviews, err := a.Repo.ReadReviews()
	if err != nil {
		HandleError(w, r, http.StatusInternalServerError, "")
		return
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "\t")
	err = enc.Encode(reviews)
	if err != nil {
		HandleError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
}

func (a API) CreateReview(w http.ResponseWriter, r *http.Request) {
	var review Review
	{
		dec := json.NewDecoder(r.Body)
		defer r.Body.Close()
		err := dec.Decode(&review)
		if err != nil {
			HandleError(w, r, http.StatusBadRequest, err.Error())
			return
		}
	}

	id, err := a.Repo.CreateReview(review)
	if err != nil {
		HandleError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	_, err = fmt.Fprintf(w, `{"id": %d}`, id)
	if err != nil {
		HandleError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
}

func (a API) ReadReview(w http.ResponseWriter, r *http.Request) {
	var id int
	{
		vars := mux.Vars(r)
		_, err := fmt.Sscanf(vars["id"], "%d", &id)
		if err != nil {
			HandleError(w, r, http.StatusBadRequest, err.Error())
			return
		}
	}

	review, err := a.Repo.ReadReview(id)
	if err != nil {
		HandleError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "\t")
	err = enc.Encode(review)
	if err != nil {
		HandleError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	return
}

func (a API) UpdateReview(w http.ResponseWriter, r *http.Request) {
	var id int
	{
		vars := mux.Vars(r)
		_, err := fmt.Sscanf(vars["id"], "%d", &id)
		if err != nil {
			HandleError(w, r, http.StatusBadRequest, err.Error())
			return
		}
	}

	var review Review
	{
		dec := json.NewDecoder(r.Body)
		defer r.Body.Close()
		err := dec.Decode(&review)
		if err != nil {
			HandleError(w, r, http.StatusBadRequest, err.Error())
			return
		}
	}

	err := a.Repo.UpdateReview(id, review)
	switch {
	case err == ReviewNotFound:
		HandleError(w, r, http.StatusBadRequest, err.Error())
		return
	case err != nil:
		HandleError(w, r, http.StatusBadRequest, err.Error())
		return
	}
}

func (a API) DeleteReview(w http.ResponseWriter, r *http.Request) {
	var id int
	{
		vars := mux.Vars(r)
		_, err := fmt.Sscanf(vars["id"], "%d", &id)
		if err != nil {
			HandleError(w, r, http.StatusBadRequest, err.Error())
			return
		}
	}

	err := a.Repo.DeleteReview(id)
	switch {
	case err == ReviewNotFound:
		HandleError(w, r, http.StatusBadRequest, err.Error())
		return
	case err != nil:
		HandleError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
}

func (a API) ReadComments(w http.ResponseWriter, r *http.Request) {
	var id int
	{
		vars := mux.Vars(r)
		_, err := fmt.Sscanf(vars["rid"], "%d", &id)
		if err != nil {
			HandleError(w, r, http.StatusBadRequest, err.Error())
			return
		}
	}

	comments, err := a.Repo.ReadComments(id)
	switch {
	case err == ReviewNotFound:
		HandleError(w, r, http.StatusBadRequest, err.Error())
		return
	case err != nil:
		HandleError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "\t")
	err = enc.Encode(&comments)
	if err != nil {
		HandleError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
}

func (a API) CreateComment(w http.ResponseWriter, r *http.Request) {
	var rid int
	{
		vars := mux.Vars(r)
		_, err := fmt.Sscanf(vars["rid"], "%d", &rid)
		if err != nil {
			HandleError(w, r, http.StatusBadRequest, err.Error())
			return
		}
	}

	var comment Comment
	{
		dec := json.NewDecoder(r.Body)
		defer r.Body.Close()
		err := dec.Decode(&comment)
		if err != nil {
			HandleError(w, r, http.StatusBadRequest, err.Error())
			return
		}
	}

	id, err := a.Repo.CreateComment(rid, comment)
	if err != nil {
		HandleError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	_, err = fmt.Fprintf(w, `{"id": %d}`, id)
	if err != nil {
		HandleError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
}

func (a API) ReadComment(w http.ResponseWriter, r *http.Request) {
	var rid, id int
	{
		vars := mux.Vars(r)
		_, err := fmt.Sscanf(vars["rid"], "%d", &rid)
		if err != nil {
			HandleError(w, r, http.StatusBadRequest, err.Error())
			return
		}

		_, err = fmt.Sscanf(vars["id"], "%d", &id)
		if err != nil {
			HandleError(w, r, http.StatusBadRequest, err.Error())
			return
		}
	}

	comment, err := a.Repo.ReadComment(rid, id)
	switch {
	case err == ReviewNotFound || err == CommentNotFound:
		HandleError(w, r, http.StatusBadRequest, err.Error())
		return
	case err != nil:
		HandleError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "\t")
	err = enc.Encode(&comment)
	if err != nil {
		HandleError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
}

func (a API) UpdateComment(w http.ResponseWriter, r *http.Request) {
	var rid, id int
	{
		vars := mux.Vars(r)
		_, err := fmt.Sscanf(vars["rid"], "%d", &rid)
		if err != nil {
			HandleError(w, r, http.StatusBadRequest, err.Error())
			return
		}

		_, err = fmt.Sscanf(vars["id"], "%d", &id)
		if err != nil {
			HandleError(w, r, http.StatusBadRequest, err.Error())
			return
		}
	}

	var comment Comment
	{
		dec := json.NewDecoder(r.Body)
		defer r.Body.Close()
		err := dec.Decode(&comment)
		if err != nil {
			HandleError(w, r, http.StatusBadRequest, err.Error())
			return
		}
	}

	err := a.Repo.UpdateComment(rid, id, comment)
	switch {
	case err == ReviewNotFound || err == CommentNotFound:
		HandleError(w, r, http.StatusBadRequest, err.Error())
		return
	case err != nil:
		HandleError(w, r, http.StatusBadRequest, err.Error())
		return
	}
}

func (a API) DeleteComment(w http.ResponseWriter, r *http.Request) {
	var rid, id int
	{
		vars := mux.Vars(r)
		_, err := fmt.Sscanf(vars["rid"], "%d", &rid)
		if err != nil {
			HandleError(w, r, http.StatusBadRequest, err.Error())
			return
		}

		_, err = fmt.Sscanf(vars["id"], "%d", &id)
		if err != nil {
			HandleError(w, r, http.StatusBadRequest, err.Error())
			return
		}
	}

	err := a.Repo.DeleteComment(rid, id)
	switch {
	case err == ReviewNotFound || err == CommentNotFound:
		HandleError(w, r, http.StatusBadRequest, err.Error())
		return
	case err != nil:
		HandleError(w, r, http.StatusBadRequest, err.Error())
		return
	}
}

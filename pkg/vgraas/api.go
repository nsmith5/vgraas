package vgraas

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nsmith5/vgraas/pkg/middleware"
)

// API implements the OpenAPI specification of vgraas.
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

// NewAPI returns an http.Handler that implements
// the OpenAPI specification for vgraas.
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

		Route{"Health", "GET", "/healthz", a.Health},
	}

	for _, route := range routes {
		a.Router.
			Methods(route.Methods).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}

	// Fall back for non-existant routers
	a.Router.NotFoundHandler = http.HandlerFunc(NotFound)

	return middleware.ContentType(a, "application/json; charset=UTF=8")
}

// HandleError sets the status code and writes a JSON object with
// error message for requests that have fallen on troubled times.
func HandleError(w http.ResponseWriter, r *http.Request, status int, err string) {
	w.WriteHeader(status)
	fmt.Fprintf(w, `{"err": "%s"}`, err)
}

// NotFound is a handy request handler for routes that don't exist.
func NotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, `{"err": "Route '%s' does not exist"}`, r.URL.Path)
}

// ReadReviews implements GET /reviews/
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

// CreateReview implements POST /reviews/
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

// ReadReview implements GET /reviews/{id}
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

// UpdateReview implements PUT /reviews/{id}
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

// DeleteReview implements DELETE /reviews/{id}
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

// ReadComments implements GET /reviews/{rid}/comments
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

// CreateComment implements POST /reviews/{rid}/comments
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

// ReadComment implements GET /reviews/{rid}/comments/{id}
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

// UpdateComment implements PUT /reviews/{rid}/comments/{id}
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

// DeleteComment implements DELETE /reviews/{rid}/comments/{id}
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

// Health implements a health monitoring endpoint at /healthz.
//
// Pop-quiz: Why is that 'z' always there? Good question. Anyways
// this endpoint is great for Kubernetes because you can use
// it for liveness and readiness probes.
func (a API) Health(w http.ResponseWriter, r *http.Request) {
}

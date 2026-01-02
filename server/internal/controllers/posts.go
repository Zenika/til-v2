package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid/v5"
	"github.com/zenika/tilv2back/internal/services"
	"github.com/zenika/tilv2back/internal/structures"
	"io"
	"net/http"
	"strings"
)

func GetPosts(w http.ResponseWriter, r *http.Request) {
	page := r.Context().Value("page").(structures.Pagination)

	var err error
	var items structures.Paginate[structures.Post]

	if r.URL.Query().Get("tags") != "" {
		items, err = services.GetPostsByTags(strings.Split(r.URL.Query().Get("tags"), ","), page)
	} else {
		items, err = services.GetPosts(page)
	}

	if err != nil {
		ErrorResponder(err, w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	marshaled, _ := json.Marshal(items)
	_, _ = w.Write(marshaled)
}

func GetPost(w http.ResponseWriter, r *http.Request) {
	postId := uuid.FromStringOrNil(chi.URLParam(r, "postId"))
	if postId == uuid.Nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	data, err := services.GetPostById(postId)
	if err != nil {
		ErrorResponder(err, w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	marshaled, _ := json.Marshal(data)
	_, _ = w.Write(marshaled)
}

func CreatePost(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(structures.User)

	payload, _ := io.ReadAll(r.Body)
	var post structures.Post
	err := json.Unmarshal(payload, &post)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id, err := services.CreatePost(post, user)
	if err != nil {
		ErrorResponder(err, w)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Location", fmt.Sprintf("/posts/%s", id.String()))
	w.Header().Set("X-Post-Id", id.String())
	w.WriteHeader(http.StatusCreated)

}

func DeletePost(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(structures.User)

	postId := uuid.FromStringOrNil(chi.URLParam(r, "postId"))
	if postId == uuid.Nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := services.DeletePostById(postId, user)
	if err != nil {
		ErrorResponder(err, w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func RealTimePostsUpdates(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")

	clientGone := r.Context().Done()
	channel := services.Broadcast.Subscribe()

	// Send message every time a channel have something to say.
	for {
		select {
		case <-clientGone:
			services.Broadcast.Unsubscribe(channel)
			return
		case a := <-channel:
			d, _ := json.Marshal(a.Data)
			_, _ = fmt.Fprintf(w, "event: %s\n", a.Kind)
			_, _ = fmt.Fprintf(w, "data: %s\n\n", d)
			w.(http.Flusher).Flush()
		}
	}
}

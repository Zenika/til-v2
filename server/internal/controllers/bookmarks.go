package controllers

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid/v5"
	"github.com/zenika/tilv2back/internal/services"
	"github.com/zenika/tilv2back/internal/structures"
	"net/http"
)

func GetBookmarks(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(structures.User)

	bookmarks, err := services.GetBookmarksForCurrentUser(user)
	if err != nil {
		ErrorResponder(err, w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if len(bookmarks) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	marshaled, _ := json.Marshal(bookmarks)
	_, _ = w.Write(marshaled)
}

func AddBookmark(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(structures.User)

	postId := uuid.FromStringOrNil(chi.URLParam(r, "postId"))
	if postId == uuid.Nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := services.AddBookmarkForCurrentUser(postId, user); err != nil {
		ErrorResponder(err, w)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}

func DeleteBookmark(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(structures.User)

	postId := uuid.FromStringOrNil(chi.URLParam(r, "postId"))
	if postId == uuid.Nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := services.DeleteBookmarkForCurrentUser(postId, user); err != nil {
		ErrorResponder(err, w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

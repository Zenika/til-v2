package controllers

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid/v5"
	"github.com/zenika/tilv2back/internal/services"
	"github.com/zenika/tilv2back/internal/structures"
	"io"
	"net/http"
)

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(structures.User)
	userId := getPathUserId(r)

	payload, _ := io.ReadAll(r.Body)
	var updatingUser structures.User
	err := json.Unmarshal(payload, &updatingUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	updatingUser.Id = userId

	err = services.UpdateUser(updatingUser, user)
	if err != nil {
		ErrorResponder(err, w)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(structures.User)
	userId := getPathUserId(r)

	err := services.DeleteUser(userId, user)
	if err != nil {
		ErrorResponder(err, w)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	page := r.Context().Value("page").(structures.Pagination)
	user := r.Context().Value("user").(structures.User)

	items, err := services.GetUsers(page, user)

	if err != nil {
		ErrorResponder(err, w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	marshaled, _ := json.Marshal(items)
	_, _ = w.Write(marshaled)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(structures.User)
	userId := getPathUserId(r)

	user, err := services.GetUser(userId, &user)
	if err != nil {
		ErrorResponder(err, w)
		return
	}

	data, _ := json.Marshal(user)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}

func RenewFeedKey(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(structures.User)

	userId := getPathUserId(r)
	if userId == uuid.Nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := services.RegenerateFeedKey(userId, user); err != nil {
		ErrorResponder(err, w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// getPathUserId is a stub to handle properly "self" as ID
func getPathUserId(r *http.Request) uuid.UUID {
	user := r.Context().Value("user").(structures.User)

	pathUserId := chi.URLParam(r, "userId")
	if pathUserId == "self" {
		return user.Id
	}

	return uuid.FromStringOrNil(chi.URLParam(r, "userId"))

}

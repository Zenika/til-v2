package controllers

import (
	"encoding/json"
	"github.com/zenika/tilv2back/internal/services"
	"net/http"
)

func GetTags(w http.ResponseWriter, r *http.Request) {
	tags, err := services.GetAllTags()
	if err != nil {
		ErrorResponder(err, w)
		return
	}

	marshaled, _ := json.Marshal(tags)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(marshaled)
}

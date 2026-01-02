package controllers

import (
	"encoding/json"
	"github.com/zenika/tilv2back/internal/configuration"
	"net/http"
)

func Configuration(w http.ResponseWriter, r *http.Request) {
	config := struct {
		GoogleId string `json:"google_id"`
	}{
		configuration.Configuration.Google.ClientID,
	}

	data, _ := json.Marshal(config)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}

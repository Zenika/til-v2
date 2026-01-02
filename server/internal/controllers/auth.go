package controllers

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/zenika/tilv2back/internal/configuration"
	"github.com/zenika/tilv2back/internal/services"
	"github.com/zenika/tilv2back/internal/structures"
	"net/http"
	"net/url"
	"time"
)

func RunAuth(w http.ResponseWriter, r *http.Request) {

	authCode := r.URL.Query().Get("code")
	redirectUriParsed, _ := url.ParseQuery(r.URL.Query().Get("state"))
	redirectUri := redirectUriParsed.Get("redirect_uri")

	if authCode == "" || redirectUri == "" {
		fmt.Println(authCode, redirectUri)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userId, err := services.GoogleAuth(authCode, redirectUri)
	if err != nil {
		ErrorResponder(err, w)
		return
	}

	token, err := services.CreateToken(userId)
	if err != nil {
		ErrorResponder(err, w)
		return
	}

	w.Header().Set("Authorization", token)
	w.WriteHeader(204)

}

func RenewToken(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(structures.User)
	claims := r.Context().Value("tokenClaims").(jwt.MapClaims)

	if d, e := claims.GetExpirationTime(); e == nil && (d.Before(time.Now().Add(time.Hour*4)) || configuration.Configuration.Debug) {
		token, err := services.CreateToken(user.Id)
		if err != nil {
			ErrorResponder(err, w)
			return
		}

		w.Header().Set("Authorization", token)
		w.WriteHeader(204)
	} else {
		w.Header().Set("Authorization", r.Header.Get("Authorization"))
		w.WriteHeader(204)
	}
}

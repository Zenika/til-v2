package services

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/gofrs/uuid/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/zenika/tilv2back/internal/configuration"
	"github.com/zenika/tilv2back/internal/custom_errors"
	"github.com/zenika/tilv2back/internal/repository"
	"github.com/zenika/tilv2back/internal/structures"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func GoogleAuth(authCode string, redirectUri string) (uuid.UUID, error) {
	data := url.Values{}
	data.Add("code", authCode)
	data.Add("client_id", configuration.Configuration.Google.ClientID)
	data.Add("client_secret", configuration.Configuration.Google.ClientSecret)
	data.Add("redirect_uri", redirectUri)
	data.Add("grant_type", "authorization_code")
	data.Add("access_type", "offline")

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	
	req, err := http.Post(configuration.Configuration.Google.TokenEndpoint, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
	if err != nil {
		return uuid.Nil, custom_errors.NewInternalServerError("Unable to contact Google Auth", err)
	}

	result, err := io.ReadAll(req.Body)
	if err != nil {
		return uuid.Nil, custom_errors.NewInternalServerError("Unable to read Google Auth result", err)
	}

	var tokenData map[string]interface{}
	err = json.Unmarshal(result, &tokenData)
	if err != nil {
		return uuid.Nil, custom_errors.NewInternalServerError("Unable to parse Google Auth result", err)
	}

	if _, ok := tokenData["id_token"]; !ok {
		return uuid.Nil, custom_errors.NewInternalServerError(fmt.Sprintf("Error returned by Google; %s", tokenData["error"]), err)
	}

	token, err := jwt.Parse(tokenData["id_token"].(string), nil)
	if err != nil && err.Error() != "token is unverifiable: no keyfunc was provided" {
		return uuid.Nil, custom_errors.NewInternalServerError("Unable to parse JWT token", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if user, err := repository.GetUserByGoogleId(claims["sub"].(string)); err == nil {
			return user.Id, nil
		} else {
			userId, err := repository.CreateUser(structures.User{
				DisplayName: claims["name"].(string),
				GoogleId:    claims["sub"].(string),
			})

			if err != nil {
				return uuid.Nil, custom_errors.NewPermissionDeniedError("Failed to create user", err)
			}

			return userId, nil
		}

	} else {
		return uuid.Nil, custom_errors.NewPermissionDeniedError("Invalid token", nil)
	}
}

func CreateToken(userId uuid.UUID) (string, error) {
	user, err := repository.GetUser(userId, false, true)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"name":     user.DisplayName,
		"id":       userId,
		"is_admin": user.IsAdmin,
		"iat":      jwt.NewNumericDate(time.Now()),
		"nbf":      jwt.NewNumericDate(time.Now()),
		"exp":      jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		"iss":      "univerzity",
	})

	return token.SignedString([]byte(configuration.Configuration.JwtSecret))
}

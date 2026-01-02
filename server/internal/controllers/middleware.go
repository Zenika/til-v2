package controllers

import (
	"context"
	"errors"
	"github.com/gofrs/uuid/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/zenika/tilv2back/internal/configuration"
	"github.com/zenika/tilv2back/internal/custom_errors"
	"github.com/zenika/tilv2back/internal/services"
	"github.com/zenika/tilv2back/internal/structures"
	"net/http"
	"strconv"
	"strings"
)

func HandleCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, HEAD, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization")
		w.Header().Set("Access-Control-Expose-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization, Location, X-Post-Id")
		w.Header().Set("Access-Control-Max-Age", "86400")

		if r.Method == "OPTIONS" {
			return
		}

		next.ServeHTTP(w, r)
	})
}

func PaginationMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		page := structures.Pagination{
			Page:    0,
			PerPage: 20,
		}

		if pageNumber, err := strconv.Atoi(r.URL.Query().Get("page")); err == nil && pageNumber > 0 {
			page.Page = pageNumber
		}

		if pageItems, err := strconv.Atoi(r.URL.Query().Get("size")); err == nil && pageItems > 0 && pageItems < 100 {
			page.PerPage = pageItems
		}

		ctx := context.WithValue(r.Context(), "page", page)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func AuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Custom authentication for RSS and browser extension permanent authorization
		// Only POST /posts & GET /rss are allowed to use this method, please refer to documentation for further details
		if ((r.URL.Path == "/api/posts" && r.Method == http.MethodPost) || r.URL.Path == "/api/rss") && r.URL.Query().Get("key") != "" {
			user, err := services.GetUserByFeedKey(r.URL.Query().Get("key"))
			// If we fail, gracefully continue to standard authentication
			if err == nil {
				ctx := context.WithValue(r.Context(), "user", user)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
		}

		if r.Header.Get("Authorization") == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(r.Header.Get("Authorization"), func(token *jwt.Token) (interface{}, error) {
			return []byte(configuration.Configuration.JwtSecret), nil
		}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {

			user, err := uuid.FromString(claims["id"].(string))
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			realUser, err := services.GetUser(user, nil)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), "tokenClaims", claims)
			ctx = context.WithValue(ctx, "user", realUser)
			next.ServeHTTP(w, r.WithContext(ctx))

		} else {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

	})
}

func ErrorResponder(err error, w http.ResponseWriter) {
	var customError *custom_errors.CustomError
	if errors.As(err, &customError) {
		if customError.BaseError != nil {
			configuration.Logger.Error(customError.BaseError.Error())
		}
		w.WriteHeader(customError.ErrorCode)
		_, _ = w.Write([]byte(customError.Content))
	} else {
		if strings.HasPrefix(err.Error(), "UNIQUE constraint failed") {
			w.WriteHeader(http.StatusConflict)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}

	configuration.Logger.Error(err.Error())
}

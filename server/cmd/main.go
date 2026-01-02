package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/zenika/tilv2back/internal/configuration"
	"github.com/zenika/tilv2back/internal/controllers"
	"github.com/zenika/tilv2back/internal/repository"
	"github.com/zenika/tilv2back/internal/services"
	"github.com/zenika/tilv2back/internal/structures"
	"io"
	"net/http"
	"os"
	"path"
	"time"
)

func main() {

	go heartbeat()

	r := chi.NewRouter()

	r.Route("/api", func(r chi.Router) {
		r.Use(middleware.Logger)
		r.Use(controllers.HandleCORS)
		r.Use(controllers.PaginationMiddleware)

		// Unauthenticated endpoints
		r.Get("/auth", controllers.RunAuth)
		r.Get("/configuration", controllers.Configuration)

		// Authenticated endpoints
		r.Group(func(r chi.Router) {
			r.Use(controllers.AuthenticationMiddleware)

			r.Route("/posts", func(r chi.Router) {
				r.Get("/stream", controllers.RealTimePostsUpdates)
				r.Get("/", controllers.GetPosts)
				r.Post("/", controllers.CreatePost)
				r.Delete("/{postId}", controllers.DeletePost)
				r.Get("/{postId}", controllers.GetPost)
			})

			r.Route("/users", func(r chi.Router) {
				r.Get("/", controllers.GetUsers)
				r.Delete("/{userId}", controllers.DeleteUser)
				r.Get("/{userId}", controllers.GetUser)
				r.Put("/{userId}", controllers.UpdateUser)
				r.Put("/{userId}/renew", controllers.RenewFeedKey)
			})

			r.Route("/bookmarks", func(r chi.Router) {
				r.Get("/", controllers.GetBookmarks)
				r.Put("/{postId}", controllers.AddBookmark)
				r.Delete("/{postId}", controllers.DeleteBookmark)
			})

			r.Get("/rss", controllers.GetRSSFeed)
			r.Get("/tags", controllers.GetTags)
			r.Get("/renew", controllers.RenewToken)
		})
	})

	if configuration.Configuration.UseEmbeddedFrontend {
		r.HandleFunc("/*", func(w http.ResponseWriter, r *http.Request) {
			fs := http.Dir("/static")
			fileServer := http.FileServer(fs)
			_, err := fs.Open(path.Clean(r.URL.Path))
			if os.IsNotExist(err) {
				d, err := fs.Open("/index.html")
				if err != nil {
					configuration.Logger.Error("Unable to open fallback file", err)
					return
				}
				defer d.Close()
				di, _ := io.ReadAll(d)
				_, _ = w.Write(di)
				return
			}
			fileServer.ServeHTTP(w, r)
		})
	}

	if users, err := repository.GetUsers(structures.Pagination{}, false); err == nil && users.TotalItems == 0 && configuration.Configuration.DefaultAdmin != "" {
		_, err := repository.CreateUser(structures.User{
			DisplayName: "Admin",
			GoogleId:    configuration.Configuration.DefaultAdmin,
			IsAdmin:     true,
		})
		if err != nil {
			configuration.Logger.Error("Unable to register default admin.", err.Error())
		}
	}

	configuration.Logger.Info(fmt.Sprintf("Server is now listening on port %d", configuration.Configuration.ServerPort))
	_ = http.ListenAndServe(fmt.Sprintf(":%d", configuration.Configuration.ServerPort), r)
}

func heartbeat() {
	for {
		time.Sleep(15 * time.Second)
		services.SendHeartbeatToClients()
	}
}

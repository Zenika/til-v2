package controllers

import (
	"fmt"
	"github.com/gorilla/feeds"
	"github.com/zenika/tilv2back/internal/services"
	"github.com/zenika/tilv2back/internal/structures"
	"net/http"
	"strings"
)

func GetRSSFeed(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(structures.User)

	if r.URL.Query().Get("tags") != "" {
		user.AutomaticTagsFilter = strings.Split(r.URL.Query().Get("tags"), ",")
	}

	posts, err := services.GetPostsByTags(user.AutomaticTagsFilter, structures.Pagination{
		Page:    0,
		PerPage: 20,
	})

	if err != nil {
		ErrorResponder(err, w)
		return
	}

	feed := &feeds.Feed{
		Title:       "TIL - Zenika custom feed",
		Link:        &feeds.Link{Href: "https://til.znk.io/rss"},
		Description: fmt.Sprintf("Customized feed for %s", user.DisplayName),
		Author:      &feeds.Author{Name: "Zenika"},
		Items:       []*feeds.Item{},
	}

	for _, k := range posts.Items {
		feed.Items = append(feed.Items, &feeds.Item{
			Title:       k.Title,
			Link:        &feeds.Link{Href: k.Link},
			Author:      &feeds.Author{Name: k.User.DisplayName},
			Description: k.Content,
			Id:          k.Id.String(),
			Created:     k.CreationDate,
			Content:     k.Content,
		})
	}

	generatedFeed := &feeds.Rss{Feed: feed}
	rss, _ := generatedFeed.ToRss()

	w.Header().Set("Content-Type", "application/rss+xml")
	_, _ = w.Write([]byte(rss))

}

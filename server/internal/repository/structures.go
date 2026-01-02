package repository

import (
	"github.com/gofrs/uuid/v5"
	"time"
)

type userDatabase struct {
	Id            uuid.UUID `sql:"id"`
	DisplayName   string    `sql:"display_name"`
	GoogleId      string    `sql:"google_id"`
	IsAdmin       bool      `sql:"is_admin"`
	DefaultFilter string    `sql:"default_filter"`
	FeedKey       string    `sql:"feed_key"`
}

type tag struct {
	Tag string `sql:"tag"`
}

type postDatabase struct {
	Id           uuid.UUID `sql:"id"`
	CreationDate time.Time `sql:"creation_date"`
	UserId       uuid.UUID `sql:"user_id"`
	Title        string    `sql:"title"`
	Content      string    `sql:"content"`
	Link         string    `sql:"link"`
	Image        string    `sql:"image"`
	Icon         string    `sql:"icon"`
}

package structures

import (
	"github.com/gofrs/uuid/v5"
	"time"
)

type Post struct {
	Id           uuid.UUID `json:"id"`
	CreationDate time.Time `json:"creation_date"`
	User         User      `json:"user,omitempty"`
	Title        string    `json:"title"`
	Content      string    `json:"content,omitempty"`
	Link         string    `json:"link,omitempty"`
	Tags         []string  `json:"tags,omitempty"`
	Image        string    `json:"image,omitempty"`
	Icon         string    `json:"icon,omitempty"`
}

type User struct {
	Id                  uuid.UUID `json:"id"`
	DisplayName         string    `json:"display_name"`
	GoogleId            string    `json:"google_id,omitempty"`
	IsAdmin             bool      `json:"is_admin"`
	FavoritePosts       []Post    `json:"favorite_posts,omitempty"`
	AutomaticTagsFilter []string  `json:"automatic_tags_filter,omitempty"`
	FeedKey             string    `json:"feed_key,omitempty"`
}

type Paginate[T Post | User] struct {
	Items        []T `json:"items"`
	TotalItems   int `json:"total_items"`
	ItemsPerPage int `json:"items_per_page"`
	CurrentPage  int `json:"current_page"`
	TotalPages   int `json:"total_pages"`
}

type Pagination struct {
	Page    int `default:"0"`
	PerPage int `default:"20"`
}

type Configuration struct {
	Debug               bool                `mapstructure:"DEBUG"`
	DatabaseFileName    string              `mapstructure:"DATABASE_FILE_NAME"`
	ServerPort          int                 `mapstructure:"SERVER_PORT"`
	JwtSecret           string              `mapstructure:"JWT_SECRET"`
	Google              ConfigurationGoogle `mapstructure:"GOOGLE"`
	DefaultAdmin        string              `mapstructure:"DEFAULT_ADMIN"`
	UseEmbeddedFrontend bool                `mapstructure:"USE_EMBEDDED_FRONTEND"`
}

type ConfigurationGoogle struct {
	ClientID      string `mapstructure:"CLIENT_ID"`
	ClientSecret  string `mapstructure:"CLIENT_SECRET"`
	TokenEndpoint string `mapstructure:"TOKEN_ENDPOINT"`
}

type Event struct {
	Kind string `json:"kind"`
	Data Post   `json:"data"`
}

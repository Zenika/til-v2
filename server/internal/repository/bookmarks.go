package repository

import (
	"github.com/gofrs/uuid/v5"
	"github.com/zenika/tilv2back/internal/database"
	"github.com/zenika/tilv2back/internal/structures"
)

func GetBookmarksByUser(userId uuid.UUID) (result []structures.Post, err error) {
	posts, err := database.Query[postDatabase]("SELECT * FROM `user_favorites_posts` INNER JOIN `posts` ON `user_favorites_posts`.`post_id` = `posts`.`id` WHERE `user_favorites_posts`.`user_id` = ? ORDER BY `posts`.`creation_date` DESC", userId.String())
	if err != nil {
		return
	}

	for _, post := range posts {
		result = append(result, toServicePost(post, true))
	}

	return
}

func AddBookmarkForUser(userId uuid.UUID, postId uuid.UUID) error {
	return database.Exec("INSERT INTO `user_favorites_posts` (user_id, post_id) VALUES (?, ?)", userId, postId)
}

func DeleteBookmarkForUser(userId uuid.UUID, postId uuid.UUID) error {
	return database.Exec("DELETE FROM `user_favorites_posts` WHERE user_id = ? AND post_id = ?", userId, postId)
}

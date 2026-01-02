package services

import (
	"github.com/gofrs/uuid/v5"
	"github.com/zenika/tilv2back/internal/repository"
	"github.com/zenika/tilv2back/internal/structures"
)

func GetBookmarksForCurrentUser(currentUser structures.User) ([]structures.Post, error) {
	return repository.GetBookmarksByUser(currentUser.Id)
}

func AddBookmarkForCurrentUser(postId uuid.UUID, currentUser structures.User) error {
	// Check if post exists
	_, err := repository.GetPostById(postId)
	if err != nil {
		return err
	}

	return repository.AddBookmarkForUser(currentUser.Id, postId)
}

func DeleteBookmarkForCurrentUser(postId uuid.UUID, currentUser structures.User) error {
	return repository.DeleteBookmarkForUser(currentUser.Id, postId)
}

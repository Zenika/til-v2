package services

import (
	"github.com/gofrs/uuid/v5"
	"github.com/zenika/tilv2back/internal/custom_errors"
	"github.com/zenika/tilv2back/internal/repository"
	"github.com/zenika/tilv2back/internal/structures"
	"strings"
)

func DeleteUser(userId uuid.UUID, currentUser structures.User) error {
	if !currentUser.IsAdmin {
		return custom_errors.NewPermissionDeniedError("You are not allowed to delete this user", nil)
	}

	return repository.DeleteUser(userId)
}

func GetUser(userId uuid.UUID, currentUser *structures.User) (structures.User, error) {
	includeSensitive := false
	if currentUser == nil || currentUser.IsAdmin || currentUser.Id == userId {
		includeSensitive = true
	}

	user, err := repository.GetUser(userId, true, includeSensitive)
	if err != nil {
		return structures.User{}, err
	}

	return user, nil
}

func GetUsers(pagination structures.Pagination, currentUser structures.User) (structures.Paginate[structures.User], error) {
	users, err := repository.GetUsers(pagination, currentUser.IsAdmin)
	if err != nil {
		return structures.Paginate[structures.User]{}, err
	}

	return users, nil
}

func UpdateUser(userToUpdate structures.User, currentUser structures.User) error {
	if !currentUser.IsAdmin && userToUpdate.Id != currentUser.Id {
		return custom_errors.NewPermissionDeniedError("You are not allowed to update this user", nil)
	}

	if currentUser.IsAdmin {
		_ = repository.SetAdminState(userToUpdate.Id, userToUpdate.IsAdmin)
	}

	_ = repository.SetFilter(userToUpdate.Id, strings.Join(userToUpdate.AutomaticTagsFilter, ","))
	_ = repository.SetDisplayName(userToUpdate.Id, userToUpdate.DisplayName)

	return nil
}

func GetUserByFeedKey(feedKey string) (structures.User, error) {
	return repository.GetUserByFeedKey(feedKey)
}

func RegenerateFeedKey(userId uuid.UUID, currentUser structures.User) error {
	if !currentUser.IsAdmin && userId != currentUser.Id {
		return custom_errors.NewPermissionDeniedError("You are not allowed to renew feed key for this user", nil)
	}

	return repository.RegenerateFeedKey(userId)
}

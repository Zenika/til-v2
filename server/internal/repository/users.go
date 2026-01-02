package repository

import (
	"github.com/gofrs/uuid/v5"
	"github.com/zenika/tilv2back/internal/database"
	"github.com/zenika/tilv2back/internal/structures"
	"strings"
)

func CreateUser(user structures.User) (uuid.UUID, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return uuid.Nil, err
	}

	err = database.Exec("INSERT INTO `users`(`id`, `display_name`, `is_admin`, `google_id`, `feed_key`) VALUES(?,?,?,?,?);", id.String(), user.DisplayName, user.IsAdmin, user.GoogleId, generateRandomString(30))
	return id, err
}

func RegenerateFeedKey(userId uuid.UUID) error {
	return database.Exec("UPDATE `users` SET `feed_key`=? WHERE id=?;", generateRandomString(30), userId.String())
}

func SetAdminState(userId uuid.UUID, isAdmin bool) error {
	return database.Exec("UPDATE `users` SET `is_admin`=? WHERE id=?;", isAdmin, userId.String())
}

func SetFilter(userId uuid.UUID, filter string) error {
	return database.Exec("UPDATE `users` SET `default_filter`=? WHERE id=?;", filter, userId.String())
}

func SetDisplayName(userId uuid.UUID, displayName string) error {
	return database.Exec("UPDATE `users` SET `display_name`=? WHERE id=?;", displayName, userId.String())
}

func DeleteUser(user uuid.UUID) error {
	return database.Exec("DELETE FROM `users` WHERE `id` = ?", user.String())
}

func GetUser(user uuid.UUID, populate bool, includeSensitive bool) (structures.User, error) {
	userObject, err := database.QueryOne[userDatabase]("SELECT * FROM `users` WHERE `id` = ?", user.String())
	if err != nil {
		return structures.User{}, err
	}
	return toServiceUser(userObject, populate, includeSensitive), nil
}

func GetUserByGoogleId(googleId string) (structures.User, error) {
	userObject, err := database.QueryOne[userDatabase]("SELECT * FROM `users` WHERE `google_id` = ?", googleId)
	if err != nil {
		return structures.User{}, err
	}
	return toServiceUser(userObject, false, true), nil
}

func GetUserByFeedKey(feedKey string) (structures.User, error) {
	userObject, err := database.QueryOne[userDatabase]("SELECT * FROM `users` WHERE `feed_key` = ?", feedKey)
	if err != nil {
		return structures.User{}, err
	}
	return toServiceUser(userObject, false, true), nil
}

func GetUsers(page structures.Pagination, includeSensitive bool) (result structures.Paginate[structures.User], err error) {
	users, err := database.Query[userDatabase]("SELECT * FROM `users` LIMIT ?,?", page.Page*page.PerPage, page.PerPage)
	if err != nil {
		return
	}

	result, err = createPage[structures.User]("users", page, nil)
	if err != nil {
		return
	}

	for _, user := range users {
		result.Items = append(result.Items, toServiceUser(user, true, includeSensitive))
	}

	return
}

func toServiceUser(user userDatabase, populate bool, includeSensitive bool) structures.User {
	newUser := structures.User{
		Id:          user.Id,
		DisplayName: user.DisplayName,
	}

	if includeSensitive {
		newUser.GoogleId = user.GoogleId
		newUser.IsAdmin = user.IsAdmin
		if user.DefaultFilter != "" {
			newUser.AutomaticTagsFilter = strings.Split(user.DefaultFilter, ",")
		}
		newUser.FeedKey = user.FeedKey
	}

	if populate {
		newUser.FavoritePosts, _ = GetBookmarksByUser(user.Id)
	}
	return newUser
}

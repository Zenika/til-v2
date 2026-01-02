package repository

import (
	"github.com/gofrs/uuid/v5"
	"github.com/zenika/tilv2back/internal/database"
	"github.com/zenika/tilv2back/internal/structures"
	"strings"
)

func GetPosts(page structures.Pagination) (result structures.Paginate[structures.Post], err error) {
	posts, err := database.Query[postDatabase]("SELECT * FROM `posts`  ORDER BY `creation_date` DESC LIMIT ?,?", page.Page*page.PerPage, page.PerPage)
	if err != nil {
		return
	}

	result, err = createPage[structures.Post]("posts", page, nil)
	if err != nil {
		return
	}

	result.Items = []structures.Post{}

	for _, post := range posts {
		result.Items = append(result.Items, toServicePost(post, true))
	}

	return
}

func SetArticleImage(id uuid.UUID, postImage string, favicon string) error {
	return database.Exec("UPDATE `posts` SET `image` = ?, `icon` = ? WHERE `id` = ?", postImage, favicon, id.String())
}

func GetPostsByTags(tags []string, page structures.Pagination) (result structures.Paginate[structures.Post], err error) {
	if len(tags) == 0 {
		return GetPosts(page)
	}

	req := "SELECT * FROM `posts` WHERE posts.id IN( SELECT post_id FROM post_tags WHERE tag IN (?" + strings.Repeat(",?", len(tags)-1) + ")) ORDER BY `posts`.`creation_date` DESC LIMIT ?,?"
	reqCount := "SELECT COUNT(*) FROM `posts` WHERE posts.id IN( SELECT post_id FROM post_tags WHERE tag IN (?" + strings.Repeat(",?", len(tags)-1) + "))"

	var args []interface{}

	for _, k := range tags {
		args = append(args, k)
	}

	count := 0

	if rows, err := database.Database.Query(reqCount, args...); err == nil {
		rows.Next()
		rows.Scan(&count)
		rows.Close()
	}

	args = append(args, page.Page*page.PerPage)
	args = append(args, page.PerPage)

	posts, err := database.Query[postDatabase](req, args...)
	if err != nil {
		return
	}

	result, err = createPage[structures.Post]("posts", page, &count)
	if err != nil {
		return
	}

	result.Items = []structures.Post{}

	for _, post := range posts {
		result.Items = append(result.Items, toServicePost(post, true))
	}

	return
}

func GetPostById(id uuid.UUID) (structures.Post, error) {
	post, err := database.QueryOne[postDatabase]("SELECT * FROM `posts` WHERE `id` = ?", id)
	if err != nil {
		return structures.Post{}, err
	}

	return toServicePost(post, true), nil
}

func DeletePostById(id uuid.UUID) (err error) {
	// Will fail if we delete a tag still in use but not a big deal, so ignoring it.
	_ = database.Exec("DELETE FROM `post_tags` WHERE `post_id` = ?;", id.String())
	_ = database.Exec("DELETE FROM `user_favorites_posts` WHERE `post_id` = ?;", id.String())

	return database.Exec("DELETE FROM `posts` WHERE `id` = ?;", id.String())
}

func CreatePost(post structures.Post) (id uuid.UUID, err error) {
	id, err = uuid.NewV4()
	if err != nil {
		return uuid.Nil, err
	}
	err = database.Exec("INSERT INTO posts (id, user_id, title, link, content) VALUES(?,?,?,?,?)", id.String(), post.User.Id, post.Title, post.Link, post.Content)
	if err != nil {
		return uuid.Nil, err
	}

	for _, tag := range post.Tags {
		_ = database.Exec("INSERT INTO post_tags (post_id, tag) VALUES (?,?)", id.String(), tag)
	}

	return id, err
}

func GetAllTags() (tagList []string, err error) {
	tags, err := database.Query[tag]("SELECT DISTINCT `tag` FROM `post_tags`")
	if err != nil {
		return nil, err
	}
	for _, t := range tags {
		tagList = append(tagList, t.Tag)
	}
	return
}

func getTagsByPostId(id uuid.UUID) (stringTags []string, err error) {

	tags, err := database.Query[tag]("SELECT `tag` FROM `post_tags` WHERE `post_id` = ?", id.String())
	if err != nil {
		return []string{}, err
	}

	for i := range tags {
		stringTags = append(stringTags, tags[i].Tag)
	}

	return stringTags, nil
}

func toServicePost(post postDatabase, populate bool) structures.Post {
	newPost := structures.Post{
		Id:           post.Id,
		CreationDate: post.CreationDate,
		User: structures.User{
			Id: post.UserId,
		},
		Title:   post.Title,
		Content: post.Content,
		Link:    post.Link,
		Tags:    []string{},
		Image:   post.Image,
		Icon:    post.Icon,
	}

	if populate {
		newPost.Tags, _ = getTagsByPostId(post.Id)
		newPost.User, _ = GetUser(post.UserId, false, false)
	}

	return newPost
}

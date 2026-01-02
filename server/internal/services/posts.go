package services

import (
	"github.com/gofrs/uuid/v5"
	"github.com/zenika/tilv2back/internal/configuration"
	"github.com/zenika/tilv2back/internal/custom_errors"
	"github.com/zenika/tilv2back/internal/repository"
	"github.com/zenika/tilv2back/internal/structures"
	"io"
	"net/http"
	"regexp"
	"slices"
	"strings"
)

func CreatePost(post structures.Post, currentUser structures.User) (uuid.UUID, error) {
	if len(post.Title) == 0 {
		return uuid.Nil, custom_errors.NewInvalidArgumentError("Title could not be empty", nil)
	}

	if len(post.Content) == 0 && len(post.Link) == 0 {
		return uuid.Nil, custom_errors.NewInvalidArgumentError("You must set at least one of Content or Link", nil)
	}

	if len(post.Link) != 0 && !strings.HasPrefix(post.Link, "http://") && !strings.HasPrefix(post.Link, "https://") {
		return uuid.Nil, custom_errors.NewInvalidArgumentError("Link is not an URL", nil)
	}

	if len(post.Content) > 180 {
		return uuid.Nil, custom_errors.NewInvalidArgumentError("Content too long (180 characters max)", nil)
	}

	if !slices.ContainsFunc(post.Tags, func(d string) bool { return strings.HasPrefix(d, "lang:") }) {
		return uuid.Nil, custom_errors.NewInvalidArgumentError("Language not supported", nil)
	}

	post.User.Id = currentUser.Id

	// Tags to lowercase
	for i := range post.Tags {
		post.Tags[i] = strings.ToLower(post.Tags[i])
	}

	res, err := repository.CreatePost(post)
	if err != nil {
		return uuid.Nil, err
	}
	if len(post.Link) != 0 {
		go setImageArticleAndIcon(res, post.Link)
	}

	go sendEvent(res, "created")

	return res, nil
}

func DeletePostById(id uuid.UUID, currentUser structures.User) error {
	// Check if post exists
	p, err := repository.GetPostById(id)
	if err != nil {
		return err
	}

	if !currentUser.IsAdmin && currentUser.Id.String() != p.User.Id.String() {
		return custom_errors.NewPermissionDeniedError("You are not allowed to delete this post", nil)
	}

	return repository.DeletePostById(id)
}

func GetPostById(id uuid.UUID) (structures.Post, error) {
	return repository.GetPostById(id)
}

func GetPostsByTags(tags []string, pagination structures.Pagination) (structures.Paginate[structures.Post], error) {

	return repository.GetPostsByTags(tags, pagination)
}

func GetPosts(pagination structures.Pagination) (structures.Paginate[structures.Post], error) {
	return repository.GetPosts(pagination)
}

func GetAllTags() ([]string, error) {
	return repository.GetAllTags()
}

func setImageArticleAndIcon(id uuid.UUID, url string) {
	req, err := http.Get(url)
	if err != nil {
		configuration.Logger.Warn("Unable to retrieve image and icon from URL %s: %v", url, err)
		return
	}

	pageBody, err := io.ReadAll(req.Body)
	if err != nil {
		configuration.Logger.Warn("Unable to read content received from URL %s: %v", url, err)
		return
	}

	articleImage := ""
	favicon := ""

	// Look for main article image
	// Option 1 : look with og:image meta properties, as it is the standard way to do
	regex := regexp.MustCompile("(?ms)<meta([^>])*property=\"og:image\"([^>])*>")
	if find := regex.Find(pageBody); len(find) != 0 {
		regex = regexp.MustCompile("content=['\"]([^\"]*)['\"]")
		articleImage = string(regex.FindSubmatch(regex.Find(find))[1])
	} else {
		// Option 2 : If there is no meta, look for the first image on the website after a heading text.
		regex = regexp.MustCompile("(?ms)<(h1|h2|h3)(.*)")
		find = regex.Find(pageBody)
		if find != nil {
			regex = regexp.MustCompile("(?ms)<img([^>]*)>")
			find = regex.Find(find)
			if find != nil {
				regex = regexp.MustCompile("src=['\"]([^\"]*)['\"]")
				articleImage = string(regex.FindSubmatch(regex.Find(find))[1])
			}
		}
	}

	// Sanitize to ensure the image is a full URL address
	articleImage = sanitizeUrl(articleImage, req)

	// Lookup for favicon in meta, and sanitize
	regex = regexp.MustCompile("(?ms)<link([^>])*rel=\"icon\"([^>])*>")
	if find := regex.Find(pageBody); len(find) != 0 {
		regex = regexp.MustCompile("href=['\"]([^\"]*)['\"]")
		favicon = string(regex.FindSubmatch(regex.Find(find))[1])
		favicon = sanitizeUrl(favicon, req)
	}

	err = repository.SetArticleImage(id, articleImage, favicon)
	if err == nil {
		go sendEvent(id, "updated")
	}
}

func sanitizeUrl(url string, res *http.Response) string {
	if url != "" && !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		if strings.HasPrefix(url, "/") {
			url = res.Request.URL.Scheme + "://" + res.Request.URL.Host + url
		} else {
			url = res.Request.URL.Scheme + "://" + res.Request.URL.Host + "/" + url
		}
	}
	return url
}

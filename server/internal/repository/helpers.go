package repository

import (
	"github.com/zenika/tilv2back/internal/database"
	"github.com/zenika/tilv2back/internal/structures"
	"math"
	"math/rand"
	"time"
)

func createPage[T structures.Post | structures.User](tableName string, pagination structures.Pagination, rowNumber *int) (page structures.Paginate[T], err error) {
	var nb int

	if rowNumber == nil {
		nb, err = database.CountTable(tableName)
		rowNumber = &nb
		if err != nil {
			return
		}
	}

	page.TotalItems = *rowNumber
	page.ItemsPerPage = pagination.PerPage
	page.CurrentPage = pagination.Page
	page.TotalPages = int(math.Ceil(float64(*rowNumber) / float64(pagination.PerPage)))

	return
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

func generateRandomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

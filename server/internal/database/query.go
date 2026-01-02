package database

import (
	"fmt"
	"github.com/kisielk/sqlstruct"
	"github.com/zenika/tilv2back/internal/custom_errors"
)

func CountTable(tableName string) (int, error) {
	rows, err := Database.Query(fmt.Sprintf("SELECT COUNT(*) AS nb FROM %s;", tableName))
	if err != nil {
		return 0, err
	}

	defer rows.Close()

	for rows.Next() {
		var nb int
		if err := rows.Scan(&nb); err != nil {
			return 0, err
		}
		return nb, nil
	}
	return 0, nil
}

func Query[T any](request string, args ...any) ([]T, error) {
	rows, err := Database.Query(request, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var resultArray []T

	for rows.Next() {
		var result T
		err := sqlstruct.Scan(&result, rows)
		if err != nil {
			return nil, err
		}
		resultArray = append(resultArray, result)
	}

	return resultArray, nil
}

func QueryOne[T any](request string, args ...any) (T, error) {
	rows, err := Database.Query(request, args...)
	if err != nil {
		return *new(T), err
	}
	defer rows.Close()

	for rows.Next() {
		var result T
		err := sqlstruct.Scan(&result, rows)
		if err != nil {
			return *new(T), err
		}
		return result, nil
	}

	return *new(T), custom_errors.NewNotFoundError("Not found", nil)
}

func Exec(request string, args ...any) error {
	_, err := Database.Exec(request, args...)
	return err
}

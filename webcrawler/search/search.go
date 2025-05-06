package search

import (
	"database/sql"
	"fmt"
	"strings"
)

func SearchLinks(query string) ([]string, error) {
	// Простой LIKE-поиск
	keywords := strings.Split(query, " ")

	sqlQuery := `SELECT url FROM links WHERE `
	var conditions []string
	var args []interface{}

	for i, kw := range keywords {
		conditions = append(conditions, fmt.Sprintf("url ILIKE $%d", i+1))
		args = append(args, "%"+kw+"%")
	}

	sqlQuery += strings.Join(conditions, " OR ")
	sqlQuery += " LIMIT 10"

	rows, err := db.Query(sqlQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []string
	for rows.Next() {
		var url string
		if err := rows.Scan(&url); err == nil {
			results = append(results, url)
		}
	}

	return results, nil
}

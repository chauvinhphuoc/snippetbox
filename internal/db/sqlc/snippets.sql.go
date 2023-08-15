// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: snippets.sql

package sqlc

import (
	"context"
)

const createSnippet = `-- name: CreateSnippet :one
INSERT INTO snippets (title, content, created_at, expires)
VALUES ($1, $2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP + MAKE_INTERVAL(DAYS => $3::int)) RETURNING id
`

type CreateSnippetParams struct {
	Title    string `json:"title"`
	Content  string `json:"content"`
	Duration int32  `json:"duration"`
}

func (q *Queries) CreateSnippet(ctx context.Context, arg CreateSnippetParams) (int32, error) {
	row := q.db.QueryRowContext(ctx, createSnippet, arg.Title, arg.Content, arg.Duration)
	var id int32
	err := row.Scan(&id)
	return id, err
}

const getSnippetNotExpired = `-- name: GetSnippetNotExpired :one
SELECT id, title, content, created_at, expires
FROM snippets
WHERE expires > CURRENT_TIMESTAMP
  AND id = $1
`

func (q *Queries) GetSnippetNotExpired(ctx context.Context, id int32) (Snippet, error) {
	row := q.db.QueryRowContext(ctx, getSnippetNotExpired, id)
	var i Snippet
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.Content,
		&i.CreatedAt,
		&i.Expires,
	)
	return i, err
}

const getTenLatestSnippets = `-- name: GetTenLatestSnippets :many
SELECT id, title, content, created_at, expires
FROM snippets
WHERE expires > CURRENT_TIMESTAMP
ORDER BY id DESC LIMIT 10
`

func (q *Queries) GetTenLatestSnippets(ctx context.Context) ([]Snippet, error) {
	rows, err := q.db.QueryContext(ctx, getTenLatestSnippets)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Snippet{}
	for rows.Next() {
		var i Snippet
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Content,
			&i.CreatedAt,
			&i.Expires,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

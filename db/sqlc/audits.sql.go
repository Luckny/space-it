// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: audits.sql

package db

import (
	"context"

	"github.com/google/uuid"
)

const createAuthenticatedRequestLog = `-- name: CreateAuthenticatedRequestLog :one
INSERT INTO request_log (method, path, user_id)
VALUES ($1, $2, $3)
RETURNING id, method, path, user_id, created_at
`

type CreateAuthenticatedRequestLogParams struct {
	Method string    `json:"method"`
	Path   string    `json:"path"`
	UserID uuid.UUID `json:"user_id"`
}

func (q *Queries) CreateAuthenticatedRequestLog(ctx context.Context, arg CreateAuthenticatedRequestLogParams) (RequestLog, error) {
	row := q.db.QueryRow(ctx, createAuthenticatedRequestLog, arg.Method, arg.Path, arg.UserID)
	var i RequestLog
	err := row.Scan(
		&i.ID,
		&i.Method,
		&i.Path,
		&i.UserID,
		&i.CreatedAt,
	)
	return i, err
}

const createResponseLog = `-- name: CreateResponseLog :one
INSERT INTO response_log (id, status)
VALUES ($1, $2)
RETURNING id, status, created_at
`

type CreateResponseLogParams struct {
	ID     uuid.UUID `json:"id"`
	Status int32     `json:"status"`
}

func (q *Queries) CreateResponseLog(ctx context.Context, arg CreateResponseLogParams) (ResponseLog, error) {
	row := q.db.QueryRow(ctx, createResponseLog, arg.ID, arg.Status)
	var i ResponseLog
	err := row.Scan(&i.ID, &i.Status, &i.CreatedAt)
	return i, err
}

const createUnauthenticatedRequestLog = `-- name: CreateUnauthenticatedRequestLog :one
INSERT INTO request_log (method, path)
VALUES ($1, $2)
RETURNING id, method, path, user_id, created_at
`

type CreateUnauthenticatedRequestLogParams struct {
	Method string `json:"method"`
	Path   string `json:"path"`
}

func (q *Queries) CreateUnauthenticatedRequestLog(ctx context.Context, arg CreateUnauthenticatedRequestLogParams) (RequestLog, error) {
	row := q.db.QueryRow(ctx, createUnauthenticatedRequestLog, arg.Method, arg.Path)
	var i RequestLog
	err := row.Scan(
		&i.ID,
		&i.Method,
		&i.Path,
		&i.UserID,
		&i.CreatedAt,
	)
	return i, err
}

// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: reset.sql

package database

import (
	"context"
)

const resetTable = `-- name: ResetTable :exec
DELETE FROM users
`

func (q *Queries) ResetTable(ctx context.Context) error {
	_, err := q.db.ExecContext(ctx, resetTable)
	return err
}
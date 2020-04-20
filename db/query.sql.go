// Code generated by sqlc. DO NOT EDIT.
// source: query.sql

package main

import (
	"context"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (
  pkey
) VALUES (
  $1
)
RETURNING id, pkey
`

func (q *Queries) CreateUser(ctx context.Context, pkey []byte) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser, pkey)
	var i User
	err := row.Scan(&i.ID, &i.Pkey)
	return i, err
}

const deleteUser = `-- name: DeleteUser :exec

DELETE FROM users
WHERE id = $1
`

func (q *Queries) DeleteUser(ctx context.Context, id int32) error {
	_, err := q.db.ExecContext(ctx, deleteUser, id)
	return err
}

const getUser = `-- name: GetUser :one
SELECT id, pkey FROM users
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetUser(ctx context.Context, id int32) (User, error) {
	row := q.db.QueryRowContext(ctx, getUser, id)
	var i User
	err := row.Scan(&i.ID, &i.Pkey)
	return i, err
}

const listUsers = `-- name: ListUsers :many
SELECT pkey FROM users
`

func (q *Queries) ListUsers(ctx context.Context) ([][]byte, error) {
	rows, err := q.db.QueryContext(ctx, listUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items [][]byte
	for rows.Next() {
		var pkey []byte
		if err := rows.Scan(&pkey); err != nil {
			return nil, err
		}
		items = append(items, pkey)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

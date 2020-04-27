// Code generated by sqlc. DO NOT EDIT.

package DB

import (
	sql "QGTodo/pkg/util/JSONunMarshal"
)

type Task struct {
	ID        int32
	FkUser    sql.JsonNullInt32
	Title     sql.NullString
	Comment   sql.NullString
	Done      sql.NullBool
	CreatedAt sql.NullTime
	UpdatedAt sql.NullTime
	DeletedAt sql.NullTime
}

type User struct {
	ID        int32
	Username  sql.NullString
	Password  []byte
	CreatedAt sql.NullTime
	UpdatedAt sql.NullTime
	DeletedAt sql.NullTime
}

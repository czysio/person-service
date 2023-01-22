// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: person.sql

package db

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const createPerson = `-- name: CreatePerson :one
INSERT INTO people (
  first_name,
  surname,
  email,
  nickname,
  created_at,
  updated_at
) VALUES (
  $1, $2, $3, $4, $5, $6
)
RETURNING id, first_name, surname, email, nickname, created_at, updated_at
`

type CreatePersonParams struct {
	FirstName string    `json:"first_name"`
	Surname   string    `json:"surname"`
	Email     string    `json:"email"`
	Nickname  string    `json:"nickname"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (q *Queries) CreatePerson(ctx context.Context, arg CreatePersonParams) (Person, error) {
	row := q.queryRow(ctx, q.createPersonStmt, createPerson,
		arg.FirstName,
		arg.Surname,
		arg.Email,
		arg.Nickname,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	var i Person
	err := row.Scan(
		&i.ID,
		&i.FirstName,
		&i.Surname,
		&i.Email,
		&i.Nickname,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deletePerson = `-- name: DeletePerson :exec
DELETE FROM people
WHERE id = $1
`

func (q *Queries) DeletePerson(ctx context.Context, id uuid.UUID) error {
	_, err := q.exec(ctx, q.deletePersonStmt, deletePerson, id)
	return err
}

const getPersonById = `-- name: GetPersonById :one
SELECT id, first_name, surname, email, nickname, created_at, updated_at FROM people
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetPersonById(ctx context.Context, id uuid.UUID) (Person, error) {
	row := q.queryRow(ctx, q.getPersonByIdStmt, getPersonById, id)
	var i Person
	err := row.Scan(
		&i.ID,
		&i.FirstName,
		&i.Surname,
		&i.Email,
		&i.Nickname,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const listPeople = `-- name: ListPeople :many
SELECT id, first_name, surname, email, nickname, created_at, updated_at FROM people
ORDER BY id
LIMIT $1
OFFSET $2
`

type ListPeopleParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListPeople(ctx context.Context, arg ListPeopleParams) ([]Person, error) {
	rows, err := q.query(ctx, q.listPeopleStmt, listPeople, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Person
	for rows.Next() {
		var i Person
		if err := rows.Scan(
			&i.ID,
			&i.FirstName,
			&i.Surname,
			&i.Email,
			&i.Nickname,
			&i.CreatedAt,
			&i.UpdatedAt,
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

const updatePerson = `-- name: UpdatePerson :one
UPDATE people
set 
first_name = coalesce($1, first_name), 
surname = coalesce($2, surname), 
email = coalesce($3, email) ,
nickname = coalesce($4, nickname), 
updated_at = coalesce($5, updated_at ) 
WHERE id = $6
RETURNING id, first_name, surname, email, nickname, created_at, updated_at
`

type UpdatePersonParams struct {
	FirstName sql.NullString `json:"first_name"`
	Surname   sql.NullString `json:"surname"`
	Email     sql.NullString `json:"email"`
	Nickname  sql.NullString `json:"nickname"`
	UpdatedAt sql.NullTime   `json:"updated_at "`
	ID        uuid.UUID      `json:"id"`
}

func (q *Queries) UpdatePerson(ctx context.Context, arg UpdatePersonParams) (Person, error) {
	row := q.queryRow(ctx, q.updatePersonStmt, updatePerson,
		arg.FirstName,
		arg.Surname,
		arg.Email,
		arg.Nickname,
		arg.UpdatedAt,
		arg.ID,
	)
	var i Person
	err := row.Scan(
		&i.ID,
		&i.FirstName,
		&i.Surname,
		&i.Email,
		&i.Nickname,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: query.sql

package sqlite

import (
	"context"
)

const addCompany = `-- name: AddCompany :exec
INSERT INTO companies (
    name,
    alias,
    site,
    priority
) VALUES (
    ?, ?, ?, ?
)
`

type AddCompanyParams struct {
	Name     string
	Alias    string
	Site     string
	Priority int64
}

func (q *Queries) AddCompany(ctx context.Context, arg AddCompanyParams) error {
	_, err := q.db.ExecContext(ctx, addCompany,
		arg.Name,
		arg.Alias,
		arg.Site,
		arg.Priority,
	)
	return err
}

const getCompanies = `-- name: GetCompanies :many
SELECT name, alias, site, priority, created_at from companies
ORDER BY site, name, priority
`

func (q *Queries) GetCompanies(ctx context.Context) ([]Company, error) {
	rows, err := q.db.QueryContext(ctx, getCompanies)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Company
	for rows.Next() {
		var i Company
		if err := rows.Scan(
			&i.Name,
			&i.Alias,
			&i.Site,
			&i.Priority,
			&i.CreatedAt,
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

const getCompaniesByMaxPriority = `-- name: GetCompaniesByMaxPriority :many
SELECT name, alias, site, priority, created_at FROM companies
WHERE priority <= ?
ORDER BY site, priority
`

func (q *Queries) GetCompaniesByMaxPriority(ctx context.Context, priority int64) ([]Company, error) {
	rows, err := q.db.QueryContext(ctx, getCompaniesByMaxPriority, priority)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Company
	for rows.Next() {
		var i Company
		if err := rows.Scan(
			&i.Name,
			&i.Alias,
			&i.Site,
			&i.Priority,
			&i.CreatedAt,
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

const getCompaniesByName = `-- name: GetCompaniesByName :many
SELECT name, alias, site, priority, created_at FROM companies
WHERE name = ?
ORDER BY site, priority
`

func (q *Queries) GetCompaniesByName(ctx context.Context, name string) ([]Company, error) {
	rows, err := q.db.QueryContext(ctx, getCompaniesByName, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Company
	for rows.Next() {
		var i Company
		if err := rows.Scan(
			&i.Name,
			&i.Alias,
			&i.Site,
			&i.Priority,
			&i.CreatedAt,
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

const getCompanyByAlias = `-- name: GetCompanyByAlias :one
SELECT name, alias, site, priority, created_at FROM companies
WHERE alias = ?
LIMIT 1
`

func (q *Queries) GetCompanyByAlias(ctx context.Context, alias string) (Company, error) {
	row := q.db.QueryRowContext(ctx, getCompanyByAlias, alias)
	var i Company
	err := row.Scan(
		&i.Name,
		&i.Alias,
		&i.Site,
		&i.Priority,
		&i.CreatedAt,
	)
	return i, err
}

const getCompanyByNameAndSite = `-- name: GetCompanyByNameAndSite :one
SELECT name, alias, site, priority, created_at FROM companies
WHERE name = ? AND site = ?
LIMIT 1
`

type GetCompanyByNameAndSiteParams struct {
	Name string
	Site string
}

func (q *Queries) GetCompanyByNameAndSite(ctx context.Context, arg GetCompanyByNameAndSiteParams) (Company, error) {
	row := q.db.QueryRowContext(ctx, getCompanyByNameAndSite, arg.Name, arg.Site)
	var i Company
	err := row.Scan(
		&i.Name,
		&i.Alias,
		&i.Site,
		&i.Priority,
		&i.CreatedAt,
	)
	return i, err
}

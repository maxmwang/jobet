-- name: GetCompanies :many
SELECT * from companies
ORDER BY site, name, priority;

-- name: GetCompaniesByName :many
SELECT * FROM companies
WHERE name = ?
ORDER BY site, priority;

-- name: GetCompaniesByMaxPriority :many
SELECT * FROM companies
WHERE priority <= ?
ORDER BY site, priority;

-- name: GetCompanyByNameAndSite :one
SELECT * FROM companies
WHERE name = ? AND site = ?
LIMIT 1;

-- name: GetCompanyByAlias :one
SELECT * FROM companies
WHERE alias = ?
LIMIT 1;

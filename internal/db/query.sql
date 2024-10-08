-- name: GetCompanies :many
SELECT * from companies
ORDER BY site, name, priority;

-- name: GetCompaniesByName :many
SELECT * FROM companies
WHERE name = $1
ORDER BY site, priority;

-- name: GetCompaniesByMaxPriority :many
SELECT * FROM companies
WHERE priority <= $1
ORDER BY site, priority;

-- name: GetCompanyByNameAndSite :one
SELECT * FROM companies
WHERE name = $1 AND site = $2
LIMIT 1;

-- name: GetCompanyByAlias :one
SELECT * FROM companies
WHERE alias = $1
LIMIT 1;

-- name: AddCompany :exec
INSERT INTO companies (
    name,
    alias,
    site,
    priority
) VALUES (
     $1, $2, $3, $4
 );

-- name: GetChannels :many
SELECT id FROM discord_channels;

-- name: AddChannel :exec
INSERT INTO discord_channels (
    id
) VALUES (
    $1
);

-- name: RemoveChannel :exec
DELETE FROM discord_channels
WHERE id = $1;

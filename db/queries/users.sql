-- name: CreateUser :one
INSERT INTO users (jira_id, username)
values ($1, $2)
returning *;
-- name: GetUserByJiraID :one 
SELECT *
FROM users
WHERE jira_id = $1;
-- name: UpdateLastLogin :exec 
update users
SET last_login = now()
WHERE id = $1;
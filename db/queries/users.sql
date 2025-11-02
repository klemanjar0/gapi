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
WHERE jira_id = $1;
-- name: GetUserWithSettings :one
SELECT u.id,
    u.jira_id,
    u.username,
    u.created_at,
    u.last_login,
    us.project_id,
    us.issue_query,
    us.content_template,
    us.ticket_item_template,
    us.mail_recipient,
    us.mail_subject,
    us.mail_author,
    us.updated_at as settings_updated_at
FROM users u
    LEFT JOIN user_settings us ON us.user_id = u.id
WHERE u.jira_id = $1;
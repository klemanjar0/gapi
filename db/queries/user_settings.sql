-- name: GetSettingsByUserID :one
SELECT *
FROM user_settings
WHERE user_id = $1;
-- name: UpsertSettings :exec 
INSERT INTO user_settings (
        user_id,
        project_id,
        issue_query,
        content_template,
        ticket_item_template,
        mail_recipient,
        mail_subject,
        mail_author
    )
values ($1, $2, $3, $4, $5, $6, $7, $8) ON conflict (user_id) do
update
SET project_id = excluded.project_id,
    issue_query = excluded.issue_query,
    content_template = excluded.content_template,
    ticket_item_template = excluded.ticket_item_template,
    mail_recipient = excluded.mail_recipient,
    mail_subject = excluded.mail_subject,
    mail_author = excluded.mail_author,
    updated_at = now();
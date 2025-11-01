create table if not exists users (
    id serial primary key,
    jira_id text not null unique,
    username text not null,
    created_at timestamptz not null default now (),
    last_login timestamptz not null default now ()
);
create table if not exists user_settings (
    user_id int not null primary key references users (id) on delete cascade,
    project_id text not null default '',
    issue_query text not null default '',
    content_template text not null default '',
    ticket_item_template text not null default '',
    mail_recipient text not null default '',
    mail_subject text not null default '',
    mail_author text not null default '',
    updated_at timestamptz not null default now ()
);
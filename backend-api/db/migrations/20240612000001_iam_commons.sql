-- migrate:up
create table if not exists principal_sessions (
    id varchar not null primary key,
    principal varchar not null default '',
    ip_address inet not null,
    user_agent varchar not null default '',
    tenant varchar not null default '',
    metadata jsonb not null default '{}',
    expires_at timestamptz not null default now() + interval '1 hour',
    created_at timestamptz not null default now()
);

create table if not exists principal_tokens (
    id varchar not null primary key,
    principal varchar not null default '',
    kind varchar not null default '',
    tenant varchar not null default '',
    expires_at timestamptz not null default now() + interval '5 minutes',
    created_at timestamptz not null default now()
);

-- Audit table
create table if not exists principal_login_events (
    id          varchar     not null,
    principal   varchar     not null default '',
    tenant      varchar     not null default '',
    ip_address  inet        not null,
    user_agent  varchar     not null default '',
    method      varchar     not null,
    success     boolean     not null default false,
    created_at  timestamptz not null default now(),
    primary key (id)
);

-- migrate:down
drop table if exists principal_sessions;
drop table if exists principal_tokens;
drop table if exists principal_login_events;
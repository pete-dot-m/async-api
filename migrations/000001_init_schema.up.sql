create table users
(
    id              uuid primary key      default gen_random_uuid(),
    email           varchar(320) not null unique,
    hashed_password varchar(96)  not null, -- bcrypt password base64 encoded
    created_at      timestamptz  not null default current_timestamp
);

create table refresh_tokens
(
    user_id      UUID REFERENCES users (id) on delete cascade,
    hashed_token varchar(500) not null unique,
    created_at   timestamptz  not null default current_timestamp,
    expires_at   timestamptz  not null default current_timestamp + interval '1 day',
    primary key (user_id, hashed_token)
);

create table reports
(
    user_id                 UUID references users (id) on delete cascade,
    id                      uuid        not null default gen_random_uuid(),
    report_type             varchar     not null,
    output_file_path        varchar,
    download_url            varchar,
    download_url_expires_at timestamptz,
    error_message           varchar,
    created_at              timestamptz not null default current_timestamp,
    started_at              timestamptz,
    primary key (user_id, id)
);

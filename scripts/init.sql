select 'create database project'
where not exists(select from pg_database where datname = 'project');
\gexec

\c project

create schema if not exists converter;

do $$
    begin
        if not exists(select 1 from pg_extension where extname = 'uuid-ossp') then
            create extension "uuid-ossp";
        end if;
    end
$$;

do $$
    begin
        if not exists(select 1 from pg_type where typname = 'file_format') then
            create type file_format as enum ('jpg', 'jpeg', 'png');
        end if;
    end
$$;

do $$
    begin
        if not exists(select 1 from pg_type where typname = 'status') then
            create type status as enum ('queued', 'processing', 'failed', 'done');
        end if;
    end
$$;

create table if not exists converter.users (
    id uuid default uuid_generate_v1() primary key,
    email varchar(50) unique not null,
    password varchar(120) not null,
    created timestamp without time zone default current_timestamp not null,
    updated timestamp without time zone default current_timestamp not null
);

create table if not exists converter.images (
    id uuid default uuid_generate_v1() primary key,
    name varchar(80) not null,
    format file_format not null,
    created timestamp without time zone default current_timestamp not null,
    updated timestamp without time zone default current_timestamp not null
);

create table if not exists converter.requests (
    id uuid default uuid_generate_v1() primary key,
    user_id uuid not null,
    source_id uuid not null,
    target_id uuid,
    source_format file_format not null,
    target_format file_format not null,
    ratio int check ( ratio > 0  and ratio < 100),
    status status not null,
    created timestamp without time zone default current_timestamp not null,
    updated timestamp without time zone default current_timestamp not null,

    foreign key (user_id) references converter.users(id),
    foreign key (source_id) references converter.images(id),
    foreign key (target_id) references converter.images(id)
);
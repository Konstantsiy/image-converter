create schema if not exists converter;

create table if not exists converter.users (
    id serial primary key,
    email varchar(50) unique not null,
    password varchar(50) not null
);

do $$
    begin
        if not exists(select 1 from pg_type where typname = 'file_format') then
            create type file_format as enum ('jpg', 'png');
        end if;
    end;
$$;

do $$
    begin
        if not exists(select 1 from pg_type where typname = 'status') then
            create type status as enum ('queued', 'processed', 'failed', 'done');
        end if;
    end;
$$;

create table if not exists converter.images (
    id serial primary key,
    name varchar(80) not null,
    format file_format not null,
    location text not null
);

create table if not exists converter.requests (
    id serial primary key,
    user_id int not null,
    image_id int not null,
    source_format file_format not null,
    target_format file_format not null,
    ratio int not null,
    status status not null,
    created timestamp default now() not null,
    updated timestamp default now() not null,

    foreign key (user_id) references converter.users(id),
    foreign key (image_id) references converter.requests(id)
);
create table if not exists product (
    id uuid primary key not null,
    title text not null,
    description text,
    price bigint not null
);
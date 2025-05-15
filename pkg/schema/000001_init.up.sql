CREATE TABLE users (
    id serial unique not null,
    email varchar(255) not null,
    password varchar(255) not null
);

CREATE TABLE files (
    id serial unique not null,
    file_size int not null,
    created_at timestamp not null default now(),
    user_id int references users(id) on delete cascade not null,
    storage_url text not null
);
create table users(
	id bigserial primary key,
	email text not null,
	password_hash text not null,
	name text not null,
	created_at timestamptz default now(),
	updated_at timestamptz default now()
)
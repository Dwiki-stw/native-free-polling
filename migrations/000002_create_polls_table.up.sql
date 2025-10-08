create table polls(
	id bigserial primary key,
	user_id bigint not null references users(id) on delete cascade,
	title text not null,
	description text,
	status text not null check (status in ('draft', 'active', 'closed', 'archived')),
	starts_at timestamptz,
	ends_at timestamptz,
	created_at timestamptz not null default now(),
	updated_at timestamptz not null default now()
)
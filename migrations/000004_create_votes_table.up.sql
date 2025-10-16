create table votes(
	id bigserial primary key,
	option_id bigint not null references poll_options(id) on delete cascade,
	device_hash text not null,
	created_at timestamptz not null default now()
)
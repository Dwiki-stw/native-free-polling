create table poll_options(
	id bigserial primary key,
	poll_id bigint not null references polls(id) on delete cascade,
	label text not null,
	position int not null,
	created_at timestamptz not null default now(),
	unique(poll_id, position)
)
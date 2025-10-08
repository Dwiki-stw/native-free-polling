create table poll_questions (
	id bigserial primary key,
	poll_id bigint not null references polls(id) on delete cascade,
	type text not null check (type in('single_choice', 'multiple_choice', 'text')),
	required boolean not null default false,
	position int not null,
	prompt text not null,
	created_at timestamptz not null default now(),
	unique(poll_id, position)
)
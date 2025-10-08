create table poll_options(
	id bigserial primary key,
	question_id bigint not null references poll_questions(id) on delete cascade,
	label text not null,
	position int not null,
	created_at timestamptz not null default now(),
	unique(question_id, position)
)
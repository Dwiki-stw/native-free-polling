create table votes(
	id bigserial primary key,
	option_id bigint references poll_options(id) on delete cascade,
	question_id bigint references poll_questions(id) on delete cascade,
	text_value text,
	device_hash text not null,
	created_at timestamptz not null default now(),
	check(
		(option_id is not null and text_value is null and question_id is null) or 
		(option_id is null and text_value is not null and question_id is not null)
	)
)
create table user_votes(
	user_id bigint not null references users(id) on delete cascade,
	vote_id bigint not null references votes(id) on delete cascade,
	primary  key (user_id, vote_id)
)
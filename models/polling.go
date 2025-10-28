package models

import "time"

type Polling struct {
	ID           int64     `db:"id"`
	UserID       int64     `db:"user_id"`
	Title        string    `db:"title"`
	Description  string    `db:"description"`
	Status       string    `db:"status"`
	StartsAt     time.Time `db:"starts_at"`
	EndsAt       time.Time `db:"ends_at"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
	CreatorName  string    `db:"creator_name"`
	CreatorEmail string    `db:"creator_email"`

	Options []PollOption
	Results []VoteResult
}

type PollingSummary struct {
	ID              int64  `db:"id"`
	Title           string `db:"title"`
	Status          string `db:"status"`
	TotalVotes      int64  `db:"total_votes"`
	UserVotedOption string `db:"voted_option"`
}

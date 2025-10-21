package models

import "time"

type PollOption struct {
	ID        int64     `db:"id"`
	PollID    int64     `db:"pol_id"`
	Label     string    `db:"label"`
	Position  int       `db:"position"`
	CreatedAt time.Time `db:"created_at"`
}

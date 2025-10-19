package models

import "time"

type Vote struct {
	ID         int64     `db:"id"`
	OptionID   int64     `db:"option_at"`
	DeviceHash string    `db:"device_hash"`
	CreatedAt  time.Time `db:"created_at"`
}

type VoteResult struct {
	OptionID    int64  `db:"option_id"`
	OptionLabel string `db:"option_label"`
	Votes       int64  `db:"votes"`
}

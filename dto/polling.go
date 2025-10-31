package dto

import "time"

type CreatePollingRequest struct {
	Title       string    `json:"title" validate:"required"`
	Description string    `json:"description" validate:"required"`
	Status      string    `json:"status" validate:"required"`
	StartsAt    time.Time `json:"starts_at" validate:"required"`
	EndsAt      time.Time `json:"ends_at" validate:"required"`
	Options     []string  `json:"options" validate:"required"`
}

type UpdatePollingRequest struct {
	ID          int64     `json:"id" validate:"required"`
	Title       string    `json:"title" validate:"required"`
	Description string    `json:"description" validate:"required"`
	Status      string    `json:"status" validate:"required"`
	StartsAt    time.Time `json:"starts_at" validate:"required"`
	EndsAt      time.Time `json:"ends_at" validate:"required"`
	Options     []Option  `json:"options" validate:"required"`
}

type VoteRequest struct {
	OptionID   int64  `json:"option_id"`
	DeviceHash string `json:"device_hash"`
}

type ChangePasswordRequest struct {
	Password string `json:"password"`
}

type PollingResponse struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	StartsAt    time.Time `json:"starts_at"`
	EndsAt      time.Time `json:"ends_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Options     []Option  `json:"polling_options"`

	Creator CreatorInfo `json:"creator"`
}

type CreatorInfo struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Option struct {
	ID       int64  `json:"id"`
	Label    string `json:"label"`
	Position int    `json:"position"`
}

type ResultPolling struct {
	PollID     int64  `json:"poll_id"`
	TotalVotes int64  `json:"total_votes"`
	Result     []Vote `json:"result"`
}

type Vote struct {
	OptionID    int64  `json:"option_id"`
	OptionLabel string `json:"optoin_label"`
	Votes       int64  `json:"votes"`
}

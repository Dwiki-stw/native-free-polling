package dto

import "time"

type ProfileResponse struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UpdateProfileRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type PollingSummaryForVoter struct {
	ID        int64  `json:"id"`
	Title     string `json:"title"`
	Status    string `json:"status"`
	UserVoted string `json:"user_voted"`
}

type PollingSummaryForCreator struct {
	ID         int64  `json:"id"`
	Title      string `json:"title"`
	Status     string `json:"status"`
	TotalVotes int64  `json:"total_votes"`
}

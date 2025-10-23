package service

import (
	"context"
	"database/sql"
	"fmt"
	"native-free-pollings/domain"
	"native-free-pollings/dto"
	"native-free-pollings/helper"
	"native-free-pollings/models"
)

type polling struct {
	DB       *sql.DB
	PollRepo domain.PollRepository
	OptRepo  domain.OptionRepository
	VoteRepo domain.VoteRepository
}

func NewPolling(db *sql.DB, pollRepo domain.PollRepository, optRepo domain.OptionRepository, voteRepo domain.VoteRepository) domain.PollService {
	return &polling{DB: db, PollRepo: pollRepo, OptRepo: optRepo, VoteRepo: voteRepo}
}

func (p *polling) CreatePolling(ctx context.Context, rq *dto.CreatePollingRequest, creator dto.CreatorInfo) (*dto.PollingResponse, error) {
	tx, err := p.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, helper.NewAppError("INTERNAL_ERROR", "internal server error", err)
	}
	defer tx.Rollback()

	//insert polling
	poll := &models.Polling{
		UserID:      creator.ID,
		Title:       rq.Title,
		Description: rq.Description,
		Status:      rq.Status,
		StartsAt:    rq.StartsAt,
		EndsAt:      rq.EndsAt,
	}

	err = p.PollRepo.Create(ctx, tx, poll)
	if err != nil {
		return nil, helper.NewAppError("DB_ERROR", "failed to save polling", err)
	}

	//insert option
	var options []dto.Option
	for i, label := range rq.Options {
		opt := &models.PollOption{
			PollID:   poll.ID,
			Label:    label,
			Position: i + 1,
		}
		err := p.OptRepo.Create(ctx, tx, opt)
		if err != nil {
			return nil, helper.NewAppError("DB_ERROR", "failed to save option", err)
		}
		dOpt := dto.Option{
			ID:       opt.ID,
			Label:    opt.Label,
			Position: opt.Position,
		}
		options = append(options, dOpt)
	}

	if err := tx.Commit(); err != nil {
		return nil, helper.NewAppError("INTERNAL_ERROR", "internal server error", err)
	}

	return &dto.PollingResponse{
		ID:          poll.ID,
		Title:       poll.Title,
		Description: poll.Description,
		Status:      poll.Status,
		StartsAt:    poll.StartsAt,
		EndsAt:      poll.EndsAt,
		CreatedAt:   poll.CreatedAt,
		UpdatedAt:   poll.UpdatedAt,
		Options:     options,
		Creator:     creator,
	}, nil

}

func (p *polling) UpdatePolling(ctx context.Context, rq *dto.UpdatePollingRequest, creator dto.CreatorInfo) (*dto.PollingResponse, error) {
	oldPoll, err := p.PollRepo.GetByID(ctx, p.DB, rq.ID)
	if err != nil {
		return nil, helper.NewAppError("INTERNAL_ERROR", "internal server error", err)
	}

	if oldPoll.UserID != creator.ID {
		return nil, helper.NewAppError("FORBIDDEN_ERROR", "only creator can update this polling", nil)
	}

	oldOptions, err := p.OptRepo.GetByPollID(ctx, p.DB, rq.ID)
	if err != nil {
		return nil, helper.NewAppError("INTERNAL_ERROR", "internal server error", err)
	}

	idOptions := make(map[int64]bool)
	for _, opt := range oldOptions {
		idOptions[opt.ID] = true
	}

	tx, err := p.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, helper.NewAppError("INTERNAL_ERROR", "internal server error", err)
	}
	defer tx.Rollback()

	//update polling
	updatedPoll := &models.Polling{
		ID:          rq.ID,
		UserID:      creator.ID,
		Title:       rq.Title,
		Description: rq.Description,
		Status:      rq.Status,
		StartsAt:    rq.StartsAt,
		EndsAt:      rq.EndsAt,
	}

	err = p.PollRepo.Update(ctx, tx, updatedPoll)
	if err != nil {
		return nil, helper.NewAppError("DB_ERROR", "failed to save polling", err)
	}

	//options
	var newOptions []dto.Option
	for _, opt := range rq.Options {
		option := &models.PollOption{
			ID:       opt.ID,
			PollID:   oldPoll.ID,
			Label:    opt.Label,
			Position: opt.Position,
		}
		//update
		if opt.ID > 0 && idOptions[opt.ID] {
			err := p.OptRepo.Update(ctx, tx, option)
			if err != nil {
				return nil, helper.NewAppError("DB_ERROR", "failed to update option", err)
			}
			newOptions = append(newOptions, dto.Option{ID: option.ID, Label: option.Label, Position: option.Position})
		} else {
			//insert
			err := p.OptRepo.Create(ctx, tx, option)
			if err != nil {
				return nil, helper.NewAppError("DB_ERROR", "failed to create new option", err)
			}
		}
		idOptions[opt.ID] = false
	}

	//delete
	for id, c := range idOptions {
		if c {
			err := p.OptRepo.Delete(ctx, tx, id)
			if err != nil {
				return nil, helper.NewAppError("DB_ERROR", "failed to delete option", err)
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, helper.NewAppError("INTERNAL_ERROR", "internal server error", err)
	}

	return &dto.PollingResponse{
		ID:          updatedPoll.ID,
		Title:       updatedPoll.Title,
		Description: updatedPoll.Description,
		Status:      updatedPoll.Status,
		StartsAt:    updatedPoll.StartsAt,
		EndsAt:      updatedPoll.EndsAt,
		CreatedAt:   updatedPoll.CreatedAt,
		UpdatedAt:   updatedPoll.UpdatedAt,
		Options:     newOptions,
		Creator:     creator,
	}, nil
}

func (p *polling) VoteOptionPolling(ctx context.Context, userID, pollID, optionID int64, deviceHash string) error {
	//checking user vote
	if userID > 0 {
		exist, err := p.VoteRepo.HasUserVoted(ctx, p.DB, pollID, userID)
		if err != nil {
			return helper.NewAppError("INTERNAL_ERROR", "internal server error", err)
		}
		if exist {
			return helper.NewAppError("ALREADY_VOTED", "you have alread voted in this polling", err)
		}
	} else {
		exist, err := p.VoteRepo.HasDeviceVoted(ctx, p.DB, deviceHash, pollID)
		if err != nil {
			return helper.NewAppError("INTERNAL_ERROR", "internal server error", err)
		}
		if exist {
			return helper.NewAppError("ALREADY_VOTED", "you have alread voted in this polling", err)
		}
	}

	tx, err := p.DB.BeginTx(ctx, nil)
	if err != nil {
		return helper.NewAppError("INTERNAL_ERROR", "internal server error", err)
	}
	defer tx.Rollback()

	poll, err := p.PollRepo.GetByID(ctx, p.DB, pollID)
	if err != nil {
		return helper.NewAppError("DB_ERROR", "failed get polling", err)
	}

	if poll.Status != "active" {
		return helper.NewAppError("BAD_REQUEST", fmt.Sprintf("cannot vote polling: %s", poll.Status), err)
	}

	vote := &models.Vote{
		OptionID:   optionID,
		DeviceHash: deviceHash,
	}
	err = p.VoteRepo.Create(ctx, tx, vote)
	if err != nil {
		return helper.NewAppError("DB_ERROR", "failed save vote", err)
	}

	if userID > 0 {
		err := p.VoteRepo.CreateUserVote(ctx, tx, userID, vote.ID)
		if err != nil {
			return helper.NewAppError("DB_ERROR", "failed insert user vote", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return helper.NewAppError("INTERNAL_ERROR", "internal server error", err)
	}

	return nil
}

func (p *polling) DeletePolling(ctx context.Context, pollID, userID int64) error {
	poll, err := p.PollRepo.GetByID(ctx, p.DB, pollID)
	if err != nil {
		return helper.NewAppError("INTERNAL_ERROR", "internal server error", err)
	}

	if poll.UserID != userID {
		return helper.NewAppError("FORBIDDEN_ERROR", "only creator can delete this polling", nil)
	}

	err = p.PollRepo.Delete(ctx, p.DB, pollID)
	if err != nil {
		return helper.NewAppError("DB_ERROR", "failed delete polling", err)
	}

	return nil
}

func (p *polling) GetDetailPolling(ctx context.Context, id int64) (*dto.PollingResponse, error) {
	poll, err := p.PollRepo.GetByID(ctx, p.DB, id)
	if err != nil {
		return nil, helper.NewAppError("DB_ERROR", "failed get detail polling", err)
	}

	opt, err := p.OptRepo.GetByPollID(ctx, p.DB, id)
	if err != nil {
		return nil, helper.NewAppError("DB_ERROR", "failed get options polling", err)
	}
	var options []dto.Option
	for _, o := range opt {
		options = append(options, dto.Option{ID: o.ID, Label: o.Label, Position: o.Position})
	}

	creator := dto.CreatorInfo{
		ID:    poll.UserID,
		Name:  poll.CreatorName,
		Email: poll.CreatorEmail,
	}

	return &dto.PollingResponse{
		ID:          poll.ID,
		Title:       poll.Title,
		Description: poll.Description,
		Status:      poll.Status,
		StartsAt:    poll.StartsAt,
		EndsAt:      poll.EndsAt,
		CreatedAt:   poll.CreatedAt,
		UpdatedAt:   poll.UpdatedAt,
		Options:     options,
		Creator:     creator,
	}, nil
}
func (p *polling) GetPollingResult(ctx context.Context, pollID int64) (*dto.ResultPolling, error) {
	vr, err := p.PollRepo.GetResultsByID(ctx, p.DB, pollID)
	if err != nil {
		return nil, helper.NewAppError("DB_ERROR", "failed get votes", err)
	}

	var totalVotes int64
	var result []dto.Vote
	for _, v := range vr {
		result = append(result, dto.Vote{OptionID: v.OptionID, OptionLabel: v.OptionLabel, Votes: v.Votes})
		totalVotes += v.Votes
	}

	return &dto.ResultPolling{
		PollID:     pollID,
		TotalVotes: totalVotes,
		Result:     result,
	}, nil
}

package handler

import (
	"encoding/json"
	"native-free-pollings/domain"
	"native-free-pollings/dto"
	"native-free-pollings/helper"
	"net/http"
	"strconv"
	"strings"
)

type Polling struct {
	Service domain.PollService
}

func NewPolling(svc domain.PollService) *Polling {
	return &Polling{Service: svc}
}

func (p *Polling) CreatePolling(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"code":    "NOT_ALLOWED",
			"message": "method not allowed",
		})
		return
	}

	auth, ok := helper.GetAuthContext(r.Context())
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"code":    "INVALID_TOKEN",
			"message": "invalid user information",
		})
		return
	}
	creator := dto.CreatorInfo{
		ID:    auth.UserID,
		Name:  auth.UserName,
		Email: auth.UserEmail,
	}

	var req dto.CreatePollingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"code":    "INVALID_REQUEST",
			"message": "invalid request payload",
		})
		return
	}

	if errs, err := helper.BindAndValidate(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"code":    "VALIDATION_ERROR",
			"message": "payload validation failed",
			"details": errs,
		})
		return
	}

	resp, err := p.Service.CreatePolling(r.Context(), &req, creator)
	if err != nil {
		err.(*helper.AppError).WriteError(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"message": "created polling successfully",
		"data":    resp,
	})
}

func (p *Polling) UpdatePolling(w http.ResponseWriter, r *http.Request) {
	auth, ok := helper.GetAuthContext(r.Context())
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"code":    "INVALID_TOKEN",
			"message": "invalid user information",
		})
		return
	}
	creator := dto.CreatorInfo{
		ID:    auth.UserID,
		Name:  auth.UserName,
		Email: auth.UserEmail,
	}

	var req dto.UpdatePollingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"code":    "INVALID_REQUEST",
			"message": "invalid request payload",
		})
		return
	}

	if errs, err := helper.BindAndValidate(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"code":    "VALIDATION_ERROR",
			"message": "payload validation failed",
			"details": errs,
		})
		return
	}

	resp, err := p.Service.UpdatePolling(r.Context(), &req, creator)
	if err != nil {
		err.(*helper.AppError).WriteError(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"message": "updated polling successfully",
		"data":    resp,
	})
}

func (p *Polling) DeletePolling(w http.ResponseWriter, r *http.Request) {
	auth, ok := helper.GetAuthContext(r.Context())
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"code":    "INVALID_TOKEN",
			"message": "invalid user information",
		})
		return
	}
	creator := dto.CreatorInfo{
		ID: auth.UserID,
	}

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"code":    "INVALID_PATH",
			"message": "invalid path id polling",
		})
		return
	}

	idStr := parts[2]
	pollID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"code":    "INVALID_ID",
			"message": "invalid id polling",
		})
		return
	}

	err = p.Service.DeletePolling(r.Context(), pollID, creator.ID)
	if err != nil {
		err.(*helper.AppError).WriteError(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"message": "deleted polling successfully",
	})
}

func (p *Polling) GetDetailPolling(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"code":    "NOT_ALLOWED",
			"message": "method not allowed",
		})
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"code":    "INVALID_PATH",
			"message": "invalid path id polling",
		})
		return
	}

	idStr := parts[2]
	pollID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"code":    "INVALID_ID",
			"message": "invalid id polling",
		})
		return
	}

	resp, err := p.Service.GetDetailPolling(r.Context(), pollID)
	if err != nil {
		err.(*helper.AppError).WriteError(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"message": "get detail polling successfully",
		"data":    resp,
	})
}

func (p *Polling) VoteOptionPolling(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"code":    "NOT_ALLOWED",
			"message": "method not allowed",
		})
		return
	}
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"code":    "INVALID_PATH",
			"message": "invalid path id polling",
		})
		return
	}

	idStr := parts[2]
	pollID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"code":    "INVALID_ID",
			"message": "invalid id polling",
		})
		return
	}

	var req struct {
		OptionID   int64  `json:"option_id"`
		DeviceHash string `json:"device_hash"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"code":    "INVALID_REQUEST",
			"message": "invalid request payload",
		})
		return
	}

	auth, ok := helper.GetAuthContext(r.Context())
	if ok {
		err := p.Service.VoteOptionPolling(r.Context(), auth.UserID, pollID, req.OptionID, req.DeviceHash)
		if err != nil {
			err.(*helper.AppError).WriteError(w)
			return
		}
	} else {
		err := p.Service.VoteOptionPolling(r.Context(), 0, pollID, req.OptionID, req.DeviceHash)
		if err != nil {
			err.(*helper.AppError).WriteError(w)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"message": "vote successfully",
	})
}

func (p *Polling) GetPollingResult(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"code":    "NOT_ALLOWED",
			"message": "method not allowed",
		})
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"code":    "INVALID_PATH",
			"message": "invalid path id polling",
		})
		return
	}

	idStr := parts[2]
	pollID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"code":    "INVALID_ID",
			"message": "invalid id polling",
		})
		return
	}

	resp, err := p.Service.GetPollingResult(r.Context(), pollID)
	if err != nil {
		err.(*helper.AppError).WriteError(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"message": "get polling result successfully",
		"data":    resp,
	})
}

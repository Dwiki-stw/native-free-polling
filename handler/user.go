package handler

import (
	"encoding/json"
	"native-free-pollings/domain"
	"native-free-pollings/dto"
	"native-free-pollings/helper"
	"native-free-pollings/models"
	"net/http"
)

type UserHandler struct {
	Service domain.UserService
}

func NewUserHandler(service domain.UserService) *UserHandler {
	return &UserHandler{Service: service}
}

// Get Profile godoc
// @Summary      get profile info user login
// @Description  Retrieves profile information of the currently authenticated user.
// @Tags         User
// @Accept       json
// @Produce      json
// @Security BearerAuth
// @Success      200      {object}  dto.ProfileResponse
// @Router       /users/me [get]
func (u *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
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
			"message": "invalid user id",
		})
		return
	}

	resp, err := u.Service.GetProfile(r.Context(), auth.UserID)
	if err != nil {
		err.(*helper.AppError).WriteError(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

// Update Profile godoc
// @Summary      update profile info user login
// @Description  Updates the profile information of the currently authenticated user.
// @Tags         User
// @Accept       json
// @Produce      json
// @Security BearerAuth
// @Param        request  body     dto.UpdateProfileRequest  true "profile update payload"
// @Success      200      {object}  dto.ProfileResponse
// @Router       /users/me [patch]
func (u *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
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
			"message": "invalid user id",
		})
		return
	}

	var req dto.UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"code":    "INVALID_REQUEST",
			"message": "invalid request payload",
		})
		return
	}

	if req.Email == "" && req.Name == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"code":    "BAD_REQUEST",
			"message": "name and email cannot both be empty",
		})
		return
	}

	user := &models.User{
		ID:    auth.UserID,
		Name:  req.Name,
		Email: req.Email,
	}

	resp, err := u.Service.UpdateProfile(r.Context(), user)
	if err != nil {
		err.(*helper.AppError).WriteError(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"message": "profile updated successfully",
		"data":    resp,
	})
}

// Update Profile godoc
// @Summary      update profile info user login
// @Description  Updates the profile information of the currently authenticated user.
// @Tags         User
// @Accept       json
// @Produce      json
// @Security BearerAuth
// @Param        request  body     dto.ChangePasswordRequest  true "change password payload"
// @Success      200      {object}  map[string]string "Success message"
// @Router       /users/me/change-password [patch]
func (u *UserHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
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
			"message": "invalid user id",
		})
		return
	}

	var req dto.ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"code":    "INVALID_REQUEST",
			"message": "invalid request payload",
		})
		return
	}

	if err := u.Service.ChangePassword(r.Context(), auth.UserID, req.Password); err != nil {
		err.(*helper.AppError).WriteError(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"message": "password updated successfully",
	})
}

// Get list polls created godoc
// @Summary      get list polls creater by user login
// @Description  Retrieves a list polls created by the logged-in user.
// @Tags         User
// @Accept       json
// @Produce      json
// @Security BearerAuth
// @Success      200      {array}  dto.PollingSummaryForCreator
// @Router       /users/me/pollings/created [get]
func (u *UserHandler) GetUserCreatedPollings(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
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
			"message": "invalid user id",
		})
		return
	}

	resp, err := u.Service.GetUserCreatedPollings(r.Context(), auth.UserID)
	if err != nil {
		err.(*helper.AppError).WriteError(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"message": "get pollings successfully",
		"data":    resp,
	})
}

// Get list polls voted godoc
// @Summary      get list polls voter by user login
// @Description  Retrieves a list polls voted on by the logged-in user.
// @Tags         User
// @Accept       json
// @Produce      json
// @Security BearerAuth
// @Success      200      {object}  dto.PollingSummaryForVoter
// @Router       /users/me/pollings/voted [get]
func (u *UserHandler) GetUserVotedPollings(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
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
			"message": "invalid user id",
		})
		return
	}

	resp, err := u.Service.GetUserVotedPollings(r.Context(), auth.UserID)
	if err != nil {
		err.(*helper.AppError).WriteError(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"message": "get pollings successfully",
		"data":    resp,
	})
}

package handler

import (
	"encoding/json"
	"native-free-pollings/domain"
	"native-free-pollings/dto"
	"native-free-pollings/helper"
	"net/http"
)

type AuthHandler struct {
	Service domain.AuthService
}

func NewAuthHandler(service domain.AuthService) *AuthHandler {
	return &AuthHandler{Service: service}
}

// Login godoc
// @Summary      Login user
// @Description  Authenticates user and returns a JWT token for secure access.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request  body      dto.LoginRequest  true  "Login credentials"
// @Success      201      {object}  dto.LoginResponse
// @Router       /login [post]
func (a *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"code":    "NOT_ALLOWED",
			"message": "method not allowed",
		})
		return
	}

	var req dto.LoginRequest
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

	resp, err := a.Service.Login(r.Context(), &req)
	if err != nil {
		err.(*helper.AppError).WriteError(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(resp)
}

// Register godoc
// @Summary      Register user
// @Description  Registers a new user with email, name, and password.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request  body      dto.RegisterRequest  true  "Register credentials"
// @Success      201      {object}  dto.RegisterResponse
// @Router       /register [post]
func (a *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"code":    "NOT_ALLOWED",
			"message": "method not allowed",
		})
		return
	}

	var req dto.RegisterRequest
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

	resp, err := a.Service.Register(r.Context(), &req)
	if err != nil {
		err.(*helper.AppError).WriteError(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(resp)
}

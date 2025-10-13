package service

import (
	"context"
	"native-free-pollings/domain"
	"native-free-pollings/dto"
	"native-free-pollings/helper"
	"native-free-pollings/models"
	"time"
)

type authService struct {
	repo   domain.AuthRepository
	jwtKey []byte
	hasher helper.PasswordHasher
}

func NewAuthService(repo domain.AuthRepository, jwtKey []byte, hasher helper.PasswordHasher) domain.AuthService {
	return &authService{repo: repo, jwtKey: jwtKey, hasher: hasher}
}

func (a *authService) Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := a.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, helper.NewAppError("LOGIN_FAILED", "invalid email or password", err)
	}

	if err := a.hasher.Compare(user.PasswordHash, req.Password); err != nil {
		return nil, helper.NewAppError("LOGIN_FAILED", "invalid email or password", err)
	}

	tokenInfo := &helper.Claims{
		UserID: user.ID,
		Name:   user.Name,
		Email:  user.Email,
		Exp:    time.Now().Add(time.Hour * 72).Unix(),
	}

	token, err := helper.CreateToken(tokenInfo, a.jwtKey)
	if err != nil {
		return nil, helper.NewAppError("TOKEN_FAILED", "failed create token", err)
	}

	return &dto.LoginResponse{
		ID:    user.ID,
		Name:  user.Name,
		Token: *token,
	}, nil
}

func (a *authService) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.RegisterResponse, error) {
	existing, _ := a.repo.GetUserByEmail(ctx, req.Email)
	if existing != nil {
		return nil, helper.NewAppError("EMAIL_EXIST", "email already registered", nil)
	}

	hashed, err := a.hasher.Hash(req.Pass)
	if err != nil {
		return nil, helper.NewAppError("HASH_FAILED", "failed hash password", err)
	}

	user := &models.User{
		Email:        req.Email,
		PasswordHash: string(hashed),
		Name:         req.Name,
	}

	err = a.repo.CreateUser(ctx, user)
	if err != nil {
		return nil, helper.NewAppError("DB_ERROR", "failed to save user", err)
	}

	return &dto.RegisterResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

package services

import (
	"context"
	"errors"
	"testing"

	"github.com/iwtcode/user-order-api/internal/models"
	"github.com/iwtcode/user-order-api/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestAuthService_Login_Success(t *testing.T) {
	repo := new(mockUserRepo)
	svc := &authService{userRepo: repo}
	ctx := context.Background()

	password := "12345678"
	hash, _ := utils.HashPassword(password)
	repo.On("GetUserByEmail", ctx, "a@b.com").Return(&models.User{Email: "a@b.com", PasswordHash: hash}, nil)

	token, err := svc.Login(ctx, "a@b.com", password)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestAuthService_Login_InvalidCredentials(t *testing.T) {
	repo := new(mockUserRepo)
	svc := &authService{userRepo: repo}
	ctx := context.Background()

	hash, _ := utils.HashPassword("otherpass")
	repo.On("GetUserByEmail", ctx, "a@b.com").Return(&models.User{Email: "a@b.com", PasswordHash: hash}, nil)

	token, err := svc.Login(ctx, "a@b.com", "wrongpass")
	assert.ErrorIs(t, err, ErrInvalidCredentials)
	assert.Empty(t, token)
}

func TestAuthService_Login_UserNotFound(t *testing.T) {
	repo := new(mockUserRepo)
	svc := &authService{userRepo: repo}
	ctx := context.Background()

	repo.On("GetUserByEmail", ctx, "notfound@b.com").Return(nil, nil)

	token, err := svc.Login(ctx, "notfound@b.com", "any")
	assert.ErrorIs(t, err, ErrInvalidCredentials)
	assert.Empty(t, token)
}

func TestAuthService_Login_RepoError(t *testing.T) {
	repo := new(mockUserRepo)
	svc := &authService{userRepo: repo}
	ctx := context.Background()

	repo.On("GetUserByEmail", ctx, "a@b.com").Return(nil, errors.New("db error"))

	token, err := svc.Login(ctx, "a@b.com", "12345678")
	assert.Error(t, err)
	assert.Empty(t, token)
}

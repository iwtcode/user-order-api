package tests

import (
	"context"
	"testing"

	"github.com/iwtcode/user-order-api/internal/models"
	"github.com/iwtcode/user-order-api/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockUserRepo struct {
	mock.Mock
}

func (m *mockUserRepo) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	args := m.Called(ctx, email)
	user, _ := args.Get(0).(*models.User)
	return user, args.Error(1)
}
func (m *mockUserRepo) CreateUser(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}
func (m *mockUserRepo) ListUsers(ctx context.Context, page, limit, minAge, maxAge int) ([]models.User, int64, error) {
	args := m.Called(ctx, page, limit, minAge, maxAge)
	users, _ := args.Get(0).([]models.User)
	total, _ := args.Get(1).(int64)
	return users, total, args.Error(2)
}
func (m *mockUserRepo) UpdateUser(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}
func (m *mockUserRepo) DeleteUser(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestUserService_CreateUser(t *testing.T) {
	repo := new(mockUserRepo)
	svc := services.NewUserService(repo)
	ctx := context.Background()

	req := &models.CreateUserRequest{Name: "Test", Email: "a@b.com", Age: 20, Password: "12345678"}

	repo.On("GetUserByEmail", ctx, req.Email).Return(nil, nil)
	repo.On("CreateUser", ctx, mock.AnythingOfType("*models.User")).Return(nil)

	user, err := svc.CreateUser(ctx, req)
	assert.NoError(t, err)
	assert.Equal(t, req.Email, user.Email)
	assert.Equal(t, req.Name, user.Name)
	assert.Equal(t, req.Age, user.Age)
}

func TestUserService_CreateUser_DuplicateEmail(t *testing.T) {
	repo := new(mockUserRepo)
	svc := services.NewUserService(repo)
	ctx := context.Background()

	req := &models.CreateUserRequest{Name: "Test", Email: "a@b.com", Age: 20, Password: "12345678"}
	repo.On("GetUserByEmail", ctx, req.Email).Return(&models.User{Email: req.Email}, nil)

	user, err := svc.CreateUser(ctx, req)
	assert.ErrorIs(t, err, services.ErrEmailExists)
	assert.Nil(t, user)
}

func TestUserService_GetUserByID(t *testing.T) {
	repo := new(mockUserRepo)
	svc := services.NewUserService(repo)
	ctx := context.Background()

	repo.On("GetUserByID", ctx, uint(1)).Return(&models.User{Email: "a@b.com"}, nil)
	user, err := svc.GetUserByID(ctx, 1)
	assert.NoError(t, err)
	assert.Equal(t, "a@b.com", user.Email)
}

func TestUserService_GetUserByID_NotFound(t *testing.T) {
	repo := new(mockUserRepo)
	svc := services.NewUserService(repo)
	ctx := context.Background()

	repo.On("GetUserByID", ctx, uint(2)).Return(nil, nil)
	user, err := svc.GetUserByID(ctx, 2)
	assert.ErrorIs(t, err, services.ErrUserNotFound)
	assert.Nil(t, user)
}

func TestUserService_ListUsers(t *testing.T) {
	repo := new(mockUserRepo)
	svc := services.NewUserService(repo)
	ctx := context.Background()

	users := []models.User{{Email: "a@b.com"}, {Email: "b@b.com"}}
	repo.On("ListUsers", ctx, 1, 10, 0, 0).Return(users, int64(2), nil)

	result, total, err := svc.ListUsers(ctx, 1, 10, 0, 0)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, result, 2)
}

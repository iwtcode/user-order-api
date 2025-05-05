package repository

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/iwtcode/user-order-api/internal/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupMockDBUser(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock db: %v", err)
	}
	dialector := postgres.New(postgres.Config{Conn: db, PreferSimpleProtocol: true})
	gdb, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open gorm db: %v", err)
	}
	cleanup := func() { db.Close() }
	return gdb, mock, cleanup
}

func TestUserRepository_CreateUser(t *testing.T) {
	db, mock, cleanup := setupMockDBUser(t)
	defer cleanup()
	repo := NewUserRepository(db)
	user := &models.User{Name: "Test", Email: "test@mail.com", Age: 30, PasswordHash: "hash"}
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "users"`)).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()
	err := repo.CreateUser(context.Background(), user)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_CreateUser_Error(t *testing.T) {
	db, mock, cleanup := setupMockDBUser(t)
	defer cleanup()
	repo := NewUserRepository(db)
	user := &models.User{Name: "Test", Email: "test@mail.com", Age: 30, PasswordHash: "hash"}
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "users"`)).WillReturnError(errors.New("db error"))
	mock.ExpectRollback()
	err := repo.CreateUser(context.Background(), user)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create user")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_GetUserByEmail(t *testing.T) {
	db, mock, cleanup := setupMockDBUser(t)
	defer cleanup()
	repo := NewUserRepository(db)
	email := "test@mail.com"
	rows := sqlmock.NewRows([]string{"id", "name", "email", "age", "password_hash"}).AddRow(1, "Test", email, 30, "hash")
	mock.ExpectQuery(`SELECT \* FROM "users" WHERE email = \$1 AND "users"\."deleted_at" IS NULL ORDER BY "users"\."id" LIMIT \$2`).WithArgs(email, 1).WillReturnRows(rows)
	user, err := repo.GetUserByEmail(context.Background(), email)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, email, user.Email)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_GetUserByEmail_NotFound(t *testing.T) {
	db, mock, cleanup := setupMockDBUser(t)
	defer cleanup()
	repo := NewUserRepository(db)
	email := "notfound@mail.com"
	mock.ExpectQuery(`SELECT \* FROM "users" WHERE email = \$1 AND "users"\."deleted_at" IS NULL ORDER BY "users"\."id" LIMIT \$2`).WithArgs(email, 1).WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email", "age", "password_hash"}))
	user, err := repo.GetUserByEmail(context.Background(), email)
	assert.NoError(t, err)
	assert.Nil(t, user)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_GetUserByID(t *testing.T) {
	db, mock, cleanup := setupMockDBUser(t)
	defer cleanup()
	repo := NewUserRepository(db)
	id := uint(1)
	rows := sqlmock.NewRows([]string{"id", "name", "email", "age", "password_hash"}).AddRow(1, "Test", "test@mail.com", 30, "hash")
	mock.ExpectQuery(`SELECT \* FROM "users" WHERE "users"\."id" = \$1 AND "users"\."deleted_at" IS NULL ORDER BY "users"\."id" LIMIT \$2`).WithArgs(id, 1).WillReturnRows(rows)
	user, err := repo.GetUserByID(context.Background(), id)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, id, user.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_GetUserByID_NotFound(t *testing.T) {
	db, mock, cleanup := setupMockDBUser(t)
	defer cleanup()
	repo := NewUserRepository(db)
	id := uint(2)
	mock.ExpectQuery(`SELECT \* FROM "users" WHERE "users"\."id" = \$1 AND "users"\."deleted_at" IS NULL ORDER BY "users"\."id" LIMIT \$2`).WithArgs(id, 1).WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email", "age", "password_hash"}))
	user, err := repo.GetUserByID(context.Background(), id)
	assert.NoError(t, err)
	assert.Nil(t, user)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_ListUsers(t *testing.T) {
	db, mock, cleanup := setupMockDBUser(t)
	defer cleanup()
	repo := NewUserRepository(db)
	rows := sqlmock.NewRows([]string{"id", "name", "email", "age", "password_hash"}).AddRow(1, "Test", "test@mail.com", 30, "hash")
	mock.ExpectQuery(`SELECT count\(\*\) FROM "users" WHERE "users"\."deleted_at" IS NULL`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	mock.ExpectQuery(`SELECT \* FROM "users" WHERE "users"\."deleted_at" IS NULL LIMIT \$1`).WithArgs(10).WillReturnRows(rows)
	users, total, err := repo.ListUsers(context.Background(), 1, 10, 0, 0)
	assert.NoError(t, err)
	assert.Len(t, users, 1)
	assert.Equal(t, int64(1), total)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_ListUsers_Error(t *testing.T) {
	db, mock, cleanup := setupMockDBUser(t)
	defer cleanup()
	repo := NewUserRepository(db)
	mock.ExpectQuery(`SELECT count\(\*\) FROM "users" WHERE "users"\."deleted_at" IS NULL`).WillReturnError(errors.New("db error"))
	users, total, err := repo.ListUsers(context.Background(), 1, 10, 0, 0)
	assert.Error(t, err)
	assert.Nil(t, users)
	assert.Equal(t, int64(0), total)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_UpdateUser(t *testing.T) {
	db, mock, cleanup := setupMockDBUser(t)
	defer cleanup()
	repo := NewUserRepository(db)
	user := &models.User{Name: "Test", Email: "test@mail.com", Age: 30, PasswordHash: "hash"}
	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "users" SET`).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	err := repo.UpdateUser(context.Background(), user)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_UpdateUser_Error(t *testing.T) {
	db, mock, cleanup := setupMockDBUser(t)
	defer cleanup()
	repo := NewUserRepository(db)
	user := &models.User{Name: "Test", Email: "test@mail.com", Age: 30, PasswordHash: "hash"}
	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "users" SET`).WillReturnError(errors.New("db error"))
	mock.ExpectRollback()
	err := repo.UpdateUser(context.Background(), user)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to update user")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_DeleteUser(t *testing.T) {
	db, mock, cleanup := setupMockDBUser(t)
	defer cleanup()
	repo := NewUserRepository(db)
	id := uint(1)
	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "users" SET "deleted_at"=\$1 WHERE "users"\."id" = \$2 AND "users"\."deleted_at" IS NULL`).WithArgs(sqlmock.AnyArg(), id).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()
	err := repo.DeleteUser(context.Background(), id)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_DeleteUser_NotFound(t *testing.T) {
	db, mock, cleanup := setupMockDBUser(t)
	defer cleanup()
	repo := NewUserRepository(db)
	id := uint(2)
	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "users" SET "deleted_at"=\$1 WHERE "users"\."id" = \$2 AND "users"\."deleted_at" IS NULL`).WithArgs(sqlmock.AnyArg(), id).WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()
	err := repo.DeleteUser(context.Background(), id)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_DeleteUser_Error(t *testing.T) {
	db, mock, cleanup := setupMockDBUser(t)
	defer cleanup()
	repo := NewUserRepository(db)
	id := uint(1)
	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "users" SET "deleted_at"=\$1 WHERE "users"\."id" = \$2 AND "users"\."deleted_at" IS NULL`).WithArgs(sqlmock.AnyArg(), id).WillReturnError(errors.New("db error"))
	mock.ExpectRollback()
	err := repo.DeleteUser(context.Background(), id)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to delete user")
	assert.NoError(t, mock.ExpectationsWereMet())
}

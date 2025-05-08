package test

import (
	"context"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/iwtcode/user-order-api/internal/models"
	"github.com/iwtcode/user-order-api/internal/repository"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupMockDBRepo(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock db: %v", err)
	}
	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open gorm db: %v", err)
	}
	cleanup := func() { db.Close() }
	return gormDB, mock, cleanup
}

func TestOrderRepository_CreateOrder(t *testing.T) {
	db, mock, cleanup := setupMockDBRepo(t)
	defer cleanup()
	repo := repository.NewOrderRepository(db)
	order := &models.Order{UserID: 1, Product: "Book", Quantity: 2, Price: 10.5}
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "orders"`)).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()
	err := repo.CreateOrder(context.Background(), order)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestOrderRepository_CreateOrder_Error(t *testing.T) {
	db, mock, cleanup := setupMockDBRepo(t)
	defer cleanup()
	repo := repository.NewOrderRepository(db)
	order := &models.Order{UserID: 1, Product: "Book", Quantity: 2, Price: 10.5}
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "orders"`)).WillReturnError(errors.New("db error"))
	mock.ExpectRollback()
	err := repo.CreateOrder(context.Background(), order)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create order")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestOrderRepository_ListOrdersByUserID(t *testing.T) {
	db, mock, cleanup := setupMockDBRepo(t)
	defer cleanup()
	repo := repository.NewOrderRepository(db)
	userID := uint(1)
	rows := sqlmock.NewRows([]string{"id", "user_id", "product", "quantity", "price", "created_at"}).
		AddRow(1, userID, "Book", 2, 10.5, time.Now())
	mock.ExpectQuery(`SELECT \* FROM "orders" WHERE user_id = \$1 ORDER BY created_at desc`).WithArgs(userID).WillReturnRows(rows)
	orders, err := repo.ListOrdersByUserID(context.Background(), userID)
	assert.NoError(t, err)
	assert.Len(t, orders, 1)
	assert.Equal(t, userID, orders[0].UserID)
	assert.Equal(t, "Book", orders[0].Product)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestOrderRepository_ListOrdersByUserID_Empty(t *testing.T) {
	db, mock, cleanup := setupMockDBRepo(t)
	defer cleanup()
	repo := repository.NewOrderRepository(db)
	userID := uint(2)
	mock.ExpectQuery(`SELECT \* FROM "orders" WHERE user_id = \$1 ORDER BY created_at desc`).WithArgs(userID).WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "product", "quantity", "price", "created_at"}))
	orders, err := repo.ListOrdersByUserID(context.Background(), userID)
	assert.NoError(t, err)
	assert.Len(t, orders, 0)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestOrderRepository_ListOrdersByUserID_Error(t *testing.T) {
	db, mock, cleanup := setupMockDBRepo(t)
	defer cleanup()
	repo := repository.NewOrderRepository(db)
	userID := uint(3)
	mock.ExpectQuery(`SELECT \* FROM "orders" WHERE user_id = \$1 ORDER BY created_at desc`).WithArgs(userID).WillReturnError(errors.New("db error"))
	orders, err := repo.ListOrdersByUserID(context.Background(), userID)
	assert.Error(t, err)
	assert.Nil(t, orders)
	assert.NoError(t, mock.ExpectationsWereMet())
}

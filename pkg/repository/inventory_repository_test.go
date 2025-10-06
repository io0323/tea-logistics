package repository

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"tea-logistics/pkg/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

/*
 * 在庫管理リポジトリのSQLモックテスト
 * データベース操作の単体テストを実装する
 */

func TestSQLInventoryRepository_CreateInventory(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewInventoryRepository(db)

	tests := []struct {
		name           string
		inventory      *models.Inventory
		mockSetup      func()
		expectedError  bool
		expectedID     int64
	}{
		{
			name: "正常な在庫作成",
			inventory: &models.Inventory{
				ProductID: 1,
				Quantity:  100,
				Location:  "東京倉庫",
				Status:    models.InventoryStatusAvailable,
			},
			mockSetup: func() {
				mock.ExpectQuery(`INSERT INTO inventory`).
					WithArgs(1, 100, "東京倉庫", models.InventoryStatusAvailable, sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
			},
			expectedError: false,
			expectedID:    1,
		},
		{
			name: "データベースエラー",
			inventory: &models.Inventory{
				ProductID: 1,
				Quantity:  100,
				Location:  "東京倉庫",
				Status:    models.InventoryStatusAvailable,
			},
			mockSetup: func() {
				mock.ExpectQuery(`INSERT INTO inventory`).
					WithArgs(1, 100, "東京倉庫", models.InventoryStatusAvailable, sqlmock.AnyArg()).
					WillReturnError(sql.ErrConnDone)
			},
			expectedError: true,
			expectedID:    0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectationsWereMet()
			tt.mockSetup()

			err := repo.CreateInventory(context.Background(), tt.inventory)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedID, tt.inventory.ID)
				assert.NotZero(t, tt.inventory.CreatedAt)
				assert.NotZero(t, tt.inventory.UpdatedAt)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestSQLInventoryRepository_GetInventory(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewInventoryRepository(db)

	tests := []struct {
		name           string
		id             int64
		mockSetup      func()
		expectedError  bool
		expectedResult *models.Inventory
	}{
		{
			name: "正常な在庫取得",
			id:   1,
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{"id", "product_id", "quantity", "location", "status", "created_at", "updated_at"}).
					AddRow(1, 1, 100, "東京倉庫", models.InventoryStatusAvailable, time.Now(), time.Now())
				mock.ExpectQuery(`SELECT id, product_id, quantity, location, status, created_at, updated_at FROM inventory WHERE id = \$1`).
					WithArgs(1).
					WillReturnRows(rows)
			},
			expectedError: false,
			expectedResult: &models.Inventory{
				ID:        1,
				ProductID: 1,
				Quantity:  100,
				Location:  "東京倉庫",
				Status:    models.InventoryStatusAvailable,
			},
		},
		{
			name: "在庫が見つからない",
			id:   999,
			mockSetup: func() {
				mock.ExpectQuery(`SELECT id, product_id, quantity, location, status, created_at, updated_at FROM inventory WHERE id = \$1`).
					WithArgs(999).
					WillReturnError(sql.ErrNoRows)
			},
			expectedError:  true,
			expectedResult: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectationsWereMet()
			tt.mockSetup()

			result, err := repo.GetInventory(context.Background(), tt.id)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectedResult.ID, result.ID)
				assert.Equal(t, tt.expectedResult.ProductID, result.ProductID)
				assert.Equal(t, tt.expectedResult.Quantity, result.Quantity)
				assert.Equal(t, tt.expectedResult.Location, result.Location)
				assert.Equal(t, tt.expectedResult.Status, result.Status)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestSQLInventoryRepository_UpdateQuantity(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewInventoryRepository(db)

	tests := []struct {
		name          string
		id            int64
		quantity      int
		mockSetup     func()
		expectedError bool
	}{
		{
			name:     "正常な在庫数更新",
			id:       1,
			quantity:  150,
			mockSetup: func() {
				mock.ExpectExec(`UPDATE inventory SET quantity = \$1, updated_at = \$2 WHERE id = \$3`).
					WithArgs(150, sqlmock.AnyArg(), 1).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			expectedError: false,
		},
		{
			name:     "在庫が見つからない",
			id:       999,
			quantity:  150,
			mockSetup: func() {
				mock.ExpectExec(`UPDATE inventory SET quantity = \$1, updated_at = \$2 WHERE id = \$3`).
					WithArgs(150, sqlmock.AnyArg(), 999).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectationsWereMet()
			tt.mockSetup()

			err := repo.UpdateQuantity(context.Background(), tt.id, tt.quantity)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestSQLInventoryRepository_GetInventoryByProduct(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewInventoryRepository(db)

	tests := []struct {
		name           string
		productID      int64
		mockSetup      func()
		expectedError  bool
		expectedResult *models.Inventory
	}{
		{
			name:      "正常な商品在庫取得",
			productID: 1,
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{"id", "product_id", "quantity", "location", "status", "created_at", "updated_at"}).
					AddRow(1, 1, 100, "東京倉庫", models.InventoryStatusAvailable, time.Now(), time.Now())
				mock.ExpectQuery(`SELECT id, product_id, quantity, location, status, created_at, updated_at FROM inventory WHERE product_id = \$1`).
					WithArgs(1).
					WillReturnRows(rows)
			},
			expectedError: false,
			expectedResult: &models.Inventory{
				ID:        1,
				ProductID: 1,
				Quantity:  100,
				Location:  "東京倉庫",
				Status:    models.InventoryStatusAvailable,
			},
		},
		{
			name:      "商品在庫が見つからない",
			productID: 999,
			mockSetup: func() {
				mock.ExpectQuery(`SELECT id, product_id, quantity, location, status, created_at, updated_at FROM inventory WHERE product_id = \$1`).
					WithArgs(999).
					WillReturnError(sql.ErrNoRows)
			},
			expectedError:  true,
			expectedResult: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectationsWereMet()
			tt.mockSetup()

			result, err := repo.GetInventoryByProduct(context.Background(), tt.productID)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectedResult.ID, result.ID)
				assert.Equal(t, tt.expectedResult.ProductID, result.ProductID)
				assert.Equal(t, tt.expectedResult.Quantity, result.Quantity)
				assert.Equal(t, tt.expectedResult.Location, result.Location)
				assert.Equal(t, tt.expectedResult.Status, result.Status)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestSQLInventoryRepository_CreateMovement(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewInventoryRepository(db)

	tests := []struct {
		name          string
		movement      *models.InventoryMovement
		mockSetup     func()
		expectedError bool
		expectedID    int64
	}{
		{
			name: "正常な在庫移動作成",
			movement: &models.InventoryMovement{
				ProductID:       1,
				FromLocation:    "東京倉庫",
				ToLocation:      "大阪倉庫",
				Quantity:        50,
				MovementType:    models.MovementTypeTransfer,
				MovementDate:    time.Now(),
				ReferenceNumber: "TRF-001",
			},
			mockSetup: func() {
				mock.ExpectQuery(`INSERT INTO inventory_movements`).
					WithArgs(1, "東京倉庫", "大阪倉庫", 50, models.MovementTypeTransfer, sqlmock.AnyArg(), "TRF-001", sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
			},
			expectedError: false,
			expectedID:    1,
		},
		{
			name: "データベースエラー",
			movement: &models.InventoryMovement{
				ProductID:       1,
				FromLocation:    "東京倉庫",
				ToLocation:      "大阪倉庫",
				Quantity:        50,
				MovementType:    models.MovementTypeTransfer,
				MovementDate:    time.Now(),
				ReferenceNumber: "TRF-001",
			},
			mockSetup: func() {
				mock.ExpectQuery(`INSERT INTO inventory_movements`).
					WithArgs(1, "東京倉庫", "大阪倉庫", 50, models.MovementTypeTransfer, sqlmock.AnyArg(), "TRF-001", sqlmock.AnyArg()).
					WillReturnError(sql.ErrConnDone)
			},
			expectedError: true,
			expectedID:    0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectationsWereMet()
			tt.mockSetup()

			err := repo.CreateMovement(context.Background(), tt.movement)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedID, tt.movement.ID)
				assert.NotZero(t, tt.movement.CreatedAt)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestSQLInventoryRepository_ListMovements(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewInventoryRepository(db)

	tests := []struct {
		name           string
		productID      int64
		mockSetup      func()
		expectedError  bool
		expectedCount  int
	}{
		{
			name:      "正常な移動履歴取得",
			productID: 1,
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{"id", "product_id", "from_location", "to_location", "quantity", "movement_type", "movement_date", "reference_number", "created_at"}).
					AddRow(1, 1, "東京倉庫", "大阪倉庫", 50, models.MovementTypeTransfer, time.Now(), "TRF-001", time.Now()).
					AddRow(2, 1, "大阪倉庫", "名古屋倉庫", 30, models.MovementTypeTransfer, time.Now(), "TRF-002", time.Now())
				mock.ExpectQuery(`SELECT id, product_id, from_location, to_location, quantity, movement_type, movement_date, reference_number, created_at FROM inventory_movements WHERE product_id = \$1 ORDER BY movement_date DESC`).
					WithArgs(1).
					WillReturnRows(rows)
			},
			expectedError: false,
			expectedCount: 2,
		},
		{
			name:      "移動履歴が存在しない",
			productID: 999,
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{"id", "product_id", "from_location", "to_location", "quantity", "movement_type", "movement_date", "reference_number", "created_at"})
				mock.ExpectQuery(`SELECT id, product_id, from_location, to_location, quantity, movement_type, movement_date, reference_number, created_at FROM inventory_movements WHERE product_id = \$1 ORDER BY movement_date DESC`).
					WithArgs(999).
					WillReturnRows(rows)
			},
			expectedError: false,
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectationsWereMet()
			tt.mockSetup()

			result, err := repo.ListMovements(context.Background(), tt.productID)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Len(t, result, tt.expectedCount)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestSQLInventoryRepository_GetInventoryByLocation(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewInventoryRepository(db)

	tests := []struct {
		name           string
		location       string
		mockSetup      func()
		expectedError  bool
		expectedCount  int
	}{
		{
			name:     "正常なロケーション別在庫取得",
			location: "東京倉庫",
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{"id", "product_id", "quantity", "location", "status", "created_at", "updated_at"}).
					AddRow(1, 1, 100, "東京倉庫", models.InventoryStatusAvailable, time.Now(), time.Now()).
					AddRow(2, 2, 50, "東京倉庫", models.InventoryStatusAvailable, time.Now(), time.Now())
				mock.ExpectQuery(`SELECT id, product_id, quantity, location, status, created_at, updated_at FROM inventory WHERE location = \$1 ORDER BY id`).
					WithArgs("東京倉庫").
					WillReturnRows(rows)
			},
			expectedError: false,
			expectedCount: 2,
		},
		{
			name:     "ロケーションに在庫が存在しない",
			location: "存在しない倉庫",
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{"id", "product_id", "quantity", "location", "status", "created_at", "updated_at"})
				mock.ExpectQuery(`SELECT id, product_id, quantity, location, status, created_at, updated_at FROM inventory WHERE location = \$1 ORDER BY id`).
					WithArgs("存在しない倉庫").
					WillReturnRows(rows)
			},
			expectedError: false,
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectationsWereMet()
			tt.mockSetup()

			result, err := repo.GetInventoryByLocation(context.Background(), tt.location)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Len(t, result, tt.expectedCount)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

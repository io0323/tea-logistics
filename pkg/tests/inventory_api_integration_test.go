package tests

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"tea-logistics/pkg/handlers"
	"tea-logistics/pkg/middleware"
	"tea-logistics/pkg/models"
	"tea-logistics/pkg/repository"
	"tea-logistics/pkg/services"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

/*
 * 在庫API統合テスト
 * HTTPリクエストからレスポンスまでの統合テストを実装する
 */

// setupTestServer テスト用のサーバーをセットアップする
func setupTestServer() (*gin.Engine, sqlmock.Sqlmock, *sql.DB) {
	db, mock, err := sqlmock.New()
	if err != nil {
		panic(fmt.Sprintf("sqlmockの初期化に失敗: %v", err))
	}

	// リポジトリの初期化
	inventoryRepo := repository.NewInventoryRepository(db)

	// サービスの初期化
	inventoryService := services.NewInventoryService(inventoryRepo)

	// ハンドラの初期化
	inventoryHandler := handlers.NewInventoryHandler(inventoryService)

	// Ginルーターの設定
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// ミドルウェアの設定
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.CorsMiddleware())

	// テスト用のルーティング設定（認証なし）
	inventory := router.Group("/api/v1/inventory")
	{
		// 在庫の更新
		inventory.PUT("", inventoryHandler.UpdateInventory)
		// 在庫移動の作成
		inventory.POST("/movements", inventoryHandler.CreateMovement)
	}

	return router, mock, db
}

func TestInventoryAPI_UpdateInventory(t *testing.T) {
	router, mock, db := setupTestServer()
	defer db.Close()

	t.Run("正常な在庫更新", func(t *testing.T) {
		// モックの設定
		rows := sqlmock.NewRows([]string{"id", "product_id", "quantity", "location", "status", "created_at", "updated_at"}).
			AddRow(1, 1, 100, "東京倉庫", models.InventoryStatusAvailable, time.Now(), time.Now())
		mock.ExpectQuery(`SELECT id, product_id, quantity, location, status, created_at, updated_at FROM inventory WHERE location = \$1 ORDER BY id`).
			WithArgs("東京倉庫").
			WillReturnRows(rows)

		mock.ExpectExec(`UPDATE inventory SET quantity = \$1, updated_at = \$2 WHERE id = \$3`).
			WithArgs(150, sqlmock.AnyArg(), 1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		// リクエストボディの作成
		requestBody := map[string]interface{}{
			"product_id": 1,
			"location":   "東京倉庫",
			"quantity":   150,
		}
		jsonBody, err := json.Marshal(requestBody)
		require.NoError(t, err)

		// リクエストの作成
		req, err := http.NewRequest("PUT", "/api/v1/inventory", bytes.NewBuffer(jsonBody))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		// レスポンスの記録
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// ステータスコードの検証
		assert.Equal(t, http.StatusOK, w.Code)

		// レスポンスボディの検証
		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "在庫を更新しました", response["message"])

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("無効なリクエスト", func(t *testing.T) {
		// リクエストボディの作成
		requestBody := map[string]interface{}{
			"product_id": "invalid",
			"location":   "",
			"quantity":   -1,
		}
		jsonBody, err := json.Marshal(requestBody)
		require.NoError(t, err)

		// リクエストの作成
		req, err := http.NewRequest("PUT", "/api/v1/inventory", bytes.NewBuffer(jsonBody))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		// レスポンスの記録
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// ステータスコードの検証
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestInventoryAPI_CreateMovement(t *testing.T) {
	router, mock, db := setupTestServer()
	defer db.Close()

	t.Run("正常な在庫移動作成", func(t *testing.T) {
		// モックの設定
		fromRows := sqlmock.NewRows([]string{"id", "product_id", "quantity", "location", "status", "created_at", "updated_at"}).
			AddRow(1, 1, 100, "東京倉庫", models.InventoryStatusAvailable, time.Now(), time.Now())
		mock.ExpectQuery(`SELECT id, product_id, quantity, location, status, created_at, updated_at FROM inventory WHERE location = \$1 ORDER BY id`).
			WithArgs("東京倉庫").
			WillReturnRows(fromRows)

		toRows := sqlmock.NewRows([]string{"id", "product_id", "quantity", "location", "status", "created_at", "updated_at"})
		mock.ExpectQuery(`SELECT id, product_id, quantity, location, status, created_at, updated_at FROM inventory WHERE location = \$1 ORDER BY id`).
			WithArgs("大阪倉庫").
			WillReturnRows(toRows)

		mock.ExpectExec(`UPDATE inventory SET quantity = \$1, updated_at = \$2 WHERE id = \$3`).
			WithArgs(50, sqlmock.AnyArg(), 1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		mock.ExpectQuery(`INSERT INTO inventory`).
			WithArgs(1, 50, "大阪倉庫", models.InventoryStatusAvailable, sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(2))

		mock.ExpectQuery(`INSERT INTO inventory_movements`).
			WithArgs(1, "東京倉庫", "大阪倉庫", 50, models.MovementTypeTransfer, sqlmock.AnyArg(), "TRF-001", sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		// リクエストボディの作成
		requestBody := map[string]interface{}{
			"product_id":       1,
			"from_location":    "東京倉庫",
			"to_location":      "大阪倉庫",
			"quantity":         50,
			"movement_type":    "transfer",
			"movement_date":    time.Now().Format(time.RFC3339),
			"reference_number": "TRF-001",
		}
		jsonBody, err := json.Marshal(requestBody)
		require.NoError(t, err)

		// リクエストの作成
		req, err := http.NewRequest("POST", "/api/v1/inventory/movements", bytes.NewBuffer(jsonBody))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		// レスポンスの記録
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// ステータスコードの検証
		assert.Equal(t, http.StatusCreated, w.Code)

		// レスポンスボディの検証
		var response models.InventoryMovement
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), response.ProductID)
		assert.Equal(t, "東京倉庫", response.FromLocation)
		assert.Equal(t, "大阪倉庫", response.ToLocation)
		assert.Equal(t, 50, response.Quantity)
		assert.Equal(t, models.MovementTypeTransfer, response.MovementType)
		assert.Equal(t, "TRF-001", response.ReferenceNumber)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("無効な移動タイプ", func(t *testing.T) {
		// リクエストボディの作成
		requestBody := map[string]interface{}{
			"product_id":       1,
			"from_location":    "東京倉庫",
			"to_location":      "大阪倉庫",
			"quantity":         50,
			"movement_type":    "invalid_type",
			"movement_date":    time.Now().Format(time.RFC3339),
			"reference_number": "TRF-001",
		}
		jsonBody, err := json.Marshal(requestBody)
		require.NoError(t, err)

		// リクエストの作成
		req, err := http.NewRequest("POST", "/api/v1/inventory/movements", bytes.NewBuffer(jsonBody))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		// レスポンスの記録
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// ステータスコードの検証
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
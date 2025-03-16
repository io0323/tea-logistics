package routes

import (
	"tea-logistics/pkg/handlers"
	"tea-logistics/pkg/middleware"
	"tea-logistics/pkg/models"

	"github.com/gin-gonic/gin"
)

/*
 * 在庫ルーティング
 * 在庫管理関連のエンドポイントを定義する
 */

// SetupInventoryRoutes 在庫ルーティングを設定する
func SetupInventoryRoutes(router *gin.Engine, handler *handlers.InventoryHandler) {
	// 認証が必要なルートグループ
	inventory := router.Group("/api/v1/inventory")
	inventory.Use(middleware.AuthMiddleware())
	{
		// 在庫情報の取得（閲覧者以上）
		inventory.GET("", middleware.RoleAuth(
			models.RoleViewer,
			models.RoleOperator,
			models.RoleManager,
			models.RoleAdmin,
		), handler.GetInventory)

		// 在庫の更新（オペレーター以上）
		inventory.PUT("", middleware.RoleAuth(
			models.RoleOperator,
			models.RoleManager,
			models.RoleAdmin,
		), handler.UpdateInventory)

		// 在庫の移動（オペレーター以上）
		inventory.POST("/transfer", middleware.RoleAuth(
			models.RoleOperator,
			models.RoleManager,
			models.RoleAdmin,
		), handler.TransferInventory)

		// 在庫の可用性チェック（閲覧者以上）
		inventory.GET("/check", middleware.RoleAuth(
			models.RoleViewer,
			models.RoleOperator,
			models.RoleManager,
			models.RoleAdmin,
		), handler.CheckAvailability)
	}
}

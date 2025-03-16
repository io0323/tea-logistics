package routes

import (
	"tea-logistics/pkg/handlers"
	"tea-logistics/pkg/middleware"
	"tea-logistics/pkg/models"

	"github.com/gin-gonic/gin"
)

/*
 * 配送ルート
 * 配送関連のエンドポイントを定義する
 */

// SetupDeliveryRoutes 配送ルートを設定する
func SetupDeliveryRoutes(router *gin.Engine, handler *handlers.DeliveryHandler) {
	// 認証が必要なルート
	deliveries := router.Group("/api/deliveries")
	deliveries.Use(middleware.AuthMiddleware())

	// 配送作成 (管理者、マネージャー)
	deliveries.POST("", middleware.RoleAuth(models.RoleAdmin, models.RoleManager), handler.CreateDelivery)

	// 配送一覧取得 (全ロール)
	deliveries.GET("", handler.ListDeliveries)

	// 配送取得 (全ロール)
	deliveries.GET("/:id", handler.GetDelivery)

	// 配送ステータス更新 (管理者、マネージャー、オペレーター)
	deliveries.PUT("/:id/status", middleware.RoleAuth(models.RoleAdmin, models.RoleManager, models.RoleOperator), handler.UpdateDeliveryStatus)

	// 配送完了 (管理者、マネージャー、オペレーター)
	deliveries.POST("/:id/complete", middleware.RoleAuth(models.RoleAdmin, models.RoleManager, models.RoleOperator), handler.CompleteDelivery)

	// 配送追跡作成 (管理者、マネージャー、オペレーター)
	deliveries.POST("/:id/tracking", middleware.RoleAuth(models.RoleAdmin, models.RoleManager, models.RoleOperator), handler.CreateDeliveryTracking)

	// 配送追跡一覧取得 (全ロール)
	deliveries.GET("/:id/tracking", handler.ListDeliveryTrackings)
}

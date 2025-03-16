package routes

import (
	"tea-logistics/pkg/handlers"
	"tea-logistics/pkg/middleware"
	"tea-logistics/pkg/models"

	"github.com/gin-gonic/gin"
)

/*
 * 配送追跡ルーティング
 * 配送追跡関連のエンドポイントを定義する
 */

// SetupTrackingRoutes 配送追跡のルーティングを設定する
func SetupTrackingRoutes(router *gin.Engine, handler *handlers.TrackingHandler) {
	// 認証が必要なルートグループ
	tracking := router.Group("/api/v1/tracking")
	tracking.Use(middleware.AuthMiddleware())
	{
		// 配送追跡の初期化（オペレーター以上）
		tracking.POST("/deliveries/:delivery_id/initialize", middleware.RoleAuth(
			models.RoleOperator,
			models.RoleManager,
			models.RoleAdmin,
		), handler.InitializeTracking)

		// 配送追跡ステータスの更新（オペレーター以上）
		tracking.PUT("/:tracking_id/status", middleware.RoleAuth(
			models.RoleOperator,
			models.RoleManager,
			models.RoleAdmin,
		), handler.UpdateTrackingStatus)

		// 配送追跡イベントの追加（オペレーター以上）
		tracking.POST("/:tracking_id/events", middleware.RoleAuth(
			models.RoleOperator,
			models.RoleManager,
			models.RoleAdmin,
		), handler.AddTrackingEvent)

		// 配送追跡情報の取得（閲覧者以上）
		tracking.GET("/:tracking_id", middleware.RoleAuth(
			models.RoleViewer,
			models.RoleOperator,
			models.RoleManager,
			models.RoleAdmin,
		), handler.GetTrackingInfo)

		// 追跡条件の設定（マネージャー以上）
		tracking.POST("/:tracking_id/conditions", middleware.RoleAuth(
			models.RoleManager,
			models.RoleAdmin,
		), handler.SetTrackingCondition)

		// 追跡条件の取得（オペレーター以上）
		tracking.GET("/:tracking_id/conditions", middleware.RoleAuth(
			models.RoleOperator,
			models.RoleManager,
			models.RoleAdmin,
		), handler.GetTrackingCondition)
	}
}

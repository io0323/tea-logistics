package routes

import (
	"tea-logistics/pkg/handlers"
	"tea-logistics/pkg/middleware"

	"github.com/gin-gonic/gin"
)

/*
 * 通知ルーティング
 * 通知関連のエンドポイントを定義する
 */

// SetupNotificationRoutes 通知ルーティングを設定する
func SetupNotificationRoutes(router *gin.Engine, handler *handlers.NotificationHandler) {
	// 認証が必要なルート
	notifications := router.Group("/api/notifications")
	notifications.Use(middleware.AuthMiddleware())
	{
		notifications.POST("", handler.CreateNotification)
		notifications.GET("", handler.ListNotifications)
		notifications.GET("/:id", handler.GetNotification)
		notifications.PUT("/:id/read", handler.MarkAsRead)
		notifications.DELETE("/:id", handler.DeleteNotification)
	}
}

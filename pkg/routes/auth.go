package routes

import (
	"tea-logistics/pkg/handlers"
	"tea-logistics/pkg/middleware"

	"github.com/gin-gonic/gin"
)

/*
 * 認証ルーティング
 * 認証関連のエンドポイントを定義する
 */

// SetupAuthRoutes 認証ルーティングを設定する
func SetupAuthRoutes(router *gin.Engine, handler *handlers.UserHandler) {
	auth := router.Group("/api/v1/auth")
	{
		// 認証不要のエンドポイント
		auth.POST("/login", handler.Login)
		auth.POST("/register", handler.Register)

		// 認証が必要なエンドポイント
		secured := auth.Group("")
		secured.Use(middleware.AuthMiddleware())
		{
			secured.GET("/profile", handler.GetProfile)
			secured.PUT("/profile", handler.UpdateProfile)
			secured.PUT("/password", handler.ChangePassword)
		}
	}
}

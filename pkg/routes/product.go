package routes

import (
	"tea-logistics/pkg/handlers"
	"tea-logistics/pkg/middleware"
	"tea-logistics/pkg/models"

	"github.com/gin-gonic/gin"
)

/*
 * 商品ルーティング
 * 商品関連のエンドポイントを定義する
 */

// SetupProductRoutes 商品ルーティングを設定する
func SetupProductRoutes(router *gin.Engine, handler *handlers.ProductHandler) {
	// 認証が必要なルートグループ
	product := router.Group("/api/v1/products")
	product.Use(middleware.AuthMiddleware())
	{
		// 商品一覧の取得（閲覧者以上）
		product.GET("", middleware.RoleAuth(
			models.RoleViewer,
			models.RoleOperator,
			models.RoleManager,
			models.RoleAdmin,
		), handler.ListProducts)

		// 商品詳細の取得（閲覧者以上）
		product.GET("/:id", middleware.RoleAuth(
			models.RoleViewer,
			models.RoleOperator,
			models.RoleManager,
			models.RoleAdmin,
		), handler.GetProduct)

		// 商品の作成（マネージャー以上）
		product.POST("", middleware.RoleAuth(
			models.RoleManager,
			models.RoleAdmin,
		), handler.CreateProduct)

		// 商品の更新（マネージャー以上）
		product.PUT("/:id", middleware.RoleAuth(
			models.RoleManager,
			models.RoleAdmin,
		), handler.UpdateProduct)

		// 商品の削除（マネージャー以上）
		product.DELETE("/:id", middleware.RoleAuth(
			models.RoleManager,
			models.RoleAdmin,
		), handler.DeleteProduct)
	}
}

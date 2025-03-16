package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"tea-logistics/pkg/config"
	"tea-logistics/pkg/database"
	"tea-logistics/pkg/handlers"
	"tea-logistics/pkg/middleware"
	"tea-logistics/pkg/repository"
	"tea-logistics/pkg/routes"
	"tea-logistics/pkg/services"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"

	"github.com/joho/godotenv"
)

func init() {
	// 環境変数の読み込み
	if err := godotenv.Load(); err != nil {
		log.Printf("警告: .envファイルが読み込めません: %v\n", err)
	}
}

func main() {
	// データベース接続の初期化
	if err := database.InitDB(); err != nil {
		log.Fatalf("データベース初期化エラー: %v", err)
	}
	defer database.CloseDB()

	// マイグレーションの実行
	migrationsDir := filepath.Join("internal", "migrations")
	if err := database.RunMigrations(migrationsDir); err != nil {
		log.Fatalf("マイグレーション実行エラー: %v", err)
	}

	// サーバーポートの取得
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("サーバーを起動しました。ポート: %s\n", port)

	// 設定の読み込み
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("設定の読み込みに失敗しました: %v", err)
	}

	// データベース接続
	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("データベース接続に失敗しました: %v", err)
	}
	defer db.Close()

	// データベース接続の確認
	if err := db.Ping(); err != nil {
		log.Fatalf("データベース接続の確認に失敗しました: %v", err)
	}

	// データベースをラップ
	dbWrapper := repository.NewSQLDatabase(db)

	// リポジトリの初期化
	userRepo := repository.NewUserRepository(db)
	productRepo := repository.NewProductRepository(db)
	inventoryRepo := repository.NewInventoryRepository(db)
	trackingRepo := repository.NewTrackingRepository(db)
	notifyRepo := repository.NewSQLNotificationRepository(dbWrapper)
	deliveryRepo := repository.NewSQLDeliveryRepository(dbWrapper)

	// サービスの初期化
	userService := services.NewUserService(userRepo)
	productService := services.NewProductService(productRepo)
	inventoryService := services.NewInventoryService(inventoryRepo)
	trackingService := services.NewTrackingService(trackingRepo)
	notifyService := services.NewNotificationService(notifyRepo, deliveryRepo)
	deliveryService := services.NewDeliveryService(deliveryRepo, inventoryRepo, notifyService)

	// ハンドラの初期化
	userHandler := handlers.NewUserHandler(userService)
	productHandler := handlers.NewProductHandler(productService)
	inventoryHandler := handlers.NewInventoryHandler(inventoryService)
	trackingHandler := handlers.NewTrackingHandler(trackingService)
	notifyHandler := handlers.NewNotificationHandler(notifyService)
	deliveryHandler := handlers.NewDeliveryHandler(deliveryService)

	// Ginルーターの設定
	router := gin.Default()

	// ミドルウェアの設定
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.CorsMiddleware())

	// ルーティングの設定
	routes.SetupAuthRoutes(router, userHandler)
	routes.SetupProductRoutes(router, productHandler)
	routes.SetupInventoryRoutes(router, inventoryHandler)
	routes.SetupTrackingRoutes(router, trackingHandler)
	routes.SetupNotificationRoutes(router, notifyHandler)
	routes.SetupDeliveryRoutes(router, deliveryHandler)

	// サーバーの設定
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	// サーバーの起動
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("サーバーの起動に失敗しました: %v", err)
		}
	}()

	// シグナルの待機
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// サーバーのシャットダウン
	log.Println("サーバーをシャットダウンします...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("サーバーのシャットダウンに失敗しました: %v", err)
	}
}

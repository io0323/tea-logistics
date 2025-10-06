package main

import (
	"context"
	"database/sql"
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
	"tea-logistics/pkg/health"
	"tea-logistics/pkg/logger"
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
	// ログ機能の初期化
	if err := logger.InitFromEnv(); err != nil {
		log.Fatalf("ログ機能の初期化に失敗しました: %v", err)
	}

	logger.Info("アプリケーションを起動しています", map[string]interface{}{
		"version": "1.0.0",
		"env":     os.Getenv("GIN_MODE"),
	})

	// データベース接続の初期化
	if err := database.InitDB(); err != nil {
		logger.Fatal("データベース初期化エラー", map[string]interface{}{
			"error": err.Error(),
		})
	}
	defer database.CloseDB()

	// マイグレーションの実行
	migrationsDir := filepath.Join("internal", "migrations")
	if err := database.RunMigrations(migrationsDir); err != nil {
		logger.Fatal("マイグレーション実行エラー", map[string]interface{}{
			"error": err.Error(),
			"dir":   migrationsDir,
		})
	}

	// サーバーポートの取得
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	logger.Info("サーバー設定", map[string]interface{}{
		"port": port,
	})

	// 設定の読み込み
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatal("設定の読み込みに失敗しました", map[string]interface{}{
			"error": err.Error(),
		})
	}

	// データベース接続
	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		logger.Fatal("データベース接続に失敗しました", map[string]interface{}{
			"error": err.Error(),
		})
	}
	defer db.Close()

	// データベース接続の確認
	if err := db.Ping(); err != nil {
		logger.Fatal("データベース接続の確認に失敗しました", map[string]interface{}{
			"error": err.Error(),
		})
	}

	logger.Info("データベース接続が確立されました")

	// ヘルスチェック機能の初期化
	health.InitGlobalHealthChecker()
	health.InitGlobalMetricsManager()

	// データベースヘルスチェックを登録
	dbHealthCheck := health.NewDatabaseHealthCheck("main_database", db)
	health.RegisterGlobalCheck(dbHealthCheck)

	// メトリクス収集を開始
	ctx := context.Background()
	go health.GetGlobalMetricsManager().StartMetricsCollection(ctx, 30*time.Second)

	logger.Info("ヘルスチェック機能を初期化しました")

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

	// ログミドルウェアの設定
	router.Use(logger.RequestIDMiddleware())
	router.Use(logger.TraceIDMiddleware())
	router.Use(logger.UserIDMiddleware())
	router.Use(logger.RequestLogger(nil))
	router.Use(logger.ResponseLogger(nil))
	router.Use(logger.ErrorLogger())

	// 既存のミドルウェアの設定
	router.Use(gin.Recovery())
	router.Use(middleware.CorsMiddleware())

	// ルーティングの設定
	routes.SetupAuthRoutes(router, userHandler)
	routes.SetupProductRoutes(router, productHandler)
	routes.SetupInventoryRoutes(router, inventoryHandler)
	routes.SetupTrackingRoutes(router, trackingHandler)
	routes.SetupNotificationRoutes(router, notifyHandler)
	routes.SetupDeliveryRoutes(router, deliveryHandler)
	
	// ヘルスチェックルートの設定
	health.SetupGlobalHealthRoutes(router)

	// サーバーの設定
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	// サーバーの起動
	go func() {
		logger.Info("サーバーを起動しています", map[string]interface{}{
			"port": port,
		})
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("サーバーの起動に失敗しました", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}()

	// シグナルの待機
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// サーバーのシャットダウン
	logger.Info("サーバーをシャットダウンします...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("サーバーのシャットダウンに失敗しました", map[string]interface{}{
			"error": err.Error(),
		})
	}

	logger.Info("サーバーが正常にシャットダウンされました")
}

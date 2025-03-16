package tests

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"tea-logistics/pkg/repository"

	_ "github.com/lib/pq"
)

/*
 * テストヘルパー関数
 * テスト用のDBセットアップと後片付けを行う
 */

var testDB *sql.DB

// setupTestDB テスト用のDBをセットアップする
func setupTestDB() repository.DB {
	if testDB != nil {
		return repository.NewSQLDatabase(testDB)
	}

	// テスト用のDB接続情報
	dbHost := os.Getenv("TEST_DB_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}
	dbPort := os.Getenv("TEST_DB_PORT")
	if dbPort == "" {
		dbPort = "5432"
	}
	dbName := os.Getenv("TEST_DB_NAME")
	if dbName == "" {
		dbName = "tea_logistics_test"
	}
	dbUser := os.Getenv("TEST_DB_USER")
	if dbUser == "" {
		dbUser = "postgres"
	}
	dbPassword := os.Getenv("TEST_DB_PASSWORD")
	if dbPassword == "" {
		dbPassword = "postgres"
	}

	// DB接続
	connStr := fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
		dbHost, dbPort, dbName, dbUser, dbPassword,
	)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("テストDB接続エラー: %v", err)
	}

	// テーブルの作成
	if err := createTestTables(db); err != nil {
		log.Fatalf("テストテーブル作成エラー: %v", err)
	}

	testDB = db
	return repository.NewSQLDatabase(db)
}

// cleanupTestDB テスト用のDBをクリーンアップする
func cleanupTestDB() {
	if testDB == nil {
		return
	}

	// テーブルのクリーンアップ
	if err := cleanupTestTables(testDB); err != nil {
		log.Printf("テストテーブルクリーンアップエラー: %v", err)
	}

	if err := testDB.Close(); err != nil {
		log.Printf("テストDB切断エラー: %v", err)
	}
	testDB = nil
}

// createTestTables テストテーブルを作成する
func createTestTables(db *sql.DB) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS notifications (
			id SERIAL PRIMARY KEY,
			type VARCHAR(50) NOT NULL,
			status VARCHAR(20) NOT NULL,
			title VARCHAR(200) NOT NULL,
			message TEXT NOT NULL,
			data JSONB,
			user_id INTEGER NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS deliveries (
			id SERIAL PRIMARY KEY,
			order_id INTEGER NOT NULL,
			status VARCHAR(50) NOT NULL,
			from_warehouse_id INTEGER NOT NULL,
			to_address TEXT NOT NULL,
			estimated_time TIMESTAMP NOT NULL,
			actual_time TIMESTAMP,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS delivery_trackings (
			id SERIAL PRIMARY KEY,
			delivery_id INTEGER NOT NULL,
			location TEXT NOT NULL,
			status VARCHAR(50) NOT NULL,
			notes TEXT,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (delivery_id) REFERENCES deliveries(id)
		)`,
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			return err
		}
	}

	return nil
}

// cleanupTestTables テストテーブルをクリーンアップする
func cleanupTestTables(db *sql.DB) error {
	queries := []string{
		"DELETE FROM delivery_trackings",
		"DELETE FROM deliveries",
		"DELETE FROM notifications",
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			return err
		}
	}

	return nil
}

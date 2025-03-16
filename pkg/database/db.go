package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

/*
 * データベース接続を管理するパッケージ
 * PostgreSQLへの接続とコネクションプールの管理を行う
 */

var db *sql.DB

// InitDB データベース接続を初期化する
func InitDB() error {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("データベース接続エラー: %v", err)
	}

	// コネクションプールの設定
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	// 接続テスト
	if err = db.Ping(); err != nil {
		return fmt.Errorf("データベース接続テストエラー: %v", err)
	}

	return nil
}

// GetDB データベース接続を取得する
func GetDB() *sql.DB {
	return db
}

// CloseDB データベース接続を閉じる
func CloseDB() error {
	if db != nil {
		return db.Close()
	}
	return nil
}

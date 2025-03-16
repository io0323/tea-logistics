package repository

import (
	"context"
	"database/sql"
)

/*
 * DBインターフェース
 * SQLデータベースの操作を抽象化する
 */

// DB データベースインターフェース
type DB interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

// SQLDatabase SQLデータベースの実装
type SQLDatabase struct {
	*sql.DB
}

// NewSQLDatabase SQLデータベースを作成する
func NewSQLDatabase(db *sql.DB) DB {
	return &SQLDatabase{db}
}

// GetContext 単一の行を取得する
func (db *SQLDatabase) GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	row := db.QueryRowContext(ctx, query, args...)
	return row.Scan(dest)
}

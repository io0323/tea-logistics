package repository

import "errors"

/*
 * リポジトリエラー
 * データベース操作に関するエラーを定義する
 */

var (
	// ErrNotFound レコードが見つからない
	ErrNotFound = errors.New("record not found")

	// ErrDuplicate レコードが重複している
	ErrDuplicate = errors.New("record already exists")

	// ErrInvalidData 無効なデータ
	ErrInvalidData = errors.New("invalid data")
)

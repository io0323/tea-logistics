package database

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strings"
)

/*
 * データベースマイグレーションを管理するパッケージ
 * マイグレーションファイルの読み込みと実行を行う
 */

// RunMigrations マイグレーションを実行する
func RunMigrations(migrationsDir string) error {
	// マイグレーションファイルの取得
	files, err := ioutil.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("マイグレーションディレクトリの読み込みエラー: %v", err)
	}

	// SQLファイルのみを抽出してソート
	var migrations []string
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".sql") {
			migrations = append(migrations, f.Name())
		}
	}
	sort.Strings(migrations)

	// マイグレーションの実行
	db := GetDB()
	for _, migration := range migrations {
		fmt.Printf("マイグレーションを実行: %s\n", migration)

		// SQLファイルの読み込み
		content, err := ioutil.ReadFile(filepath.Join(migrationsDir, migration))
		if err != nil {
			return fmt.Errorf("マイグレーションファイルの読み込みエラー: %v", err)
		}

		// トランザクション内でマイグレーションを実行
		tx, err := db.Begin()
		if err != nil {
			return fmt.Errorf("トランザクション開始エラー: %v", err)
		}

		// SQLの実行
		if _, err := tx.Exec(string(content)); err != nil {
			tx.Rollback()
			return fmt.Errorf("マイグレーション実行エラー: %v", err)
		}

		// トランザクションのコミット
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("トランザクションコミットエラー: %v", err)
		}

		fmt.Printf("マイグレーション完了: %s\n", migration)
	}

	return nil
}

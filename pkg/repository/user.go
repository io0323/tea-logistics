package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"tea-logistics/pkg/models"
)

/*
 * ユーザーリポジトリ
 * データベースとのユーザー関連の操作を管理する
 */

// UserRepository ユーザーリポジトリ
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository ユーザーリポジトリを作成する
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// CreateUser ユーザーを作成する
func (r *UserRepository) CreateUser(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (
			username, email, password_hash, name, role, status,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $7)
		RETURNING id`

	now := time.Now()
	err := r.db.QueryRowContext(ctx, query,
		user.Username,
		user.Email,
		user.Password,
		user.Name,
		user.Role,
		models.UserStatusActive,
		now,
	).Scan(&user.ID)

	if err != nil {
		return fmt.Errorf("ユーザー作成エラー: %v", err)
	}

	user.CreatedAt = now
	user.UpdatedAt = now
	return nil
}

// GetUserByEmail メールアドレスでユーザーを取得する
func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, username, email, password_hash, name, role, status,
			created_at, updated_at
		FROM users
		WHERE email = $1`

	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.Name,
		&user.Role,
		&user.Status,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("ユーザーが見つかりません")
	}
	if err != nil {
		return nil, fmt.Errorf("ユーザー取得エラー: %v", err)
	}

	return user, nil
}

// GetUserByID IDでユーザーを取得する
func (r *UserRepository) GetUserByID(ctx context.Context, id int64) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, username, email, password_hash, name, role, status,
			created_at, updated_at
		FROM users
		WHERE id = $1`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.Name,
		&user.Role,
		&user.Status,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("ユーザーが見つかりません")
	}
	if err != nil {
		return nil, fmt.Errorf("ユーザー取得エラー: %v", err)
	}

	return user, nil
}

// UpdateUser ユーザー情報を更新する
func (r *UserRepository) UpdateUser(ctx context.Context, user *models.User) error {
	query := `
		UPDATE users
		SET name = $1, updated_at = $2
		WHERE id = $3`

	result, err := r.db.ExecContext(ctx, query,
		user.Name,
		time.Now(),
		user.ID,
	)
	if err != nil {
		return fmt.Errorf("ユーザー更新エラー: %v", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("結果取得エラー: %v", err)
	}
	if rows == 0 {
		return fmt.Errorf("ユーザーが見つかりません")
	}

	return nil
}

// UpdatePassword パスワードを更新する
func (r *UserRepository) UpdatePassword(ctx context.Context, userID int64, hashedPassword string) error {
	query := `
		UPDATE users
		SET password_hash = $1, updated_at = $2
		WHERE id = $3`

	result, err := r.db.ExecContext(ctx, query,
		hashedPassword,
		time.Now(),
		userID,
	)
	if err != nil {
		return fmt.Errorf("パスワード更新エラー: %v", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("結果取得エラー: %v", err)
	}
	if rows == 0 {
		return fmt.Errorf("ユーザーが見つかりません")
	}

	return nil
}

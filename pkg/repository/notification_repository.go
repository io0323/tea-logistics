package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"tea-logistics/pkg/models"
)

/*
 * 通知リポジトリ
 * 通知関連のデータベース操作を実装する
 */

// NotificationRepository 通知リポジトリインターフェース
type NotificationRepository interface {
	// CreateNotification 通知を作成する
	CreateNotification(ctx context.Context, notification *models.Notification) error

	// GetNotification 通知を取得する
	GetNotification(ctx context.Context, id int64) (*models.Notification, error)

	// ListNotifications ユーザーの通知一覧を取得する
	ListNotifications(ctx context.Context, userID int64) ([]*models.Notification, error)

	// UpdateNotificationStatus 通知ステータスを更新する
	UpdateNotificationStatus(ctx context.Context, id int64, status models.NotificationStatus) error

	// DeleteNotification 通知を削除する
	DeleteNotification(ctx context.Context, id int64) error
}

// SQLNotificationRepository SQL通知リポジトリ
type SQLNotificationRepository struct {
	db DB
}

// NewSQLNotificationRepository SQL通知リポジトリを作成する
func NewSQLNotificationRepository(db DB) NotificationRepository {
	return &SQLNotificationRepository{db: db}
}

// CreateNotification 通知を作成する
func (r *SQLNotificationRepository) CreateNotification(ctx context.Context, notification *models.Notification) error {
	query := `
		INSERT INTO notifications (
			type, status, title, message,
			data, user_id, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5::jsonb, $6, $7, $7)
		RETURNING id`

	now := time.Now()

	// データをJSONに変換
	var jsonData []byte
	var err error
	if notification.Data != nil {
		jsonData, err = json.Marshal(notification.Data)
		if err != nil {
			return fmt.Errorf("データのJSON変換エラー: %v", err)
		}
	} else {
		jsonData = []byte("{}")
	}

	err = r.db.QueryRowContext(ctx, query,
		notification.Type,
		notification.Status,
		notification.Title,
		notification.Message,
		jsonData,
		notification.UserID,
		now,
	).Scan(&notification.ID)

	if err != nil {
		return fmt.Errorf("通知作成エラー: %v", err)
	}

	notification.CreatedAt = now
	notification.UpdatedAt = now
	return nil
}

// GetNotification 通知を取得する
func (r *SQLNotificationRepository) GetNotification(ctx context.Context, id int64) (*models.Notification, error) {
	notification := &models.Notification{}
	query := `
		SELECT id, type, status, title, message,
			data, user_id, created_at, updated_at
		FROM notifications
		WHERE id = $1`

	var jsonData []byte
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&notification.ID,
		&notification.Type,
		&notification.Status,
		&notification.Title,
		&notification.Message,
		&jsonData,
		&notification.UserID,
		&notification.CreatedAt,
		&notification.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("通知取得エラー: %v", err)
	}

	// JSONデータをmapに変換
	if len(jsonData) > 0 {
		if err := json.Unmarshal(jsonData, &notification.Data); err != nil {
			return nil, fmt.Errorf("データのJSON変換エラー: %v", err)
		}
	} else {
		notification.Data = make(map[string]interface{})
	}

	return notification, nil
}

// ListNotifications 通知一覧を取得する
func (r *SQLNotificationRepository) ListNotifications(ctx context.Context, userID int64) ([]*models.Notification, error) {
	query := `
		SELECT id, type, status, title, message,
			data, user_id, created_at, updated_at
		FROM notifications
		WHERE user_id = $1
		ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("通知一覧取得エラー: %v", err)
	}
	defer rows.Close()

	var notifications []*models.Notification
	for rows.Next() {
		notification := &models.Notification{}
		var jsonData []byte
		err := rows.Scan(
			&notification.ID,
			&notification.Type,
			&notification.Status,
			&notification.Title,
			&notification.Message,
			&jsonData,
			&notification.UserID,
			&notification.CreatedAt,
			&notification.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("通知データ読み取りエラー: %v", err)
		}

		// JSONデータをmapに変換
		if len(jsonData) > 0 {
			if err := json.Unmarshal(jsonData, &notification.Data); err != nil {
				return nil, fmt.Errorf("データのJSON変換エラー: %v", err)
			}
		} else {
			notification.Data = make(map[string]interface{})
		}

		notifications = append(notifications, notification)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("通知一覧読み取りエラー: %v", err)
	}

	return notifications, nil
}

// UpdateNotificationStatus 通知ステータスを更新する
func (r *SQLNotificationRepository) UpdateNotificationStatus(ctx context.Context, id int64, status models.NotificationStatus) error {
	query := `
		UPDATE notifications
		SET status = $1, updated_at = $2
		WHERE id = $3`

	result, err := r.db.ExecContext(ctx, query,
		status,
		time.Now(),
		id,
	)
	if err != nil {
		return fmt.Errorf("通知ステータス更新エラー: %v", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("結果取得エラー: %v", err)
	}
	if rows == 0 {
		return ErrNotFound
	}

	return nil
}

// DeleteNotification 通知を削除する
func (r *SQLNotificationRepository) DeleteNotification(ctx context.Context, id int64) error {
	query := `DELETE FROM notifications WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("通知削除エラー: %v", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("結果取得エラー: %v", err)
	}
	if rows == 0 {
		return ErrNotFound
	}

	return nil
}

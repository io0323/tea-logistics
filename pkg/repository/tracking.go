package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"tea-logistics/pkg/models"
)

/*
 * 配送追跡リポジトリ
 * データベースとの配送追跡関連の操作を管理する
 */

// TrackingRepository 配送追跡リポジトリ
type TrackingRepository struct {
	db *sql.DB
}

// NewTrackingRepository 配送追跡リポジトリを作成する
func NewTrackingRepository(db *sql.DB) *TrackingRepository {
	return &TrackingRepository{db: db}
}

// CreateTracking 配送追跡情報を作成する
func (r *TrackingRepository) CreateTracking(ctx context.Context, tracking *models.TrackingInfo) error {
	query := `
		INSERT INTO tracking_info (
			id, delivery_id, status, current_location,
			estimated_time, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $6)
		RETURNING created_at, updated_at`

	now := time.Now()
	err := r.db.QueryRowContext(ctx, query,
		tracking.ID,
		tracking.DeliveryID,
		tracking.Status,
		tracking.CurrentLocation,
		tracking.EstimatedTime,
		now,
	).Scan(&tracking.CreatedAt, &tracking.UpdatedAt)

	if err != nil {
		return fmt.Errorf("配送追跡情報作成エラー: %v", err)
	}

	return nil
}

// GetTracking 配送追跡情報を取得する
func (r *TrackingRepository) GetTracking(ctx context.Context, trackingID string) (*models.TrackingInfo, error) {
	tracking := &models.TrackingInfo{}
	query := `
		SELECT id, delivery_id, status, current_location,
			estimated_time, created_at, updated_at
		FROM tracking_info
		WHERE id = $1`

	err := r.db.QueryRowContext(ctx, query, trackingID).Scan(
		&tracking.ID,
		&tracking.DeliveryID,
		&tracking.Status,
		&tracking.CurrentLocation,
		&tracking.EstimatedTime,
		&tracking.CreatedAt,
		&tracking.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("配送追跡情報が見つかりません")
	}
	if err != nil {
		return nil, fmt.Errorf("配送追跡情報取得エラー: %v", err)
	}

	// イベント履歴の取得
	events, err := r.GetTrackingEvents(ctx, trackingID)
	if err != nil {
		return nil, err
	}
	tracking.Events = events

	return tracking, nil
}

// UpdateTrackingStatus 配送追跡ステータスを更新する
func (r *TrackingRepository) UpdateTrackingStatus(ctx context.Context, trackingID string, status models.TrackingStatus, location string) error {
	query := `
		UPDATE tracking_info
		SET status = $1, current_location = $2, updated_at = $3
		WHERE id = $4`

	result, err := r.db.ExecContext(ctx, query, status, location, time.Now(), trackingID)
	if err != nil {
		return fmt.Errorf("配送追跡ステータス更新エラー: %v", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("結果取得エラー: %v", err)
	}
	if rows == 0 {
		return fmt.Errorf("配送追跡情報が見つかりません")
	}

	return nil
}

// AddTrackingEvent 配送追跡イベントを追加する
func (r *TrackingRepository) AddTrackingEvent(ctx context.Context, event *models.TrackingEvent) error {
	query := `
		INSERT INTO tracking_events (
			tracking_id, status, location, description,
			latitude, longitude, temperature, humidity,
			created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id`

	now := time.Now()
	err := r.db.QueryRowContext(ctx, query,
		event.TrackingID,
		event.Status,
		event.Location,
		event.Description,
		event.Latitude,
		event.Longitude,
		event.Temperature,
		event.Humidity,
		now,
	).Scan(&event.ID)

	if err != nil {
		return fmt.Errorf("配送追跡イベント追加エラー: %v", err)
	}

	event.CreatedAt = now
	return nil
}

// GetTrackingEvents 配送追跡イベントを取得する
func (r *TrackingRepository) GetTrackingEvents(ctx context.Context, trackingID string) ([]*models.TrackingEvent, error) {
	query := `
		SELECT id, tracking_id, status, location,
			description, latitude, longitude,
			temperature, humidity, created_at
		FROM tracking_events
		WHERE tracking_id = $1
		ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, trackingID)
	if err != nil {
		return nil, fmt.Errorf("配送追跡イベント取得エラー: %v", err)
	}
	defer rows.Close()

	var events []*models.TrackingEvent
	for rows.Next() {
		event := &models.TrackingEvent{}
		err := rows.Scan(
			&event.ID,
			&event.TrackingID,
			&event.Status,
			&event.Location,
			&event.Description,
			&event.Latitude,
			&event.Longitude,
			&event.Temperature,
			&event.Humidity,
			&event.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("配送追跡イベントスキャンエラー: %v", err)
		}
		events = append(events, event)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("配送追跡イベント行処理エラー: %v", err)
	}

	return events, nil
}

// CreateTrackingCondition 追跡条件を作成する
func (r *TrackingRepository) CreateTrackingCondition(ctx context.Context, condition *models.TrackingCondition) error {
	query := `
		INSERT INTO tracking_conditions (
			tracking_id, min_temperature, max_temperature,
			min_humidity, max_humidity, check_interval,
			notify_email, notify_phone, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $9)
		RETURNING id`

	now := time.Now()
	err := r.db.QueryRowContext(ctx, query,
		condition.TrackingID,
		condition.MinTemperature,
		condition.MaxTemperature,
		condition.MinHumidity,
		condition.MaxHumidity,
		condition.CheckInterval,
		condition.NotifyEmail,
		condition.NotifyPhone,
		now,
	).Scan(&condition.ID)

	if err != nil {
		return fmt.Errorf("追跡条件作成エラー: %v", err)
	}

	condition.CreatedAt = now
	condition.UpdatedAt = now
	return nil
}

// GetTrackingCondition 追跡条件を取得する
func (r *TrackingRepository) GetTrackingCondition(ctx context.Context, trackingID string) (*models.TrackingCondition, error) {
	condition := &models.TrackingCondition{}
	query := `
		SELECT id, tracking_id, min_temperature, max_temperature,
			min_humidity, max_humidity, check_interval,
			notify_email, notify_phone, created_at, updated_at
		FROM tracking_conditions
		WHERE tracking_id = $1`

	err := r.db.QueryRowContext(ctx, query, trackingID).Scan(
		&condition.ID,
		&condition.TrackingID,
		&condition.MinTemperature,
		&condition.MaxTemperature,
		&condition.MinHumidity,
		&condition.MaxHumidity,
		&condition.CheckInterval,
		&condition.NotifyEmail,
		&condition.NotifyPhone,
		&condition.CreatedAt,
		&condition.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("追跡条件が見つかりません")
	}
	if err != nil {
		return nil, fmt.Errorf("追跡条件取得エラー: %v", err)
	}

	return condition, nil
}

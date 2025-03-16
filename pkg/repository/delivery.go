package repository

import (
	"context"
	"fmt"
	"time"

	"tea-logistics/pkg/models"
)

/*
 * 配送管理リポジトリ
 * データベースとの配送関連の操作を管理する
 */

// DeliveryRepository 配送リポジトリインターフェース
type DeliveryRepository interface {
	CreateDelivery(ctx context.Context, delivery *models.Delivery) error
	GetDelivery(ctx context.Context, id int64) (*models.Delivery, error)
	ListDeliveries(ctx context.Context) ([]*models.Delivery, error)
	UpdateDelivery(ctx context.Context, delivery *models.Delivery) error
	CreateDeliveryItem(ctx context.Context, item *models.DeliveryItem) error
	ListDeliveryItems(ctx context.Context, deliveryID int64) ([]*models.DeliveryItem, error)
	CreateDeliveryTracking(ctx context.Context, tracking *models.DeliveryTracking) error
	ListDeliveryTrackings(ctx context.Context, deliveryID int64) ([]*models.DeliveryTracking, error)
}

// SQLDeliveryRepository SQL配送管理リポジトリ
type SQLDeliveryRepository struct {
	db DB
}

// NewSQLDeliveryRepository SQL配送管理リポジトリを作成する
func NewSQLDeliveryRepository(db DB) DeliveryRepository {
	return &SQLDeliveryRepository{db: db}
}

// CreateDelivery 配送を作成する
func (r *SQLDeliveryRepository) CreateDelivery(ctx context.Context, delivery *models.Delivery) error {
	query := `
		INSERT INTO deliveries (
			order_id, status, from_warehouse_id,
			to_address, estimated_time, actual_time,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $7)
		RETURNING id`

	now := time.Now()
	err := r.db.QueryRowContext(ctx, query,
		delivery.OrderID,
		delivery.Status,
		delivery.FromWarehouseID,
		delivery.ToAddress,
		delivery.EstimatedTime,
		delivery.ActualTime,
		now,
	).Scan(&delivery.ID)

	if err != nil {
		return fmt.Errorf("配送作成エラー: %v", err)
	}

	delivery.CreatedAt = now
	delivery.UpdatedAt = now
	return nil
}

// GetDelivery 配送を取得する
func (r *SQLDeliveryRepository) GetDelivery(ctx context.Context, id int64) (*models.Delivery, error) {
	delivery := &models.Delivery{}
	query := `
		SELECT id, order_id, status, from_warehouse_id,
			to_address, estimated_time, actual_time,
			created_at, updated_at
		FROM deliveries
		WHERE id = $1`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&delivery.ID,
		&delivery.OrderID,
		&delivery.Status,
		&delivery.FromWarehouseID,
		&delivery.ToAddress,
		&delivery.EstimatedTime,
		&delivery.ActualTime,
		&delivery.CreatedAt,
		&delivery.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("配送取得エラー: %v", err)
	}

	return delivery, nil
}

// ListDeliveries 配送一覧を取得する
func (r *SQLDeliveryRepository) ListDeliveries(ctx context.Context) ([]*models.Delivery, error) {
	query := `
		SELECT id, order_id, status, from_warehouse_id,
			to_address, estimated_time, actual_time,
			created_at, updated_at
		FROM deliveries
		ORDER BY id`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("配送一覧取得エラー: %v", err)
	}
	defer rows.Close()

	var deliveries []*models.Delivery
	for rows.Next() {
		delivery := &models.Delivery{}
		err := rows.Scan(
			&delivery.ID,
			&delivery.OrderID,
			&delivery.Status,
			&delivery.FromWarehouseID,
			&delivery.ToAddress,
			&delivery.EstimatedTime,
			&delivery.ActualTime,
			&delivery.CreatedAt,
			&delivery.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("配送データ読み取りエラー: %v", err)
		}
		deliveries = append(deliveries, delivery)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("配送一覧読み取りエラー: %v", err)
	}

	return deliveries, nil
}

// UpdateDelivery 配送を更新する
func (r *SQLDeliveryRepository) UpdateDelivery(ctx context.Context, delivery *models.Delivery) error {
	query := `
		UPDATE deliveries
		SET order_id = $1, status = $2, from_warehouse_id = $3,
			to_address = $4, estimated_time = $5, actual_time = $6,
			updated_at = $7
		WHERE id = $8`

	result, err := r.db.ExecContext(ctx, query,
		delivery.OrderID,
		delivery.Status,
		delivery.FromWarehouseID,
		delivery.ToAddress,
		delivery.EstimatedTime,
		delivery.ActualTime,
		time.Now(),
		delivery.ID,
	)
	if err != nil {
		return fmt.Errorf("配送更新エラー: %v", err)
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

// CreateDeliveryItem 配送商品を作成する
func (r *SQLDeliveryRepository) CreateDeliveryItem(ctx context.Context, item *models.DeliveryItem) error {
	query := `
		INSERT INTO delivery_items (
			delivery_id, product_id, quantity
		) VALUES ($1, $2, $3)
		RETURNING id`

	err := r.db.QueryRowContext(ctx, query,
		item.DeliveryID,
		item.ProductID,
		item.Quantity,
	).Scan(&item.ID)

	if err != nil {
		return fmt.Errorf("配送商品作成エラー: %v", err)
	}

	return nil
}

// ListDeliveryItems 配送商品一覧を取得する
func (r *SQLDeliveryRepository) ListDeliveryItems(ctx context.Context, deliveryID int64) ([]*models.DeliveryItem, error) {
	query := `
		SELECT id, delivery_id, product_id, quantity
		FROM delivery_items
		WHERE delivery_id = $1
		ORDER BY id`

	rows, err := r.db.QueryContext(ctx, query, deliveryID)
	if err != nil {
		return nil, fmt.Errorf("配送商品一覧取得エラー: %v", err)
	}
	defer rows.Close()

	var items []*models.DeliveryItem
	for rows.Next() {
		item := &models.DeliveryItem{}
		err := rows.Scan(
			&item.ID,
			&item.DeliveryID,
			&item.ProductID,
			&item.Quantity,
		)
		if err != nil {
			return nil, fmt.Errorf("配送商品データ読み取りエラー: %v", err)
		}
		items = append(items, item)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("配送商品一覧読み取りエラー: %v", err)
	}

	return items, nil
}

// CreateDeliveryTracking 配送追跡を作成する
func (r *SQLDeliveryRepository) CreateDeliveryTracking(ctx context.Context, tracking *models.DeliveryTracking) error {
	query := `
		INSERT INTO delivery_trackings (
			delivery_id, location, status,
			notes, created_at
		) VALUES ($1, $2, $3, $4, $5)
		RETURNING id`

	now := time.Now()
	err := r.db.QueryRowContext(ctx, query,
		tracking.DeliveryID,
		tracking.Location,
		tracking.Status,
		tracking.Notes,
		now,
	).Scan(&tracking.ID)

	if err != nil {
		return fmt.Errorf("配送追跡作成エラー: %v", err)
	}

	tracking.CreatedAt = now
	return nil
}

// ListDeliveryTrackings 配送追跡一覧を取得する
func (r *SQLDeliveryRepository) ListDeliveryTrackings(ctx context.Context, deliveryID int64) ([]*models.DeliveryTracking, error) {
	query := `
		SELECT id, delivery_id, location, status, notes, created_at
		FROM delivery_trackings
		WHERE delivery_id = $1
		ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, deliveryID)
	if err != nil {
		return nil, fmt.Errorf("配送追跡一覧取得エラー: %v", err)
	}
	defer rows.Close()

	var trackings []*models.DeliveryTracking
	for rows.Next() {
		tracking := &models.DeliveryTracking{}
		err := rows.Scan(
			&tracking.ID,
			&tracking.DeliveryID,
			&tracking.Location,
			&tracking.Status,
			&tracking.Notes,
			&tracking.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("配送追跡データ読み取りエラー: %v", err)
		}
		trackings = append(trackings, tracking)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("配送追跡一覧読み取りエラー: %v", err)
	}

	return trackings, nil
}

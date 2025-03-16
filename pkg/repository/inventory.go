package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"tea-logistics/pkg/models"
)

/*
 * 在庫管理リポジトリ
 * データベースとの在庫関連の操作を管理する
 */

// SQLInventoryRepository SQL在庫管理リポジトリ
type SQLInventoryRepository struct {
	db *sql.DB
}

// NewInventoryRepository 在庫管理リポジトリを作成する
func NewInventoryRepository(db *sql.DB) InventoryRepository {
	return &SQLInventoryRepository{db: db}
}

// CreateInventory 在庫を作成する
func (r *SQLInventoryRepository) CreateInventory(ctx context.Context, inventory *models.Inventory) error {
	query := `
		INSERT INTO inventory (
			product_id, quantity, location, status,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $5)
		RETURNING id`

	now := time.Now()
	err := r.db.QueryRowContext(ctx, query,
		inventory.ProductID,
		inventory.Quantity,
		inventory.Location,
		inventory.Status,
		now,
	).Scan(&inventory.ID)

	if err != nil {
		return fmt.Errorf("在庫作成エラー: %v", err)
	}

	inventory.CreatedAt = now
	inventory.UpdatedAt = now
	return nil
}

// GetInventory 在庫を取得する
func (r *SQLInventoryRepository) GetInventory(ctx context.Context, id int64) (*models.Inventory, error) {
	inventory := &models.Inventory{}
	query := `
		SELECT id, product_id, quantity, location, status,
			created_at, updated_at
		FROM inventory
		WHERE id = $1`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&inventory.ID,
		&inventory.ProductID,
		&inventory.Quantity,
		&inventory.Location,
		&inventory.Status,
		&inventory.CreatedAt,
		&inventory.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("在庫が見つかりません")
	}
	if err != nil {
		return nil, fmt.Errorf("在庫取得エラー: %v", err)
	}

	return inventory, nil
}

// ListInventories 在庫一覧を取得する
func (r *SQLInventoryRepository) ListInventories(ctx context.Context) ([]*models.Inventory, error) {
	query := `
		SELECT id, product_id, quantity, location, status,
			created_at, updated_at
		FROM inventory
		ORDER BY id`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("在庫一覧取得エラー: %v", err)
	}
	defer rows.Close()

	var inventories []*models.Inventory
	for rows.Next() {
		inventory := &models.Inventory{}
		err := rows.Scan(
			&inventory.ID,
			&inventory.ProductID,
			&inventory.Quantity,
			&inventory.Location,
			&inventory.Status,
			&inventory.CreatedAt,
			&inventory.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("在庫データ読み取りエラー: %v", err)
		}
		inventories = append(inventories, inventory)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("在庫一覧読み取りエラー: %v", err)
	}

	return inventories, nil
}

// UpdateInventory 在庫を更新する
func (r *SQLInventoryRepository) UpdateInventory(ctx context.Context, inventory *models.Inventory) error {
	query := `
		UPDATE inventory
		SET product_id = $1, quantity = $2, location = $3,
			status = $4, updated_at = $5
		WHERE id = $6`

	result, err := r.db.ExecContext(ctx, query,
		inventory.ProductID,
		inventory.Quantity,
		inventory.Location,
		inventory.Status,
		time.Now(),
		inventory.ID,
	)
	if err != nil {
		return fmt.Errorf("在庫更新エラー: %v", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("結果取得エラー: %v", err)
	}
	if rows == 0 {
		return fmt.Errorf("在庫が見つかりません")
	}

	return nil
}

// DeleteInventory 在庫を削除する
func (r *SQLInventoryRepository) DeleteInventory(ctx context.Context, id int64) error {
	query := `DELETE FROM inventory WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("在庫削除エラー: %v", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("結果取得エラー: %v", err)
	}
	if rows == 0 {
		return fmt.Errorf("在庫が見つかりません")
	}

	return nil
}

// GetInventoryByProduct 商品IDから在庫を取得する
func (r *SQLInventoryRepository) GetInventoryByProduct(ctx context.Context, productID int64) (*models.Inventory, error) {
	inventory := &models.Inventory{}
	query := `
		SELECT id, product_id, quantity, location, status,
			created_at, updated_at
		FROM inventory
		WHERE product_id = $1`

	err := r.db.QueryRowContext(ctx, query, productID).Scan(
		&inventory.ID,
		&inventory.ProductID,
		&inventory.Quantity,
		&inventory.Location,
		&inventory.Status,
		&inventory.CreatedAt,
		&inventory.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("在庫が見つかりません")
	}
	if err != nil {
		return nil, fmt.Errorf("在庫取得エラー: %v", err)
	}

	return inventory, nil
}

// UpdateQuantity 在庫数を更新する
func (r *SQLInventoryRepository) UpdateQuantity(ctx context.Context, id int64, quantity int) error {
	query := `
		UPDATE inventory
		SET quantity = $1, updated_at = $2
		WHERE id = $3`

	result, err := r.db.ExecContext(ctx, query,
		quantity,
		time.Now(),
		id,
	)
	if err != nil {
		return fmt.Errorf("在庫数更新エラー: %v", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("結果取得エラー: %v", err)
	}
	if rows == 0 {
		return fmt.Errorf("在庫が見つかりません")
	}

	return nil
}

// GetInventoryByLocation 場所から在庫を取得する
func (r *SQLInventoryRepository) GetInventoryByLocation(ctx context.Context, location string) ([]*models.Inventory, error) {
	query := `
		SELECT id, product_id, quantity, location, status,
			created_at, updated_at
		FROM inventory
		WHERE location = $1
		ORDER BY id`

	rows, err := r.db.QueryContext(ctx, query, location)
	if err != nil {
		return nil, fmt.Errorf("在庫一覧取得エラー: %v", err)
	}
	defer rows.Close()

	var inventories []*models.Inventory
	for rows.Next() {
		inventory := &models.Inventory{}
		err := rows.Scan(
			&inventory.ID,
			&inventory.ProductID,
			&inventory.Quantity,
			&inventory.Location,
			&inventory.Status,
			&inventory.CreatedAt,
			&inventory.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("在庫データ読み取りエラー: %v", err)
		}
		inventories = append(inventories, inventory)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("在庫一覧読み取りエラー: %v", err)
	}

	return inventories, nil
}

// CreateMovement 在庫移動を作成する
func (r *SQLInventoryRepository) CreateMovement(ctx context.Context, movement *models.InventoryMovement) error {
	query := `
		INSERT INTO inventory_movements (
			product_id, from_location, to_location,
			quantity, movement_type, movement_date,
			reference_number, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id`

	now := time.Now()
	err := r.db.QueryRowContext(ctx, query,
		movement.ProductID,
		movement.FromLocation,
		movement.ToLocation,
		movement.Quantity,
		movement.MovementType,
		movement.MovementDate,
		movement.ReferenceNumber,
		now,
	).Scan(&movement.ID)

	if err != nil {
		return fmt.Errorf("在庫移動作成エラー: %v", err)
	}

	movement.CreatedAt = now
	return nil
}

// ListMovements 在庫移動履歴を取得する
func (r *SQLInventoryRepository) ListMovements(ctx context.Context, productID int64) ([]*models.InventoryMovement, error) {
	query := `
		SELECT id, product_id, from_location, to_location,
			quantity, movement_type, movement_date,
			reference_number, created_at
		FROM inventory_movements
		WHERE product_id = $1
		ORDER BY movement_date DESC`

	rows, err := r.db.QueryContext(ctx, query, productID)
	if err != nil {
		return nil, fmt.Errorf("在庫移動履歴取得エラー: %v", err)
	}
	defer rows.Close()

	var movements []*models.InventoryMovement
	for rows.Next() {
		movement := &models.InventoryMovement{}
		err := rows.Scan(
			&movement.ID,
			&movement.ProductID,
			&movement.FromLocation,
			&movement.ToLocation,
			&movement.Quantity,
			&movement.MovementType,
			&movement.MovementDate,
			&movement.ReferenceNumber,
			&movement.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("在庫移動データ読み取りエラー: %v", err)
		}
		movements = append(movements, movement)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("在庫移動履歴読み取りエラー: %v", err)
	}

	return movements, nil
}

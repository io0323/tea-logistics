package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"tea-logistics/pkg/models"
)

/*
 * 商品リポジトリ
 * データベースとの商品関連の操作を管理する
 */

// SQLProductRepository SQL商品リポジトリ
type SQLProductRepository struct {
	db *sql.DB
}

// NewProductRepository 商品リポジトリを作成する
func NewProductRepository(db *sql.DB) ProductRepository {
	return &SQLProductRepository{db: db}
}

// CreateProduct 商品を作成する
func (r *SQLProductRepository) CreateProduct(ctx context.Context, product *models.Product) error {
	query := `
		INSERT INTO products (
			name, description, price, status,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $5)
		RETURNING id`

	now := time.Now()
	err := r.db.QueryRowContext(ctx, query,
		product.Name,
		product.Description,
		product.Price,
		product.Status,
		now,
	).Scan(&product.ID)

	if err != nil {
		return fmt.Errorf("商品作成エラー: %v", err)
	}

	product.CreatedAt = now
	product.UpdatedAt = now
	return nil
}

// GetProduct 商品を取得する
func (r *SQLProductRepository) GetProduct(ctx context.Context, id int64) (*models.Product, error) {
	product := &models.Product{}
	query := `
		SELECT id, name, description, price, status,
			created_at, updated_at
		FROM products
		WHERE id = $1`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.Status,
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("商品が見つかりません")
	}
	if err != nil {
		return nil, fmt.Errorf("商品取得エラー: %v", err)
	}

	return product, nil
}

// ListProducts 商品一覧を取得する
func (r *SQLProductRepository) ListProducts(ctx context.Context) ([]*models.Product, error) {
	query := `
		SELECT id, name, description, price, status,
			created_at, updated_at
		FROM products
		ORDER BY id`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("商品一覧取得エラー: %v", err)
	}
	defer rows.Close()

	var products []*models.Product
	for rows.Next() {
		product := &models.Product{}
		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.Status,
			&product.CreatedAt,
			&product.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("商品データ読み取りエラー: %v", err)
		}
		products = append(products, product)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("商品一覧読み取りエラー: %v", err)
	}

	return products, nil
}

// UpdateProduct 商品を更新する
func (r *SQLProductRepository) UpdateProduct(ctx context.Context, product *models.Product) error {
	query := `
		UPDATE products
		SET name = $1, description = $2, price = $3,
			status = $4, updated_at = $5
		WHERE id = $6`

	result, err := r.db.ExecContext(ctx, query,
		product.Name,
		product.Description,
		product.Price,
		product.Status,
		time.Now(),
		product.ID,
	)
	if err != nil {
		return fmt.Errorf("商品更新エラー: %v", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("結果取得エラー: %v", err)
	}
	if rows == 0 {
		return fmt.Errorf("商品が見つかりません")
	}

	return nil
}

// DeleteProduct 商品を削除する
func (r *SQLProductRepository) DeleteProduct(ctx context.Context, id int64) error {
	query := `DELETE FROM products WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("商品削除エラー: %v", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("結果取得エラー: %v", err)
	}
	if rows == 0 {
		return fmt.Errorf("商品が見つかりません")
	}

	return nil
}

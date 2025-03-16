package repository

import (
	"context"

	"tea-logistics/pkg/models"
)

/*
 * 商品リポジトリインターフェース
 * データベースとの商品関連の操作を定義する
 */

// ProductRepository 商品リポジトリインターフェース
type ProductRepository interface {
	CreateProduct(ctx context.Context, product *models.Product) error
	GetProduct(ctx context.Context, id int64) (*models.Product, error)
	ListProducts(ctx context.Context) ([]*models.Product, error)
	UpdateProduct(ctx context.Context, product *models.Product) error
	DeleteProduct(ctx context.Context, id int64) error
}

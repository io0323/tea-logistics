package repository

import (
	"context"

	"tea-logistics/pkg/models"
)

/*
 * 在庫リポジトリインターフェース
 * データベースとの在庫関連の操作を定義する
 */

// InventoryRepository 在庫リポジトリインターフェース
type InventoryRepository interface {
	// 基本的なCRUD操作
	CreateInventory(ctx context.Context, inventory *models.Inventory) error
	GetInventory(ctx context.Context, id int64) (*models.Inventory, error)
	ListInventories(ctx context.Context) ([]*models.Inventory, error)
	UpdateInventory(ctx context.Context, inventory *models.Inventory) error
	DeleteInventory(ctx context.Context, id int64) error

	// 在庫特有の操作
	GetInventoryByProduct(ctx context.Context, productID int64) (*models.Inventory, error)
	UpdateQuantity(ctx context.Context, id int64, quantity int) error
	GetInventoryByLocation(ctx context.Context, location string) ([]*models.Inventory, error)
	CreateMovement(ctx context.Context, movement *models.InventoryMovement) error
	ListMovements(ctx context.Context, productID int64) ([]*models.InventoryMovement, error)
}

package services

import (
	"context"
	"fmt"

	"tea-logistics/pkg/models"
	"tea-logistics/pkg/repository"
)

/*
 * 在庫管理サービス
 * 在庫関連のビジネスロジックを実装する
 */

// InventoryService 在庫管理サービス
type InventoryService struct {
	repo repository.InventoryRepository
}

// NewInventoryService 在庫管理サービスを作成する
func NewInventoryService(repo repository.InventoryRepository) *InventoryService {
	return &InventoryService{repo: repo}
}

// CreateInventory 在庫を作成する
func (s *InventoryService) CreateInventory(ctx context.Context, req *models.CreateInventoryRequest) (*models.Inventory, error) {
	inventory := &models.Inventory{
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
		Location:  req.Location,
		Status:    req.Status,
	}

	if err := s.repo.CreateInventory(ctx, inventory); err != nil {
		return nil, fmt.Errorf("在庫作成エラー: %v", err)
	}

	return inventory, nil
}

// GetInventory 在庫を取得する
func (s *InventoryService) GetInventory(ctx context.Context, id int64) (*models.Inventory, error) {
	inventory, err := s.repo.GetInventory(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("在庫取得エラー: %v", err)
	}

	return inventory, nil
}

// ListInventories 在庫一覧を取得する
func (s *InventoryService) ListInventories(ctx context.Context) ([]*models.Inventory, error) {
	inventories, err := s.repo.ListInventories(ctx)
	if err != nil {
		return nil, fmt.Errorf("在庫一覧取得エラー: %v", err)
	}

	return inventories, nil
}

// UpdateInventory 在庫を更新する
func (s *InventoryService) UpdateInventory(ctx context.Context, id int64, req *models.UpdateInventoryRequest) (*models.Inventory, error) {
	inventory, err := s.repo.GetInventory(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("在庫取得エラー: %v", err)
	}

	inventory.Quantity = req.Quantity
	inventory.Location = req.Location
	inventory.Status = req.Status

	if err := s.repo.UpdateInventory(ctx, inventory); err != nil {
		return nil, fmt.Errorf("在庫更新エラー: %v", err)
	}

	return inventory, nil
}

// DeleteInventory 在庫を削除する
func (s *InventoryService) DeleteInventory(ctx context.Context, id int64) error {
	if err := s.repo.DeleteInventory(ctx, id); err != nil {
		return fmt.Errorf("在庫削除エラー: %v", err)
	}

	return nil
}

// GetInventoryByProduct 商品IDから在庫を取得する
func (s *InventoryService) GetInventoryByProduct(ctx context.Context, productID int64) (*models.Inventory, error) {
	inventory, err := s.repo.GetInventoryByProduct(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("在庫取得エラー: %v", err)
	}

	return inventory, nil
}

// UpdateQuantity 在庫数を更新する
func (s *InventoryService) UpdateQuantity(ctx context.Context, id int64, quantity int) error {
	if err := s.repo.UpdateQuantity(ctx, id, quantity); err != nil {
		return fmt.Errorf("在庫数更新エラー: %v", err)
	}

	return nil
}

// GetInventoryByLocation 場所から在庫を取得する
func (s *InventoryService) GetInventoryByLocation(ctx context.Context, location string) ([]*models.Inventory, error) {
	inventories, err := s.repo.GetInventoryByLocation(ctx, location)
	if err != nil {
		return nil, fmt.Errorf("在庫取得エラー: %v", err)
	}

	return inventories, nil
}

// CreateMovement 在庫移動を作成する
func (s *InventoryService) CreateMovement(ctx context.Context, req *models.CreateMovementRequest) (*models.InventoryMovement, error) {
	// 移動元の在庫を確認
	fromInventory, err := s.repo.GetInventoryByProduct(ctx, req.ProductID)
	if err != nil {
		return nil, fmt.Errorf("移動元在庫取得エラー: %v", err)
	}

	// 在庫数のチェック
	if fromInventory.Quantity < req.Quantity {
		return nil, fmt.Errorf("在庫が不足しています")
	}

	// 移動元の在庫を減らす
	if err := s.repo.UpdateQuantity(ctx, fromInventory.ID, fromInventory.Quantity-req.Quantity); err != nil {
		return nil, fmt.Errorf("移動元在庫更新エラー: %v", err)
	}

	// 移動先の在庫を増やす
	toInventory, err := s.repo.GetInventoryByProduct(ctx, req.ProductID)
	if err == nil {
		// 既存の在庫がある場合は更新
		if err := s.repo.UpdateQuantity(ctx, toInventory.ID, toInventory.Quantity+req.Quantity); err != nil {
			return nil, fmt.Errorf("移動先在庫更新エラー: %v", err)
		}
	} else {
		// 新規在庫を作成
		newInventory := &models.Inventory{
			ProductID: req.ProductID,
			Quantity:  req.Quantity,
			Location:  req.ToLocation,
			Status:    "available",
		}
		if err := s.repo.CreateInventory(ctx, newInventory); err != nil {
			return nil, fmt.Errorf("移動先在庫作成エラー: %v", err)
		}
	}

	// 在庫移動を記録
	movement := &models.InventoryMovement{
		ProductID:       req.ProductID,
		FromLocation:    req.FromLocation,
		ToLocation:      req.ToLocation,
		Quantity:        req.Quantity,
		MovementType:    req.MovementType,
		MovementDate:    req.MovementDate,
		ReferenceNumber: req.ReferenceNumber,
	}

	if err := s.repo.CreateMovement(ctx, movement); err != nil {
		return nil, fmt.Errorf("在庫移動作成エラー: %v", err)
	}

	return movement, nil
}

// ListMovements 在庫移動履歴を取得する
func (s *InventoryService) ListMovements(ctx context.Context, productID int64) ([]*models.InventoryMovement, error) {
	movements, err := s.repo.ListMovements(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("在庫移動履歴取得エラー: %v", err)
	}

	return movements, nil
}

// GetProductInventory 商品の在庫を取得する
func (s *InventoryService) GetProductInventory(ctx context.Context, productID int64, location string) (*models.Inventory, error) {
	inventories, err := s.repo.GetInventoryByLocation(ctx, location)
	if err != nil {
		return nil, fmt.Errorf("在庫取得エラー: %v", err)
	}

	for _, inv := range inventories {
		if inv.ProductID == productID {
			return inv, nil
		}
	}

	return nil, fmt.Errorf("在庫が見つかりません")
}

// UpdateInventoryQuantity 在庫数を更新する
func (s *InventoryService) UpdateInventoryQuantity(ctx context.Context, productID int64, location string, quantity int) error {
	inventory, err := s.GetProductInventory(ctx, productID, location)
	if err != nil {
		return err
	}

	if err := s.repo.UpdateQuantity(ctx, inventory.ID, quantity); err != nil {
		return fmt.Errorf("在庫数更新エラー: %v", err)
	}

	return nil
}

// TransferInventory 在庫を移動する
func (s *InventoryService) TransferInventory(ctx context.Context, productID int64, fromLocation, toLocation string, quantity int) error {
	// 移動元の在庫を確認
	fromInventory, err := s.GetProductInventory(ctx, productID, fromLocation)
	if err != nil {
		return err
	}

	if fromInventory.Quantity < quantity {
		return fmt.Errorf("在庫が不足しています")
	}

	// 移動先の在庫を確認
	toInventory, err := s.GetProductInventory(ctx, productID, toLocation)
	if err != nil {
		// 移動先に在庫がない場合は新規作成
		toInventory = &models.Inventory{
			ProductID: productID,
			Quantity:  0,
			Location:  toLocation,
			Status:    models.InventoryStatusAvailable,
		}
		if err := s.repo.CreateInventory(ctx, toInventory); err != nil {
			return fmt.Errorf("移動先在庫作成エラー: %v", err)
		}
	}

	// 移動元の在庫を減らす
	if err := s.repo.UpdateQuantity(ctx, fromInventory.ID, fromInventory.Quantity-quantity); err != nil {
		return fmt.Errorf("移動元在庫更新エラー: %v", err)
	}

	// 移動先の在庫を増やす
	if err := s.repo.UpdateQuantity(ctx, toInventory.ID, toInventory.Quantity+quantity); err != nil {
		// 移動先の更新に失敗した場合、移動元を元に戻す
		if restoreErr := s.repo.UpdateQuantity(ctx, fromInventory.ID, fromInventory.Quantity); restoreErr != nil {
			return fmt.Errorf("移動先在庫更新エラー: %v, 移動元在庫復元エラー: %v", err, restoreErr)
		}
		return fmt.Errorf("移動先在庫更新エラー: %v", err)
	}

	return nil
}

// CheckAvailability 在庫の利用可能性をチェックする
func (s *InventoryService) CheckAvailability(ctx context.Context, productID int64, location string, quantity int) (bool, error) {
	inventory, err := s.GetProductInventory(ctx, productID, location)
	if err != nil {
		return false, err
	}

	return inventory.Quantity >= quantity, nil
}

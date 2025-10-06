package services

import (
	"context"
	"fmt"

	"tea-logistics/pkg/logger"
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
	logger.Info("在庫移動処理開始", map[string]interface{}{
		"product_id":       req.ProductID,
		"from_location":    req.FromLocation,
		"to_location":      req.ToLocation,
		"quantity":         req.Quantity,
		"movement_type":    req.MovementType,
		"reference_number": req.ReferenceNumber,
	})

	// 移動元の在庫をロケーション単位で取得
	fromInventory, err := s.GetProductInventory(ctx, req.ProductID, req.FromLocation)
	if err != nil {
		logger.Error("移動元在庫取得エラー", map[string]interface{}{
			"product_id":    req.ProductID,
			"from_location": req.FromLocation,
			"error":         err.Error(),
		})
		return nil, fmt.Errorf("移動元在庫取得エラー: %v", err)
	}

	logger.Info("移動元在庫情報", map[string]interface{}{
		"inventory_id": fromInventory.ID,
		"quantity":     fromInventory.Quantity,
		"location":     fromInventory.Location,
	})

	// 在庫数のチェック
	if fromInventory.Quantity < req.Quantity {
		logger.Warn("在庫不足", map[string]interface{}{
			"product_id":    req.ProductID,
			"from_location": req.FromLocation,
			"required":      req.Quantity,
			"available":     fromInventory.Quantity,
		})
		return nil, fmt.Errorf("在庫が不足しています")
	}

	// 移動元の在庫を減らす
	newFromQuantity := fromInventory.Quantity - req.Quantity
	if err := s.repo.UpdateQuantity(ctx, fromInventory.ID, newFromQuantity); err != nil {
		logger.Error("移動元在庫更新エラー", map[string]interface{}{
			"inventory_id": fromInventory.ID,
			"new_quantity": newFromQuantity,
			"error":        err.Error(),
		})
		return nil, fmt.Errorf("移動元在庫更新エラー: %v", err)
	}

	logger.Info("移動元在庫更新完了", map[string]interface{}{
		"inventory_id": fromInventory.ID,
		"old_quantity": fromInventory.Quantity,
		"new_quantity": newFromQuantity,
	})

	// 移動先の在庫をロケーション単位で取得
	toInventory, err := s.GetProductInventory(ctx, req.ProductID, req.ToLocation)
	if err == nil {
		// 既存の移動先在庫がある場合は加算
		newToQuantity := toInventory.Quantity + req.Quantity
		if err := s.repo.UpdateQuantity(ctx, toInventory.ID, newToQuantity); err != nil {
			logger.Error("移動先在庫更新エラー", map[string]interface{}{
				"inventory_id": toInventory.ID,
				"new_quantity": newToQuantity,
				"error":        err.Error(),
			})
			return nil, fmt.Errorf("移動先在庫更新エラー: %v", err)
		}
		logger.Info("移動先在庫更新完了", map[string]interface{}{
			"inventory_id": toInventory.ID,
			"old_quantity": toInventory.Quantity,
			"new_quantity": newToQuantity,
		})
	} else {
		// 移動先に在庫がない場合は新規作成
		newInventory := &models.Inventory{
			ProductID: req.ProductID,
			Quantity:  req.Quantity,
			Location:  req.ToLocation,
			Status:    models.InventoryStatusAvailable,
		}
		if err := s.repo.CreateInventory(ctx, newInventory); err != nil {
			logger.Error("移動先在庫作成エラー", map[string]interface{}{
				"product_id":  req.ProductID,
				"to_location": req.ToLocation,
				"quantity":    req.Quantity,
				"error":       err.Error(),
			})
			return nil, fmt.Errorf("移動先在庫作成エラー: %v", err)
		}
		logger.Info("移動先在庫新規作成完了", map[string]interface{}{
			"inventory_id": newInventory.ID,
			"product_id":   newInventory.ProductID,
			"location":     newInventory.Location,
			"quantity":     newInventory.Quantity,
		})
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
		logger.Error("在庫移動作成エラー", map[string]interface{}{
			"product_id":       req.ProductID,
			"from_location":    req.FromLocation,
			"to_location":      req.ToLocation,
			"quantity":         req.Quantity,
			"movement_type":    req.MovementType,
			"reference_number": req.ReferenceNumber,
			"error":            err.Error(),
		})
		return nil, fmt.Errorf("在庫移動作成エラー: %v", err)
	}

	logger.Info("在庫移動処理完了", map[string]interface{}{
		"movement_id":      movement.ID,
		"product_id":       movement.ProductID,
		"from_location":    movement.FromLocation,
		"to_location":      movement.ToLocation,
		"quantity":         movement.Quantity,
		"movement_type":    movement.MovementType,
		"reference_number": movement.ReferenceNumber,
	})

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

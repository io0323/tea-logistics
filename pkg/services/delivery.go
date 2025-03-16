package services

import (
	"context"
	"fmt"
	"time"

	"tea-logistics/pkg/models"
	"tea-logistics/pkg/repository"
)

/*
 * 配送サービス
 * 配送関連のビジネスロジックを実装する
 */

// DeliveryService 配送サービス
type DeliveryService struct {
	repo          repository.DeliveryRepository
	inventoryRepo repository.InventoryRepository
	notifyService NotificationService
}

// NewDeliveryService 配送サービスを作成する
func NewDeliveryService(
	repo repository.DeliveryRepository,
	inventoryRepo repository.InventoryRepository,
	notifyService NotificationService,
) *DeliveryService {
	return &DeliveryService{
		repo:          repo,
		inventoryRepo: inventoryRepo,
		notifyService: notifyService,
	}
}

// CreateDelivery 配送を作成する
func (s *DeliveryService) CreateDelivery(ctx context.Context, req *models.CreateDeliveryRequest) (*models.Delivery, error) {
	// 在庫の確認
	inventory, err := s.inventoryRepo.GetInventory(ctx, req.ProductID)
	if err != nil {
		return nil, fmt.Errorf("在庫確認エラー: %v", err)
	}

	if inventory.Quantity < req.Quantity {
		return nil, fmt.Errorf("在庫が不足しています")
	}

	// 在庫の更新
	inventory.Quantity -= req.Quantity
	if err := s.inventoryRepo.UpdateInventory(ctx, inventory); err != nil {
		return nil, fmt.Errorf("在庫更新エラー: %v", err)
	}

	// 配送の作成
	delivery := &models.Delivery{
		OrderID:         req.OrderID,
		Status:          "pending",
		FromWarehouseID: req.FromWarehouseID,
		ToAddress:       req.ToAddress,
		EstimatedTime:   req.EstimatedTime,
	}

	if err := s.repo.CreateDelivery(ctx, delivery); err != nil {
		// 配送作成に失敗した場合、在庫を戻す
		inventory.Quantity += req.Quantity
		if updateErr := s.inventoryRepo.UpdateInventory(ctx, inventory); updateErr != nil {
			return nil, fmt.Errorf("配送作成エラー: %v, 在庫復元エラー: %v", err, updateErr)
		}
		return nil, fmt.Errorf("配送作成エラー: %v", err)
	}

	// 配送商品の作成
	item := &models.DeliveryItem{
		DeliveryID: delivery.ID,
		ProductID:  req.ProductID,
		Quantity:   req.Quantity,
	}

	if err := s.repo.CreateDeliveryItem(ctx, item); err != nil {
		return nil, fmt.Errorf("配送商品作成エラー: %v", err)
	}

	// 配送作成の通知
	if s.notifyService != nil {
		if err := s.notifyService.NotifyDeliveryStatusChange(ctx, delivery); err != nil {
			// 通知エラーはログに記録するだけで、配送作成自体は成功とする
			fmt.Printf("通知エラー: %v\n", err)
		}
	}

	return delivery, nil
}

// GetDelivery 配送を取得する
func (s *DeliveryService) GetDelivery(ctx context.Context, id int64) (*models.Delivery, error) {
	delivery, err := s.repo.GetDelivery(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("配送取得エラー: %v", err)
	}

	return delivery, nil
}

// ListDeliveries 配送一覧を取得する
func (s *DeliveryService) ListDeliveries(ctx context.Context) ([]*models.Delivery, error) {
	deliveries, err := s.repo.ListDeliveries(ctx)
	if err != nil {
		return nil, fmt.Errorf("配送一覧取得エラー: %v", err)
	}

	return deliveries, nil
}

// UpdateDeliveryStatus 配送ステータスを更新する
func (s *DeliveryService) UpdateDeliveryStatus(ctx context.Context, id int64, status string) error {
	delivery, err := s.repo.GetDelivery(ctx, id)
	if err != nil {
		return fmt.Errorf("配送取得エラー: %v", err)
	}

	delivery.Status = status
	if status == "delivered" {
		delivery.ActualTime = time.Now()
	}

	if err := s.repo.UpdateDelivery(ctx, delivery); err != nil {
		return fmt.Errorf("配送更新エラー: %v", err)
	}

	// ステータス更新の通知
	if s.notifyService != nil {
		if err := s.notifyService.NotifyDeliveryStatusChange(ctx, delivery); err != nil {
			// 通知エラーはログに記録するだけで、ステータス更新自体は成功とする
			fmt.Printf("通知エラー: %v\n", err)
		}
	}

	return nil
}

// CreateDeliveryTracking 配送追跡を作成する
func (s *DeliveryService) CreateDeliveryTracking(ctx context.Context, req *models.CreateTrackingRequest) (*models.DeliveryTracking, error) {
	tracking := &models.DeliveryTracking{
		DeliveryID: req.DeliveryID,
		Location:   req.Location,
		Status:     req.Status,
		Notes:      req.Notes,
	}

	if err := s.repo.CreateDeliveryTracking(ctx, tracking); err != nil {
		return nil, fmt.Errorf("配送追跡作成エラー: %v", err)
	}

	// 配送追跡の通知
	if s.notifyService != nil {
		if err := s.notifyService.NotifyDeliveryTracking(ctx, tracking); err != nil {
			// 通知エラーはログに記録するだけで、追跡作成自体は成功とする
			fmt.Printf("通知エラー: %v\n", err)
		}
	}

	return tracking, nil
}

// ListDeliveryTrackings 配送追跡履歴を取得する
func (s *DeliveryService) ListDeliveryTrackings(ctx context.Context, deliveryID int64) ([]*models.DeliveryTracking, error) {
	trackings, err := s.repo.ListDeliveryTrackings(ctx, deliveryID)
	if err != nil {
		return nil, fmt.Errorf("配送追跡履歴取得エラー: %v", err)
	}

	return trackings, nil
}

// CompleteDelivery 配送を完了する
func (s *DeliveryService) CompleteDelivery(ctx context.Context, id int64) error {
	delivery, err := s.repo.GetDelivery(ctx, id)
	if err != nil {
		return fmt.Errorf("配送取得エラー: %v", err)
	}

	if delivery.Status != "in_transit" {
		return fmt.Errorf("配送中の配送のみ完了できます")
	}

	// 在庫の更新
	items, err := s.repo.ListDeliveryItems(ctx, id)
	if err != nil {
		return fmt.Errorf("配送商品取得エラー: %v", err)
	}

	for _, item := range items {
		inventory, err := s.inventoryRepo.GetInventory(ctx, item.ProductID)
		if err != nil {
			return fmt.Errorf("在庫取得エラー: %v", err)
		}

		inventory.Quantity -= item.Quantity
		if err := s.inventoryRepo.UpdateInventory(ctx, inventory); err != nil {
			return fmt.Errorf("在庫更新エラー: %v", err)
		}
	}

	// 配送ステータスの更新
	delivery.Status = "delivered"
	delivery.ActualTime = time.Now()

	if err := s.repo.UpdateDelivery(ctx, delivery); err != nil {
		return fmt.Errorf("配送更新エラー: %v", err)
	}

	// 配送完了の通知
	if s.notifyService != nil {
		if err := s.notifyService.NotifyDeliveryComplete(ctx, delivery); err != nil {
			// 通知エラーはログに記録するだけで、配送完了自体は成功とする
			fmt.Printf("通知エラー: %v\n", err)
		}
	}

	return nil
}

package services

import (
	"context"
	"fmt"
	"time"

	"tea-logistics/pkg/models"
	"tea-logistics/pkg/repository"
)

/*
 * 通知サービス
 * 通知関連のビジネスロジックを実装する
 */

// NotificationService 通知サービスインターフェース
type NotificationService interface {
	CreateNotification(ctx context.Context, req *models.CreateNotificationRequest) (*models.Notification, error)
	GetNotification(ctx context.Context, id int64) (*models.Notification, error)
	ListNotifications(ctx context.Context, userID int64) ([]*models.Notification, error)
	MarkAsRead(ctx context.Context, id int64) error
	DeleteNotification(ctx context.Context, id int64) error
	NotifyDeliveryStatusChange(ctx context.Context, delivery *models.Delivery) error
	NotifyDeliveryComplete(ctx context.Context, delivery *models.Delivery) error
	NotifyDeliveryTracking(ctx context.Context, tracking *models.DeliveryTracking) error
}

// NotificationServiceImpl 通知サービス実装
type NotificationServiceImpl struct {
	repo         repository.NotificationRepository
	deliveryRepo repository.DeliveryRepository
}

// NewNotificationService 通知サービスを作成する
func NewNotificationService(repo repository.NotificationRepository, deliveryRepo repository.DeliveryRepository) NotificationService {
	return &NotificationServiceImpl{
		repo:         repo,
		deliveryRepo: deliveryRepo,
	}
}

// CreateNotification 通知を作成する
func (s *NotificationServiceImpl) CreateNotification(ctx context.Context, req *models.CreateNotificationRequest) (*models.Notification, error) {
	notification := &models.Notification{
		Type:      req.Type,
		Status:    models.NotificationStatusUnread,
		Title:     req.Title,
		Message:   req.Message,
		Data:      req.Data,
		UserID:    req.UserID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.repo.CreateNotification(ctx, notification); err != nil {
		return nil, fmt.Errorf("通知作成エラー: %v", err)
	}

	return notification, nil
}

// GetNotification 通知を取得する
func (s *NotificationServiceImpl) GetNotification(ctx context.Context, id int64) (*models.Notification, error) {
	notification, err := s.repo.GetNotification(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("通知取得エラー: %v", err)
	}

	return notification, nil
}

// ListNotifications ユーザーの通知一覧を取得する
func (s *NotificationServiceImpl) ListNotifications(ctx context.Context, userID int64) ([]*models.Notification, error) {
	notifications, err := s.repo.ListNotifications(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("通知一覧取得エラー: %v", err)
	}

	return notifications, nil
}

// MarkAsRead 通知を既読にする
func (s *NotificationServiceImpl) MarkAsRead(ctx context.Context, id int64) error {
	if err := s.repo.UpdateNotificationStatus(ctx, id, models.NotificationStatusRead); err != nil {
		return fmt.Errorf("通知ステータス更新エラー: %v", err)
	}

	return nil
}

// DeleteNotification 通知を削除する
func (s *NotificationServiceImpl) DeleteNotification(ctx context.Context, id int64) error {
	if err := s.repo.DeleteNotification(ctx, id); err != nil {
		return fmt.Errorf("通知削除エラー: %v", err)
	}

	return nil
}

// NotifyDeliveryStatusChange 配送ステータス変更を通知する
func (s *NotificationServiceImpl) NotifyDeliveryStatusChange(ctx context.Context, delivery *models.Delivery) error {
	req := &models.CreateNotificationRequest{
		Type:    models.NotificationTypeDeliveryStatus,
		Title:   "配送ステータスが更新されました",
		Message: fmt.Sprintf("配送ID: %d のステータスが「%s」に更新されました", delivery.ID, delivery.Status),
		Data: map[string]interface{}{
			"delivery_id": delivery.ID,
			"status":      delivery.Status,
		},
		UserID: delivery.OrderID, // 注文IDをユーザーIDとして使用
	}

	_, err := s.CreateNotification(ctx, req)
	if err != nil {
		return fmt.Errorf("配送ステータス変更通知エラー: %v", err)
	}

	return nil
}

// NotifyDeliveryComplete 配送完了を通知する
func (s *NotificationServiceImpl) NotifyDeliveryComplete(ctx context.Context, delivery *models.Delivery) error {
	req := &models.CreateNotificationRequest{
		Type:    models.NotificationTypeDeliveryComplete,
		Title:   "配送が完了しました",
		Message: fmt.Sprintf("配送ID: %d の配送が完了しました", delivery.ID),
		Data: map[string]interface{}{
			"delivery_id":  delivery.ID,
			"completed_at": delivery.ActualTime,
		},
		UserID: delivery.OrderID, // 注文IDをユーザーIDとして使用
	}

	_, err := s.CreateNotification(ctx, req)
	if err != nil {
		return fmt.Errorf("配送完了通知エラー: %v", err)
	}

	return nil
}

// NotifyDeliveryTracking 配送追跡を通知する
func (s *NotificationServiceImpl) NotifyDeliveryTracking(ctx context.Context, tracking *models.DeliveryTracking) error {
	// 配送情報を取得
	delivery, err := s.deliveryRepo.GetDelivery(ctx, tracking.DeliveryID)
	if err != nil {
		return fmt.Errorf("配送情報取得エラー: %v", err)
	}

	req := &models.CreateNotificationRequest{
		Type:    models.NotificationTypeDeliveryTracking,
		Title:   "配送状況が更新されました",
		Message: fmt.Sprintf("配送ID: %d の現在位置: %s", tracking.DeliveryID, tracking.Location),
		Data: map[string]interface{}{
			"delivery_id": tracking.DeliveryID,
			"location":    tracking.Location,
			"status":      tracking.Status,
		},
		UserID: delivery.OrderID, // 配送のOrderIDを使用
	}

	_, err = s.CreateNotification(ctx, req)
	if err != nil {
		return fmt.Errorf("配送追跡通知エラー: %v", err)
	}

	return nil
}

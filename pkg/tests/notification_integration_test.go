package tests

import (
	"context"
	"testing"
	"time"

	"tea-logistics/pkg/models"
	"tea-logistics/pkg/repository"
	"tea-logistics/pkg/services"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

/*
 * 通知機能の統合テスト
 * 通知の作成から取得、配送ステータス変更時の通知送信などをテストする
 */

type NotificationIntegrationTestSuite struct {
	suite.Suite
	notifyService services.NotificationService
	deliveryRepo  repository.DeliveryRepository
	notifyRepo    repository.NotificationRepository
}

func (s *NotificationIntegrationTestSuite) SetupSuite() {
	// テスト用のDBセットアップ
	db := setupTestDB()
	s.notifyRepo = repository.NewSQLNotificationRepository(db)
	s.deliveryRepo = repository.NewSQLDeliveryRepository(db)
	s.notifyService = services.NewNotificationService(s.notifyRepo, s.deliveryRepo)
}

func (s *NotificationIntegrationTestSuite) TearDownSuite() {
	// テスト用のDBクリーンアップ
	cleanupTestDB()
}

func (s *NotificationIntegrationTestSuite) TestNotificationFlow() {
	ctx := context.Background()

	// 1. 通知の作成
	req := &models.CreateNotificationRequest{
		Type:    models.NotificationTypeDeliveryStatus,
		Title:   "配送ステータス更新",
		Message: "配送が開始されました",
		Data: map[string]interface{}{
			"delivery_id": 1,
			"status":      "in_transit",
		},
		UserID: 1,
	}

	notification, err := s.notifyService.CreateNotification(ctx, req)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), notification)
	assert.Equal(s.T(), req.Type, notification.Type)
	assert.Equal(s.T(), models.NotificationStatusUnread, notification.Status)

	// 2. 通知の取得
	fetched, err := s.notifyService.GetNotification(ctx, notification.ID)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), notification.ID, fetched.ID)
	assert.Equal(s.T(), notification.Title, fetched.Title)

	// 3. 通知を既読にする
	err = s.notifyService.MarkAsRead(ctx, notification.ID)
	assert.NoError(s.T(), err)

	// 4. 既読状態の確認
	fetched, err = s.notifyService.GetNotification(ctx, notification.ID)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), models.NotificationStatusRead, fetched.Status)
}

func (s *NotificationIntegrationTestSuite) TestDeliveryStatusChangeNotification() {
	ctx := context.Background()

	// 1. 配送の作成
	delivery := &models.Delivery{
		OrderID:         1,
		Status:          "pending",
		FromWarehouseID: 1,
		ToAddress:       "東京都渋谷区",
		EstimatedTime:   time.Now().Add(24 * time.Hour),
	}
	err := s.deliveryRepo.CreateDelivery(ctx, delivery)
	assert.NoError(s.T(), err)

	// 2. 配送ステータス変更時の通知
	err = s.notifyService.NotifyDeliveryStatusChange(ctx, delivery)
	assert.NoError(s.T(), err)

	// 3. 通知の確認
	notifications, err := s.notifyService.ListNotifications(ctx, delivery.OrderID)
	assert.NoError(s.T(), err)
	assert.NotEmpty(s.T(), notifications)
	assert.Equal(s.T(), models.NotificationTypeDeliveryStatus, notifications[0].Type)
}

func (s *NotificationIntegrationTestSuite) TestDeliveryCompleteNotification() {
	ctx := context.Background()

	// 1. 配送の作成
	delivery := &models.Delivery{
		OrderID:         2,
		Status:          "delivered",
		FromWarehouseID: 1,
		ToAddress:       "東京都新宿区",
		EstimatedTime:   time.Now().Add(24 * time.Hour),
		ActualTime:      time.Now(),
	}
	err := s.deliveryRepo.CreateDelivery(ctx, delivery)
	assert.NoError(s.T(), err)

	// 2. 配送完了時の通知
	err = s.notifyService.NotifyDeliveryComplete(ctx, delivery)
	assert.NoError(s.T(), err)

	// 3. 通知の確認
	notifications, err := s.notifyService.ListNotifications(ctx, delivery.OrderID)
	assert.NoError(s.T(), err)
	assert.NotEmpty(s.T(), notifications)
	assert.Equal(s.T(), models.NotificationTypeDeliveryComplete, notifications[0].Type)
}

func (s *NotificationIntegrationTestSuite) TestDeliveryTrackingNotification() {
	ctx := context.Background()

	// 1. 配送の作成
	delivery := &models.Delivery{
		OrderID:         3,
		Status:          "in_transit",
		FromWarehouseID: 1,
		ToAddress:       "東京都渋谷区",
		EstimatedTime:   time.Now().Add(24 * time.Hour),
	}
	err := s.deliveryRepo.CreateDelivery(ctx, delivery)
	assert.NoError(s.T(), err)

	// 2. 配送追跡の作成
	tracking := &models.DeliveryTracking{
		DeliveryID: delivery.ID, // 作成した配送のIDを使用
		Location:   "東京都渋谷区",
		Status:     "配送中",
		Notes:      "順調に配送中です",
	}
	err = s.deliveryRepo.CreateDeliveryTracking(ctx, tracking)
	assert.NoError(s.T(), err)

	// 3. 配送追跡時の通知
	err = s.notifyService.NotifyDeliveryTracking(ctx, tracking)
	assert.NoError(s.T(), err)

	// 4. 通知の確認
	notifications, err := s.notifyService.ListNotifications(ctx, delivery.OrderID) // OrderIDを使用
	assert.NoError(s.T(), err)
	assert.NotEmpty(s.T(), notifications)
	assert.Equal(s.T(), models.NotificationTypeDeliveryTracking, notifications[0].Type)
}

func TestNotificationIntegration(t *testing.T) {
	suite.Run(t, new(NotificationIntegrationTestSuite))
}

package mocks

import (
	"context"
	"tea-logistics/pkg/models"

	"github.com/stretchr/testify/mock"
)

// MockNotificationService モック通知サービス
type MockNotificationService struct {
	mock.Mock
}

func (m *MockNotificationService) CreateNotification(ctx context.Context, req *models.CreateNotificationRequest) (*models.Notification, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Notification), args.Error(1)
}

func (m *MockNotificationService) GetNotification(ctx context.Context, id int64) (*models.Notification, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Notification), args.Error(1)
}

func (m *MockNotificationService) ListNotifications(ctx context.Context, userID int64) ([]*models.Notification, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Notification), args.Error(1)
}

func (m *MockNotificationService) MarkAsRead(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockNotificationService) DeleteNotification(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockNotificationService) NotifyDeliveryStatusChange(ctx context.Context, delivery *models.Delivery) error {
	args := m.Called(ctx, delivery)
	return args.Error(0)
}

func (m *MockNotificationService) NotifyDeliveryComplete(ctx context.Context, delivery *models.Delivery) error {
	args := m.Called(ctx, delivery)
	return args.Error(0)
}

func (m *MockNotificationService) NotifyDeliveryTracking(ctx context.Context, tracking *models.DeliveryTracking) error {
	args := m.Called(ctx, tracking)
	return args.Error(0)
}

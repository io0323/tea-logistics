package services

import (
	"context"
	"testing"
	"time"

	"tea-logistics/pkg/models"
	"tea-logistics/pkg/services/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

/*
 * 配送サービステスト
 * 配送関連のビジネスロジックのテストを実装する
 */

func setupTest() (*mocks.MockDeliveryRepository, *mocks.MockInventoryRepository, *mocks.MockNotificationService) {
	mockRepo := new(mocks.MockDeliveryRepository)
	mockInventoryRepo := new(mocks.MockInventoryRepository)
	mockNotifyService := new(mocks.MockNotificationService)
	return mockRepo, mockInventoryRepo, mockNotifyService
}

func TestCreateDelivery(t *testing.T) {
	mockRepo, mockInventoryRepo, mockNotifyService := setupTest()
	service := NewDeliveryService(mockRepo, mockInventoryRepo, mockNotifyService)

	ctx := context.Background()
	req := &models.CreateDeliveryRequest{
		OrderID:         1,
		ProductID:       1,
		Quantity:        10,
		FromWarehouseID: 1,
		ToAddress:       "東京都渋谷区",
		EstimatedTime:   time.Now().Add(24 * time.Hour),
	}

	inventory := &models.Inventory{
		ID:        1,
		ProductID: 1,
		Quantity:  100,
		Location:  "東京倉庫",
		Status:    models.InventoryStatusAvailable,
	}

	mockInventoryRepo.On("GetInventory", ctx, int64(1)).Return(inventory, nil)
	mockInventoryRepo.On("UpdateInventory", ctx, mock.AnythingOfType("*models.Inventory")).Return(nil)
	mockRepo.On("CreateDelivery", ctx, mock.AnythingOfType("*models.Delivery")).Return(nil)
	mockRepo.On("CreateDeliveryItem", ctx, mock.AnythingOfType("*models.DeliveryItem")).Return(nil)
	mockNotifyService.On("NotifyDeliveryStatusChange", ctx, mock.AnythingOfType("*models.Delivery")).Return(nil)

	delivery, err := service.CreateDelivery(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, delivery)
	assert.Equal(t, req.OrderID, delivery.OrderID)
	assert.Equal(t, req.FromWarehouseID, delivery.FromWarehouseID)
	assert.Equal(t, req.ToAddress, delivery.ToAddress)
	assert.Equal(t, req.EstimatedTime, delivery.EstimatedTime)

	mockRepo.AssertExpectations(t)
	mockInventoryRepo.AssertExpectations(t)
	mockNotifyService.AssertExpectations(t)
}

func TestGetDelivery(t *testing.T) {
	mockRepo, mockInventoryRepo, mockNotifyService := setupTest()
	service := NewDeliveryService(mockRepo, mockInventoryRepo, mockNotifyService)

	ctx := context.Background()
	expectedDelivery := &models.Delivery{
		ID:              1,
		OrderID:         1,
		Status:          "pending",
		FromWarehouseID: 1,
		ToAddress:       "東京都渋谷区",
		EstimatedTime:   time.Now().Add(24 * time.Hour),
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	mockRepo.On("GetDelivery", ctx, int64(1)).Return(expectedDelivery, nil)

	delivery, err := service.GetDelivery(ctx, 1)

	assert.NoError(t, err)
	assert.NotNil(t, delivery)
	assert.Equal(t, expectedDelivery.ID, delivery.ID)
	assert.Equal(t, expectedDelivery.OrderID, delivery.OrderID)
	assert.Equal(t, expectedDelivery.Status, delivery.Status)
	assert.Equal(t, expectedDelivery.FromWarehouseID, delivery.FromWarehouseID)
	assert.Equal(t, expectedDelivery.ToAddress, delivery.ToAddress)
	assert.Equal(t, expectedDelivery.EstimatedTime, delivery.EstimatedTime)

	mockRepo.AssertExpectations(t)
}

func TestUpdateDeliveryStatus(t *testing.T) {
	mockRepo, mockInventoryRepo, mockNotifyService := setupTest()
	service := NewDeliveryService(mockRepo, mockInventoryRepo, mockNotifyService)

	ctx := context.Background()
	delivery := &models.Delivery{
		ID:              1,
		OrderID:         1,
		Status:          "pending",
		FromWarehouseID: 1,
		ToAddress:       "東京都渋谷区",
		EstimatedTime:   time.Now().Add(24 * time.Hour),
	}

	mockRepo.On("GetDelivery", ctx, int64(1)).Return(delivery, nil)
	mockRepo.On("UpdateDelivery", ctx, mock.AnythingOfType("*models.Delivery")).Return(nil)
	mockNotifyService.On("NotifyDeliveryStatusChange", ctx, mock.AnythingOfType("*models.Delivery")).Return(nil)

	err := service.UpdateDeliveryStatus(ctx, 1, "in_transit")

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockNotifyService.AssertExpectations(t)
}

func TestCompleteDelivery(t *testing.T) {
	mockRepo, mockInventoryRepo, mockNotifyService := setupTest()
	service := NewDeliveryService(mockRepo, mockInventoryRepo, mockNotifyService)

	ctx := context.Background()
	delivery := &models.Delivery{
		ID:              1,
		OrderID:         1,
		Status:          "in_transit",
		FromWarehouseID: 1,
		ToAddress:       "東京都渋谷区",
		EstimatedTime:   time.Now().Add(24 * time.Hour),
	}

	items := []*models.DeliveryItem{
		{
			ID:         1,
			DeliveryID: 1,
			ProductID:  1,
			Quantity:   10,
		},
	}

	inventory := &models.Inventory{
		ID:        1,
		ProductID: 1,
		Quantity:  100,
		Location:  "東京倉庫",
		Status:    models.InventoryStatusAvailable,
	}

	mockRepo.On("GetDelivery", ctx, int64(1)).Return(delivery, nil)
	mockRepo.On("ListDeliveryItems", ctx, int64(1)).Return(items, nil)
	mockInventoryRepo.On("GetInventory", ctx, int64(1)).Return(inventory, nil)
	mockInventoryRepo.On("UpdateInventory", ctx, mock.AnythingOfType("*models.Inventory")).Return(nil)
	mockRepo.On("UpdateDelivery", ctx, mock.AnythingOfType("*models.Delivery")).Return(nil)
	mockNotifyService.On("NotifyDeliveryComplete", ctx, mock.AnythingOfType("*models.Delivery")).Return(nil)

	err := service.CompleteDelivery(ctx, 1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockInventoryRepo.AssertExpectations(t)
	mockNotifyService.AssertExpectations(t)
}

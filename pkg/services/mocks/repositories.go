package mocks

import (
	"context"

	"tea-logistics/pkg/models"
	"tea-logistics/pkg/repository"

	"github.com/stretchr/testify/mock"
)

/*
 * モックリポジトリ
 * テスト用のモックリポジトリを定義する
 */

// MockDeliveryRepository モック配送リポジトリ
type MockDeliveryRepository struct {
	mock.Mock
}

// Ensure MockDeliveryRepository implements DeliveryRepository interface
var _ repository.DeliveryRepository = (*MockDeliveryRepository)(nil)

func (m *MockDeliveryRepository) CreateDelivery(ctx context.Context, delivery *models.Delivery) error {
	args := m.Called(ctx, delivery)
	return args.Error(0)
}

func (m *MockDeliveryRepository) GetDelivery(ctx context.Context, id int64) (*models.Delivery, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Delivery), args.Error(1)
}

func (m *MockDeliveryRepository) ListDeliveries(ctx context.Context) ([]*models.Delivery, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Delivery), args.Error(1)
}

func (m *MockDeliveryRepository) UpdateDelivery(ctx context.Context, delivery *models.Delivery) error {
	args := m.Called(ctx, delivery)
	return args.Error(0)
}

func (m *MockDeliveryRepository) CreateDeliveryItem(ctx context.Context, item *models.DeliveryItem) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}

func (m *MockDeliveryRepository) ListDeliveryItems(ctx context.Context, deliveryID int64) ([]*models.DeliveryItem, error) {
	args := m.Called(ctx, deliveryID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.DeliveryItem), args.Error(1)
}

func (m *MockDeliveryRepository) CreateDeliveryTracking(ctx context.Context, tracking *models.DeliveryTracking) error {
	args := m.Called(ctx, tracking)
	return args.Error(0)
}

func (m *MockDeliveryRepository) ListDeliveryTrackings(ctx context.Context, deliveryID int64) ([]*models.DeliveryTracking, error) {
	args := m.Called(ctx, deliveryID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.DeliveryTracking), args.Error(1)
}

// MockInventoryRepository モック在庫リポジトリ
type MockInventoryRepository struct {
	mock.Mock
}

// Ensure MockInventoryRepository implements InventoryRepository interface
var _ repository.InventoryRepository = (*MockInventoryRepository)(nil)

func (m *MockInventoryRepository) CreateInventory(ctx context.Context, inventory *models.Inventory) error {
	args := m.Called(ctx, inventory)
	return args.Error(0)
}

func (m *MockInventoryRepository) GetInventory(ctx context.Context, id int64) (*models.Inventory, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Inventory), args.Error(1)
}

func (m *MockInventoryRepository) ListInventories(ctx context.Context) ([]*models.Inventory, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Inventory), args.Error(1)
}

func (m *MockInventoryRepository) UpdateInventory(ctx context.Context, inventory *models.Inventory) error {
	args := m.Called(ctx, inventory)
	return args.Error(0)
}

func (m *MockInventoryRepository) DeleteInventory(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockInventoryRepository) GetInventoryByProduct(ctx context.Context, productID int64) (*models.Inventory, error) {
	args := m.Called(ctx, productID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Inventory), args.Error(1)
}

func (m *MockInventoryRepository) UpdateQuantity(ctx context.Context, id int64, quantity int) error {
	args := m.Called(ctx, id, quantity)
	return args.Error(0)
}

func (m *MockInventoryRepository) GetInventoryByLocation(ctx context.Context, location string) ([]*models.Inventory, error) {
	args := m.Called(ctx, location)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Inventory), args.Error(1)
}

func (m *MockInventoryRepository) CreateMovement(ctx context.Context, movement *models.InventoryMovement) error {
	args := m.Called(ctx, movement)
	return args.Error(0)
}

func (m *MockInventoryRepository) ListMovements(ctx context.Context, productID int64) ([]*models.InventoryMovement, error) {
	args := m.Called(ctx, productID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.InventoryMovement), args.Error(1)
}

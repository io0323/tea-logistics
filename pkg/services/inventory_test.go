package services

import (
	"context"
	"fmt"
	"testing"
	"time"

	"tea-logistics/pkg/models"
	"tea-logistics/pkg/repository"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

/*
 * 在庫管理サービステスト
 * 在庫関連のビジネスロジックのテストを実装する
 */

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

func TestCreateInventory(t *testing.T) {
	mockRepo := new(MockInventoryRepository)
	service := NewInventoryService(mockRepo)

	ctx := context.Background()
	req := &models.CreateInventoryRequest{
		ProductID: 1,
		Quantity:  100,
		Location:  "東京倉庫",
		Status:    models.InventoryStatusAvailable,
	}

	mockRepo.On("CreateInventory", ctx, mock.AnythingOfType("*models.Inventory")).Return(nil)

	inventory, err := service.CreateInventory(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, inventory)
	assert.Equal(t, req.ProductID, inventory.ProductID)
	assert.Equal(t, req.Quantity, inventory.Quantity)
	assert.Equal(t, req.Location, inventory.Location)
	assert.Equal(t, req.Status, inventory.Status)

	mockRepo.AssertExpectations(t)
}

func TestGetInventory(t *testing.T) {
	mockRepo := new(MockInventoryRepository)
	service := NewInventoryService(mockRepo)

	ctx := context.Background()
	expectedInventory := &models.Inventory{
		ID:        1,
		ProductID: 1,
		Quantity:  100,
		Location:  "東京倉庫",
		Status:    models.InventoryStatusAvailable,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockRepo.On("GetInventory", ctx, int64(1)).Return(expectedInventory, nil)

	inventory, err := service.GetInventory(ctx, 1)

	assert.NoError(t, err)
	assert.NotNil(t, inventory)
	assert.Equal(t, expectedInventory.ID, inventory.ID)
	assert.Equal(t, expectedInventory.ProductID, inventory.ProductID)
	assert.Equal(t, expectedInventory.Quantity, inventory.Quantity)
	assert.Equal(t, expectedInventory.Location, inventory.Location)
	assert.Equal(t, expectedInventory.Status, inventory.Status)

	mockRepo.AssertExpectations(t)
}

func TestListInventories(t *testing.T) {
	mockRepo := new(MockInventoryRepository)
	service := NewInventoryService(mockRepo)

	ctx := context.Background()
	expectedInventories := []*models.Inventory{
		{
			ID:        1,
			ProductID: 1,
			Quantity:  100,
			Location:  "東京倉庫",
			Status:    models.InventoryStatusAvailable,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        2,
			ProductID: 2,
			Quantity:  50,
			Location:  "大阪倉庫",
			Status:    models.InventoryStatusAvailable,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	mockRepo.On("ListInventories", ctx).Return(expectedInventories, nil)

	inventories, err := service.ListInventories(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, inventories)
	assert.Len(t, inventories, 2)
	assert.Equal(t, expectedInventories[0].ID, inventories[0].ID)
	assert.Equal(t, expectedInventories[1].ID, inventories[1].ID)

	mockRepo.AssertExpectations(t)
}

func TestUpdateInventory(t *testing.T) {
	mockRepo := new(MockInventoryRepository)
	service := NewInventoryService(mockRepo)

	ctx := context.Background()
	existingInventory := &models.Inventory{
		ID:        1,
		ProductID: 1,
		Quantity:  100,
		Location:  "東京倉庫",
		Status:    models.InventoryStatusAvailable,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	req := &models.UpdateInventoryRequest{
		Quantity: 150,
		Location: "大阪倉庫",
		Status:   models.InventoryStatusAvailable,
	}

	mockRepo.On("GetInventory", ctx, int64(1)).Return(existingInventory, nil)
	mockRepo.On("UpdateInventory", ctx, mock.AnythingOfType("*models.Inventory")).Return(nil)

	inventory, err := service.UpdateInventory(ctx, 1, req)

	assert.NoError(t, err)
	assert.NotNil(t, inventory)
	assert.Equal(t, req.Quantity, inventory.Quantity)
	assert.Equal(t, req.Location, inventory.Location)
	assert.Equal(t, req.Status, inventory.Status)

	mockRepo.AssertExpectations(t)
}

func TestCreateMovement(t *testing.T) {
	mockRepo := new(MockInventoryRepository)
	service := NewInventoryService(mockRepo)

	ctx := context.Background()
	fromInventory := &models.Inventory{
		ID:        1,
		ProductID: 1,
		Quantity:  100,
		Location:  "東京倉庫",
		Status:    models.InventoryStatusAvailable,
	}

	req := &models.CreateMovementRequest{
		ProductID:       1,
		FromLocation:    "東京倉庫",
		ToLocation:      "大阪倉庫",
		Quantity:        50,
		MovementType:    models.MovementTypeTransfer,
		MovementDate:    time.Now(),
		ReferenceNumber: "TRF-001",
	}

	mockRepo.On("GetInventoryByProduct", ctx, int64(1)).Return(fromInventory, nil)
	mockRepo.On("UpdateQuantity", ctx, int64(1), 50).Return(nil)
	mockRepo.On("GetInventoryByProduct", ctx, int64(1)).Return(nil, fmt.Errorf("在庫が見つかりません"))
	mockRepo.On("CreateInventory", ctx, mock.AnythingOfType("*models.Inventory")).Return(nil)
	mockRepo.On("CreateMovement", ctx, mock.AnythingOfType("*models.InventoryMovement")).Return(nil)

	movement, err := service.CreateMovement(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, movement)
	assert.Equal(t, req.ProductID, movement.ProductID)
	assert.Equal(t, req.FromLocation, movement.FromLocation)
	assert.Equal(t, req.ToLocation, movement.ToLocation)
	assert.Equal(t, req.Quantity, movement.Quantity)
	assert.Equal(t, req.MovementType, movement.MovementType)
	assert.Equal(t, req.ReferenceNumber, movement.ReferenceNumber)

	mockRepo.AssertExpectations(t)
}

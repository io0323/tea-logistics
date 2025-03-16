package services

import (
	"context"
	"testing"
	"time"

	"tea-logistics/pkg/models"
	"tea-logistics/pkg/repository"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

/*
 * 商品サービステスト
 * 商品関連のビジネスロジックのテストを実装する
 */

// MockProductRepository モック商品リポジトリ
type MockProductRepository struct {
	mock.Mock
}

// Ensure MockProductRepository implements ProductRepository interface
var _ repository.ProductRepository = (*MockProductRepository)(nil)

func (m *MockProductRepository) CreateProduct(ctx context.Context, product *models.Product) error {
	args := m.Called(ctx, product)
	return args.Error(0)
}

func (m *MockProductRepository) GetProduct(ctx context.Context, id int64) (*models.Product, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Product), args.Error(1)
}

func (m *MockProductRepository) ListProducts(ctx context.Context) ([]*models.Product, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Product), args.Error(1)
}

func (m *MockProductRepository) UpdateProduct(ctx context.Context, product *models.Product) error {
	args := m.Called(ctx, product)
	return args.Error(0)
}

func (m *MockProductRepository) DeleteProduct(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestCreateProduct(t *testing.T) {
	mockRepo := new(MockProductRepository)
	service := NewProductService(mockRepo)

	ctx := context.Background()
	req := &models.CreateProductRequest{
		Name:        "テスト商品",
		Description: "テスト商品の説明",
		SKU:         "TEST-001",
		Category:    "テストカテゴリ",
		Price:       1000,
		Status:      models.ProductStatusActive,
	}

	mockRepo.On("CreateProduct", ctx, mock.AnythingOfType("*models.Product")).Return(nil)

	product, err := service.CreateProduct(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, product)
	assert.Equal(t, req.Name, product.Name)
	assert.Equal(t, req.Description, product.Description)
	assert.Equal(t, req.SKU, product.SKU)
	assert.Equal(t, req.Category, product.Category)
	assert.Equal(t, req.Price, product.Price)
	assert.Equal(t, req.Status, product.Status)

	mockRepo.AssertExpectations(t)
}

func TestGetProduct(t *testing.T) {
	mockRepo := new(MockProductRepository)
	service := NewProductService(mockRepo)

	ctx := context.Background()
	expectedProduct := &models.Product{
		ID:          1,
		Name:        "テスト商品",
		Description: "テスト商品の説明",
		SKU:         "TEST-001",
		Category:    "テストカテゴリ",
		Price:       1000,
		Status:      models.ProductStatusActive,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	mockRepo.On("GetProduct", ctx, int64(1)).Return(expectedProduct, nil)

	product, err := service.GetProduct(ctx, 1)

	assert.NoError(t, err)
	assert.NotNil(t, product)
	assert.Equal(t, expectedProduct.ID, product.ID)
	assert.Equal(t, expectedProduct.Name, product.Name)
	assert.Equal(t, expectedProduct.Description, product.Description)
	assert.Equal(t, expectedProduct.SKU, product.SKU)
	assert.Equal(t, expectedProduct.Category, product.Category)
	assert.Equal(t, expectedProduct.Price, product.Price)
	assert.Equal(t, expectedProduct.Status, product.Status)

	mockRepo.AssertExpectations(t)
}

func TestListProducts(t *testing.T) {
	mockRepo := new(MockProductRepository)
	service := NewProductService(mockRepo)

	ctx := context.Background()
	expectedProducts := []*models.Product{
		{
			ID:          1,
			Name:        "テスト商品1",
			Description: "テスト商品1の説明",
			SKU:         "TEST-001",
			Category:    "テストカテゴリ",
			Price:       1000,
			Status:      models.ProductStatusActive,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          2,
			Name:        "テスト商品2",
			Description: "テスト商品2の説明",
			SKU:         "TEST-002",
			Category:    "テストカテゴリ",
			Price:       2000,
			Status:      models.ProductStatusActive,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	mockRepo.On("ListProducts", ctx).Return(expectedProducts, nil)

	products, err := service.ListProducts(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, products)
	assert.Len(t, products, 2)
	assert.Equal(t, expectedProducts[0].ID, products[0].ID)
	assert.Equal(t, expectedProducts[1].ID, products[1].ID)

	mockRepo.AssertExpectations(t)
}

func TestUpdateProduct(t *testing.T) {
	mockRepo := new(MockProductRepository)
	service := NewProductService(mockRepo)

	ctx := context.Background()
	req := &models.UpdateProductRequest{
		Name:        "更新商品",
		Description: "更新商品の説明",
		SKU:         "TEST-003",
		Category:    "更新カテゴリ",
		Price:       3000,
		Status:      models.ProductStatusActive,
	}

	existingProduct := &models.Product{
		ID:          1,
		Name:        "テスト商品",
		Description: "テスト商品の説明",
		SKU:         "TEST-001",
		Category:    "テストカテゴリ",
		Price:       1000,
		Status:      models.ProductStatusActive,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	mockRepo.On("GetProduct", ctx, int64(1)).Return(existingProduct, nil)
	mockRepo.On("UpdateProduct", ctx, mock.AnythingOfType("*models.Product")).Return(nil)

	product, err := service.UpdateProduct(ctx, 1, req)

	assert.NoError(t, err)
	assert.NotNil(t, product)
	assert.Equal(t, req.Name, product.Name)
	assert.Equal(t, req.Description, product.Description)
	assert.Equal(t, req.SKU, product.SKU)
	assert.Equal(t, req.Category, product.Category)
	assert.Equal(t, req.Price, product.Price)
	assert.Equal(t, req.Status, product.Status)

	mockRepo.AssertExpectations(t)
}

func TestDeleteProduct(t *testing.T) {
	mockRepo := new(MockProductRepository)
	service := NewProductService(mockRepo)

	ctx := context.Background()

	mockRepo.On("DeleteProduct", ctx, int64(1)).Return(nil)

	err := service.DeleteProduct(ctx, 1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

package services

import (
	"context"
	"fmt"

	"tea-logistics/pkg/models"
	"tea-logistics/pkg/repository"
)

/*
 * 商品サービス
 * 商品関連のビジネスロジックを実装する
 */

// ProductService 商品サービス
type ProductService struct {
	repo repository.ProductRepository
}

// NewProductService 商品サービスを作成する
func NewProductService(repo repository.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

// CreateProduct 商品を作成する
func (s *ProductService) CreateProduct(ctx context.Context, req *models.CreateProductRequest) (*models.Product, error) {
	product := &models.Product{
		Name:        req.Name,
		Description: req.Description,
		SKU:         req.SKU,
		Category:    req.Category,
		Price:       req.Price,
		Status:      req.Status,
	}

	if err := s.repo.CreateProduct(ctx, product); err != nil {
		return nil, fmt.Errorf("商品作成エラー: %v", err)
	}

	return product, nil
}

// GetProduct 商品を取得する
func (s *ProductService) GetProduct(ctx context.Context, id int64) (*models.Product, error) {
	product, err := s.repo.GetProduct(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("商品取得エラー: %v", err)
	}

	return product, nil
}

// ListProducts 商品一覧を取得する
func (s *ProductService) ListProducts(ctx context.Context) ([]*models.Product, error) {
	products, err := s.repo.ListProducts(ctx)
	if err != nil {
		return nil, fmt.Errorf("商品一覧取得エラー: %v", err)
	}

	return products, nil
}

// UpdateProduct 商品を更新する
func (s *ProductService) UpdateProduct(ctx context.Context, id int64, req *models.UpdateProductRequest) (*models.Product, error) {
	product, err := s.repo.GetProduct(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("商品取得エラー: %v", err)
	}

	product.Name = req.Name
	product.Description = req.Description
	product.SKU = req.SKU
	product.Category = req.Category
	product.Price = req.Price
	product.Status = req.Status

	if err := s.repo.UpdateProduct(ctx, product); err != nil {
		return nil, fmt.Errorf("商品更新エラー: %v", err)
	}

	return product, nil
}

// DeleteProduct 商品を削除する
func (s *ProductService) DeleteProduct(ctx context.Context, id int64) error {
	if err := s.repo.DeleteProduct(ctx, id); err != nil {
		return fmt.Errorf("商品削除エラー: %v", err)
	}

	return nil
}

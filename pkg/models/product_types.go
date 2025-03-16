package models

import (
	"time"
)

/*
 * 商品モデル
 * 商品関連のデータ構造を定義する
 */

// ProductStatus 商品ステータス
type ProductStatus string

const (
	// ProductStatusActive アクティブ
	ProductStatusActive ProductStatus = "active"
	// ProductStatusInactive 非アクティブ
	ProductStatusInactive ProductStatus = "inactive"
	// ProductStatusDiscontinued 廃盤
	ProductStatusDiscontinued ProductStatus = "discontinued"
)

// Product 商品情報
type Product struct {
	ID          int64         `json:"id" db:"id"`
	Name        string        `json:"name" db:"name"`
	Description string        `json:"description" db:"description"`
	SKU         string        `json:"sku" db:"sku"`
	Category    string        `json:"category" db:"category"`
	Price       float64       `json:"price" db:"price"`
	Status      ProductStatus `json:"status" db:"status"`
	CreatedAt   time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at" db:"updated_at"`
}

// CreateProductRequest 商品作成リクエスト
type CreateProductRequest struct {
	Name        string        `json:"name" binding:"required"`
	Description string        `json:"description" binding:"required"`
	SKU         string        `json:"sku" binding:"required"`
	Category    string        `json:"category" binding:"required"`
	Price       float64       `json:"price" binding:"required,gt=0"`
	Status      ProductStatus `json:"status" binding:"required,oneof=active inactive discontinued"`
}

// UpdateProductRequest 商品更新リクエスト
type UpdateProductRequest struct {
	Name        string        `json:"name" binding:"required"`
	Description string        `json:"description" binding:"required"`
	SKU         string        `json:"sku" binding:"required"`
	Category    string        `json:"category" binding:"required"`
	Price       float64       `json:"price" binding:"required,gt=0"`
	Status      ProductStatus `json:"status" binding:"required,oneof=active inactive discontinued"`
}

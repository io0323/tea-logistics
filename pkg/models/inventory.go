package models

import (
	"time"
)

/*
 * 在庫管理モデル
 * 商品の在庫情報を管理する
 */

// InventoryStatus 在庫ステータス
type InventoryStatus string

const (
	// InventoryStatusAvailable 利用可能
	InventoryStatusAvailable InventoryStatus = "available"
	// InventoryStatusOutOfStock 在庫切れ
	InventoryStatusOutOfStock InventoryStatus = "out_of_stock"
	// InventoryStatusReserved 予約済み
	InventoryStatusReserved InventoryStatus = "reserved"
	// InventoryStatusDiscontinued 取扱終了
	InventoryStatusDiscontinued InventoryStatus = "discontinued"
)

// MovementType 在庫移動タイプ
type MovementType string

const (
	// MovementTypeInbound 入庫
	MovementTypeInbound MovementType = "inbound"
	// MovementTypeOutbound 出庫
	MovementTypeOutbound MovementType = "outbound"
	// MovementTypeTransfer 移動
	MovementTypeTransfer MovementType = "transfer"
	// MovementTypeAdjustment 調整
	MovementTypeAdjustment MovementType = "adjustment"
)

// Inventory 在庫情報
type Inventory struct {
	ID        int64           `json:"id"`
	ProductID int64           `json:"product_id"`
	Quantity  int             `json:"quantity"`
	Location  string          `json:"location"`
	Status    InventoryStatus `json:"status"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

// InventoryMovement 在庫移動履歴
type InventoryMovement struct {
	ID              int64        `json:"id"`
	ProductID       int64        `json:"product_id"`
	FromLocation    string       `json:"from_location"`
	ToLocation      string       `json:"to_location"`
	Quantity        int          `json:"quantity"`
	MovementType    MovementType `json:"movement_type"`
	MovementDate    time.Time    `json:"movement_date"`
	ReferenceNumber string       `json:"reference_number"`
	CreatedAt       time.Time    `json:"created_at"`
}

// Warehouse 倉庫情報
type Warehouse struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Address   string    `json:"address"`
	Capacity  int       `json:"capacity"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateInventoryRequest 在庫作成リクエスト
type CreateInventoryRequest struct {
	ProductID int64           `json:"product_id" binding:"required"`
	Quantity  int             `json:"quantity" binding:"required,min=0"`
	Location  string          `json:"location" binding:"required"`
	Status    InventoryStatus `json:"status" binding:"required"`
}

// UpdateInventoryRequest 在庫更新リクエスト
type UpdateInventoryRequest struct {
	Quantity int             `json:"quantity" binding:"required,min=0"`
	Location string          `json:"location" binding:"required"`
	Status   InventoryStatus `json:"status" binding:"required"`
}

// CreateMovementRequest 在庫移動作成リクエスト
type CreateMovementRequest struct {
	ProductID       int64        `json:"product_id" binding:"required"`
	FromLocation    string       `json:"from_location" binding:"required"`
	ToLocation      string       `json:"to_location" binding:"required"`
	Quantity        int          `json:"quantity" binding:"required,min=1"`
	MovementType    MovementType `json:"movement_type" binding:"required"`
	MovementDate    time.Time    `json:"movement_date" binding:"required"`
	ReferenceNumber string       `json:"reference_number" binding:"required"`
}

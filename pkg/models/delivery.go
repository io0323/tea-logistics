package models

import (
	"time"
)

/*
 * 配送管理モデル
 * 配送関連のデータ構造を定義する
 */

// DeliveryStatus 配送ステータス
type DeliveryStatus string

const (
	// DeliveryStatusPending 配送待ち
	DeliveryStatusPending DeliveryStatus = "pending"
	// DeliveryStatusScheduled 配送予定
	DeliveryStatusScheduled DeliveryStatus = "scheduled"
	// DeliveryStatusInTransit 配送中
	DeliveryStatusInTransit DeliveryStatus = "in_transit"
	// DeliveryStatusDelivered 配送完了
	DeliveryStatusDelivered DeliveryStatus = "delivered"
	// DeliveryStatusCancelled キャンセル
	DeliveryStatusCancelled DeliveryStatus = "cancelled"
)

// DeliveryOrder 配送オーダー
type DeliveryOrder struct {
	ID            int64          `json:"id"`
	OrderNumber   string         `json:"order_number"`
	CustomerID    int64          `json:"customer_id"`
	ProductID     int64          `json:"product_id"`
	Quantity      int            `json:"quantity"`
	FromLocation  string         `json:"from_location"`
	ToLocation    string         `json:"to_location"`
	Status        DeliveryStatus `json:"status"`
	ScheduledDate time.Time      `json:"scheduled_date"`
	DeliveryDate  *time.Time     `json:"delivery_date,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
}

// DeliverySchedule 配送スケジュール
type DeliverySchedule struct {
	ID         int64     `json:"id"`
	DeliveryID int64     `json:"delivery_id"`
	DriverID   int64     `json:"driver_id"`
	VehicleID  string    `json:"vehicle_id"`
	StartTime  time.Time `json:"start_time"`
	EndTime    time.Time `json:"end_time"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// DeliveryTracking 配送追跡情報
type DeliveryTracking struct {
	ID         int64     `json:"id"`
	DeliveryID int64     `json:"delivery_id"`
	Location   string    `json:"location"`
	Status     string    `json:"status"`
	Notes      string    `json:"notes"`
	CreatedAt  time.Time `json:"created_at"`
}

// Delivery 配送情報
type Delivery struct {
	ID              int64     `json:"id"`
	OrderID         int64     `json:"order_id"`
	Status          string    `json:"status"`
	FromWarehouseID int64     `json:"from_warehouse_id"`
	ToAddress       string    `json:"to_address"`
	EstimatedTime   time.Time `json:"estimated_time"`
	ActualTime      time.Time `json:"actual_time"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// DeliveryItem 配送商品情報
type DeliveryItem struct {
	ID         int64 `json:"id"`
	DeliveryID int64 `json:"delivery_id"`
	ProductID  int64 `json:"product_id"`
	Quantity   int   `json:"quantity"`
}

// Route 配送ルート情報
type Route struct {
	ID          int64     `json:"id"`
	DeliveryID  int64     `json:"delivery_id"`
	Sequence    int       `json:"sequence"`
	Location    string    `json:"location"`
	ArrivalTime time.Time `json:"arrival_time"`
	Distance    float64   `json:"distance"`
	Duration    int       `json:"duration"` // 分単位
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Vehicle 配送車両情報
type Vehicle struct {
	ID            int64     `json:"id"`
	VehicleNumber string    `json:"vehicle_number"`
	Type          string    `json:"type"`
	Capacity      float64   `json:"capacity"`
	Status        string    `json:"status"`
	LastLocation  string    `json:"last_location"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// CreateDeliveryRequest 配送作成リクエスト
type CreateDeliveryRequest struct {
	OrderID         int64     `json:"order_id" binding:"required"`
	ProductID       int64     `json:"product_id" binding:"required"`
	Quantity        int       `json:"quantity" binding:"required,min=1"`
	FromWarehouseID int64     `json:"from_warehouse_id" binding:"required"`
	ToAddress       string    `json:"to_address" binding:"required"`
	EstimatedTime   time.Time `json:"estimated_time" binding:"required"`
}

// CreateTrackingRequest 配送追跡作成リクエスト
type CreateTrackingRequest struct {
	DeliveryID int64  `json:"delivery_id" binding:"required"`
	Location   string `json:"location" binding:"required"`
	Status     string `json:"status" binding:"required"`
	Notes      string `json:"notes" binding:"required"`
}

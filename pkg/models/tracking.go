package models

import (
	"time"
)

/*
 * 配送追跡モデル
 * 配送追跡に関連するデータ構造を定義する
 */

// TrackingStatus 追跡ステータス
type TrackingStatus string

const (
	// TrackingStatusRegistered 登録済み
	TrackingStatusRegistered TrackingStatus = "registered"
	// TrackingStatusInTransit 輸送中
	TrackingStatusInTransit TrackingStatus = "in_transit"
	// TrackingStatusDelivered 配送完了
	TrackingStatusDelivered TrackingStatus = "delivered"
	// TrackingStatusException 例外発生
	TrackingStatusException TrackingStatus = "exception"
)

// TrackingEvent 追跡イベント
type TrackingEvent struct {
	ID          int64          `json:"id"`
	TrackingID  string         `json:"tracking_id"`
	Status      TrackingStatus `json:"status"`
	Location    string         `json:"location"`
	Description string         `json:"description"`
	Latitude    float64        `json:"latitude,omitempty"`
	Longitude   float64        `json:"longitude,omitempty"`
	Temperature float64        `json:"temperature,omitempty"`
	Humidity    float64        `json:"humidity,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
}

// TrackingInfo 追跡情報
type TrackingInfo struct {
	ID              string           `json:"id"`
	DeliveryID      int64            `json:"delivery_id"`
	Status          TrackingStatus   `json:"status"`
	CurrentLocation string           `json:"current_location"`
	EstimatedTime   *time.Time       `json:"estimated_time,omitempty"`
	Events          []*TrackingEvent `json:"events"`
	CreatedAt       time.Time        `json:"created_at"`
	UpdatedAt       time.Time        `json:"updated_at"`
}

// TrackingException 追跡例外
type TrackingException struct {
	ID          int64      `json:"id"`
	TrackingID  string     `json:"tracking_id"`
	Type        string     `json:"type"`
	Description string     `json:"description"`
	Location    string     `json:"location"`
	ResolvedAt  *time.Time `json:"resolved_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
}

// TrackingCondition 追跡条件
type TrackingCondition struct {
	ID             int64     `json:"id"`
	TrackingID     string    `json:"tracking_id"`
	MinTemperature float64   `json:"min_temperature"`
	MaxTemperature float64   `json:"max_temperature"`
	MinHumidity    float64   `json:"min_humidity"`
	MaxHumidity    float64   `json:"max_humidity"`
	CheckInterval  int       `json:"check_interval"` // 分単位
	NotifyEmail    string    `json:"notify_email"`
	NotifyPhone    string    `json:"notify_phone"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

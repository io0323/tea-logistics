package models

import "time"

/*
 * 通知モデル
 * システム内の通知関連のモデルを定義する
 */

// NotificationType 通知タイプ
type NotificationType string

const (
	// NotificationTypeDeliveryStatus 配送ステータス変更通知
	NotificationTypeDeliveryStatus NotificationType = "delivery_status"
	// NotificationTypeDeliveryComplete 配送完了通知
	NotificationTypeDeliveryComplete NotificationType = "delivery_complete"
	// NotificationTypeDeliveryTracking 配送追跡通知
	NotificationTypeDeliveryTracking NotificationType = "delivery_tracking"
)

// NotificationStatus 通知ステータス
type NotificationStatus string

const (
	// NotificationStatusUnread 未読
	NotificationStatusUnread NotificationStatus = "unread"
	// NotificationStatusRead 既読
	NotificationStatusRead NotificationStatus = "read"
)

// Notification 通知
type Notification struct {
	ID        int64                  `json:"id" db:"id"`
	Type      NotificationType       `json:"type" db:"type"`
	Status    NotificationStatus     `json:"status" db:"status"`
	Title     string                 `json:"title" db:"title"`
	Message   string                 `json:"message" db:"message"`
	Data      map[string]interface{} `json:"data" db:"data"`
	UserID    int64                  `json:"user_id" db:"user_id"`
	CreatedAt time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt time.Time              `json:"updated_at" db:"updated_at"`
}

// CreateNotificationRequest 通知作成リクエスト
type CreateNotificationRequest struct {
	Type    NotificationType       `json:"type" binding:"required"`
	Title   string                 `json:"title" binding:"required"`
	Message string                 `json:"message" binding:"required"`
	Data    map[string]interface{} `json:"data"`
	UserID  int64                  `json:"user_id" binding:"required"`
}

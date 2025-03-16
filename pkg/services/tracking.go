package services

import (
	"context"
	"fmt"
	"time"

	"tea-logistics/pkg/models"
	"tea-logistics/pkg/repository"
)

/*
 * 配送追跡サービス
 * 配送追跡関連のビジネスロジックを実装する
 */

// TrackingService 配送追跡サービス
type TrackingService struct {
	trackingRepo *repository.TrackingRepository
}

// NewTrackingService 配送追跡サービスを作成する
func NewTrackingService(trackingRepo *repository.TrackingRepository) *TrackingService {
	return &TrackingService{trackingRepo: trackingRepo}
}

// InitializeTracking 配送追跡を初期化する
func (s *TrackingService) InitializeTracking(ctx context.Context, deliveryID int64, fromLocation string) (*models.TrackingInfo, error) {
	tracking := &models.TrackingInfo{
		ID:              fmt.Sprintf("TRK-%d", time.Now().Unix()),
		DeliveryID:      deliveryID,
		Status:          models.TrackingStatusRegistered,
		CurrentLocation: fromLocation,
		Events:          make([]*models.TrackingEvent, 0),
	}

	if err := s.trackingRepo.CreateTracking(ctx, tracking); err != nil {
		return nil, err
	}

	// 初期イベントの追加
	event := &models.TrackingEvent{
		TrackingID:  tracking.ID,
		Status:      models.TrackingStatusRegistered,
		Location:    fromLocation,
		Description: "配送追跡が開始されました",
	}
	if err := s.trackingRepo.AddTrackingEvent(ctx, event); err != nil {
		return nil, err
	}

	return tracking, nil
}

// UpdateTrackingStatus 配送追跡ステータスを更新する
func (s *TrackingService) UpdateTrackingStatus(ctx context.Context, trackingID string, status models.TrackingStatus, location string, description string) error {
	// ステータスの更新
	if err := s.trackingRepo.UpdateTrackingStatus(ctx, trackingID, status, location); err != nil {
		return err
	}

	// イベントの追加
	event := &models.TrackingEvent{
		TrackingID:  trackingID,
		Status:      status,
		Location:    location,
		Description: description,
	}
	if err := s.trackingRepo.AddTrackingEvent(ctx, event); err != nil {
		return err
	}

	return nil
}

// AddTrackingEvent 配送追跡イベントを追加する
func (s *TrackingService) AddTrackingEvent(ctx context.Context, event *models.TrackingEvent) error {
	// 追跡情報の存在確認
	tracking, err := s.trackingRepo.GetTracking(ctx, event.TrackingID)
	if err != nil {
		return err
	}

	// イベントの追加
	if err := s.trackingRepo.AddTrackingEvent(ctx, event); err != nil {
		return err
	}

	// 条件チェック
	condition, err := s.trackingRepo.GetTrackingCondition(ctx, tracking.ID)
	if err == nil {
		// 温度チェック
		if event.Temperature < condition.MinTemperature || event.Temperature > condition.MaxTemperature {
			// TODO: 通知を送信する
			return fmt.Errorf("温度が範囲外です: %.2f°C", event.Temperature)
		}

		// 湿度チェック
		if event.Humidity < condition.MinHumidity || event.Humidity > condition.MaxHumidity {
			// TODO: 通知を送信する
			return fmt.Errorf("湿度が範囲外です: %.2f%%", event.Humidity)
		}
	}

	return nil
}

// GetTrackingInfo 配送追跡情報を取得する
func (s *TrackingService) GetTrackingInfo(ctx context.Context, trackingID string) (*models.TrackingInfo, error) {
	return s.trackingRepo.GetTracking(ctx, trackingID)
}

// SetTrackingCondition 追跡条件を設定する
func (s *TrackingService) SetTrackingCondition(ctx context.Context, condition *models.TrackingCondition) error {
	// 追跡情報の存在確認
	if _, err := s.trackingRepo.GetTracking(ctx, condition.TrackingID); err != nil {
		return err
	}

	// 条件の作成
	if err := s.trackingRepo.CreateTrackingCondition(ctx, condition); err != nil {
		return err
	}

	return nil
}

// GetTrackingCondition 追跡条件を取得する
func (s *TrackingService) GetTrackingCondition(ctx context.Context, trackingID string) (*models.TrackingCondition, error) {
	return s.trackingRepo.GetTrackingCondition(ctx, trackingID)
}

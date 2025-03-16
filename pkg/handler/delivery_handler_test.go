package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"tea-logistics/pkg/models"
	"tea-logistics/pkg/services"
	"tea-logistics/pkg/services/mocks"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

/*
 * 配送ハンドラーテスト
 * 配送関連のHTTP APIエンドポイントのテストを実装する
 */

func setupDeliveryTest() (*gin.Engine, *mocks.MockDeliveryRepository, *mocks.MockInventoryRepository, services.NotificationService) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	mockDeliveryRepo := new(mocks.MockDeliveryRepository)
	mockInventoryRepo := new(mocks.MockInventoryRepository)
	mockNotifyService := new(mocks.MockNotificationService)
	service := services.NewDeliveryService(mockDeliveryRepo, mockInventoryRepo, mockNotifyService)
	handler := NewDeliveryHandler(service)
	handler.RegisterRoutes(router)

	return router, mockDeliveryRepo, mockInventoryRepo, mockNotifyService
}

func TestCreateDelivery(t *testing.T) {
	router, mockDeliveryRepo, mockInventoryRepo, mockNotifyService := setupDeliveryTest()

	req := &models.CreateDeliveryRequest{
		OrderID:         1,
		ProductID:       1,
		Quantity:        10,
		FromWarehouseID: 1,
		ToAddress:       "東京都渋谷区",
		EstimatedTime:   time.Now().Add(24 * time.Hour),
	}

	// 在庫の設定
	inventory := &models.Inventory{
		ID:        1,
		ProductID: 1,
		Quantity:  100,
		Location:  "東京倉庫",
		Status:    models.InventoryStatusAvailable,
	}

	mockInventoryRepo.On("GetInventory", mock.Anything, int64(1)).Return(inventory, nil)
	mockInventoryRepo.On("UpdateInventory", mock.Anything, mock.AnythingOfType("*models.Inventory")).Return(nil)
	mockDeliveryRepo.On("CreateDelivery", mock.Anything, mock.AnythingOfType("*models.Delivery")).Return(nil)
	mockDeliveryRepo.On("CreateDeliveryItem", mock.Anything, mock.AnythingOfType("*models.DeliveryItem")).Return(nil)
	mockNotifyService.(*mocks.MockNotificationService).On("NotifyDeliveryStatusChange", mock.Anything, mock.AnythingOfType("*models.Delivery")).Return(nil)

	body, _ := json.Marshal(req)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/api/deliveries", bytes.NewBuffer(body))
	r.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, r)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockDeliveryRepo.AssertExpectations(t)
	mockInventoryRepo.AssertExpectations(t)
	mockNotifyService.(*mocks.MockNotificationService).AssertExpectations(t)
}

func TestGetDelivery(t *testing.T) {
	router, mockDeliveryRepo, _, _ := setupDeliveryTest()

	expectedDelivery := &models.Delivery{
		ID:              1,
		OrderID:         1,
		Status:          "pending",
		FromWarehouseID: 1,
		ToAddress:       "東京都渋谷区",
		EstimatedTime:   time.Now().Add(24 * time.Hour),
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	mockDeliveryRepo.On("GetDelivery", mock.Anything, int64(1)).Return(expectedDelivery, nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/api/deliveries/1", nil)
	router.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.Delivery
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedDelivery.ID, response.ID)
	assert.Equal(t, expectedDelivery.Status, response.Status)

	mockDeliveryRepo.AssertExpectations(t)
}

func TestUpdateDeliveryStatus(t *testing.T) {
	router, mockDeliveryRepo, _, mockNotifyService := setupDeliveryTest()

	delivery := &models.Delivery{
		ID:              1,
		OrderID:         1,
		Status:          "pending",
		FromWarehouseID: 1,
		ToAddress:       "東京都渋谷区",
		EstimatedTime:   time.Now().Add(24 * time.Hour),
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	mockDeliveryRepo.On("GetDelivery", mock.Anything, int64(1)).Return(delivery, nil)
	mockDeliveryRepo.On("UpdateDelivery", mock.Anything, mock.AnythingOfType("*models.Delivery")).Return(nil)
	mockNotifyService.(*mocks.MockNotificationService).On("NotifyDeliveryStatusChange", mock.Anything, mock.AnythingOfType("*models.Delivery")).Return(nil)

	req := struct {
		Status string `json:"status"`
	}{
		Status: "in_transit",
	}

	body, _ := json.Marshal(req)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPut, "/api/deliveries/1/status", bytes.NewBuffer(body))
	r.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	mockDeliveryRepo.AssertExpectations(t)
	mockNotifyService.(*mocks.MockNotificationService).AssertExpectations(t)
}

func TestCompleteDelivery(t *testing.T) {
	router, mockDeliveryRepo, mockInventoryRepo, mockNotifyService := setupDeliveryTest()

	delivery := &models.Delivery{
		ID:              1,
		OrderID:         1,
		Status:          "in_transit",
		FromWarehouseID: 1,
		ToAddress:       "東京都渋谷区",
		EstimatedTime:   time.Now().Add(24 * time.Hour),
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	items := []*models.DeliveryItem{
		{
			ID:         1,
			DeliveryID: 1,
			ProductID:  1,
			Quantity:   10,
		},
	}

	inventory := &models.Inventory{
		ID:        1,
		ProductID: 1,
		Quantity:  100,
		Location:  "東京倉庫",
		Status:    models.InventoryStatusAvailable,
	}

	mockDeliveryRepo.On("GetDelivery", mock.Anything, int64(1)).Return(delivery, nil)
	mockDeliveryRepo.On("ListDeliveryItems", mock.Anything, int64(1)).Return(items, nil)
	mockInventoryRepo.On("GetInventory", mock.Anything, int64(1)).Return(inventory, nil)
	mockInventoryRepo.On("UpdateInventory", mock.Anything, mock.AnythingOfType("*models.Inventory")).Return(nil)
	mockDeliveryRepo.On("UpdateDelivery", mock.Anything, mock.AnythingOfType("*models.Delivery")).Return(nil)
	mockNotifyService.(*mocks.MockNotificationService).On("NotifyDeliveryComplete", mock.Anything, mock.AnythingOfType("*models.Delivery")).Return(nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/api/deliveries/1/complete", nil)
	router.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	mockDeliveryRepo.AssertExpectations(t)
	mockInventoryRepo.AssertExpectations(t)
	mockNotifyService.(*mocks.MockNotificationService).AssertExpectations(t)
}

func TestCreateDeliveryTracking(t *testing.T) {
	router, mockDeliveryRepo, _, mockNotifyService := setupDeliveryTest()

	req := &models.CreateTrackingRequest{
		DeliveryID: 1,
		Location:   "東京都渋谷区",
		Status:     "配送中",
		Notes:      "順調に配送中です",
	}

	mockDeliveryRepo.On("CreateDeliveryTracking", mock.Anything, mock.AnythingOfType("*models.DeliveryTracking")).Return(nil)
	mockNotifyService.(*mocks.MockNotificationService).On("NotifyDeliveryTracking", mock.Anything, mock.AnythingOfType("*models.DeliveryTracking")).Return(nil)

	body, _ := json.Marshal(req)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/api/deliveries/1/tracking", bytes.NewBuffer(body))
	r.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, r)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockDeliveryRepo.AssertExpectations(t)
	mockNotifyService.(*mocks.MockNotificationService).AssertExpectations(t)
}

func TestListDeliveryTrackings(t *testing.T) {
	router, mockDeliveryRepo, _, _ := setupDeliveryTest()

	expectedTrackings := []*models.DeliveryTracking{
		{
			ID:         1,
			DeliveryID: 1,
			Location:   "東京都渋谷区",
			Status:     "配送中",
			Notes:      "順調に配送中です",
			CreatedAt:  time.Now(),
		},
	}

	mockDeliveryRepo.On("ListDeliveryTrackings", mock.Anything, int64(1)).Return(expectedTrackings, nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/api/deliveries/1/tracking", nil)
	router.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code)

	var response []*models.DeliveryTracking
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response, 1)
	assert.Equal(t, expectedTrackings[0].ID, response[0].ID)
	assert.Equal(t, expectedTrackings[0].Status, response[0].Status)

	mockDeliveryRepo.AssertExpectations(t)
}

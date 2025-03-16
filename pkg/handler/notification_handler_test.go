package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"tea-logistics/pkg/models"
	"tea-logistics/pkg/services/mocks"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

/*
 * 通知ハンドラーのテスト
 */

func TestCreateNotification(t *testing.T) {
	mockService := new(mocks.MockNotificationService)
	handler := NewNotificationHandler(mockService)

	req := &models.CreateNotificationRequest{
		Type:    models.NotificationTypeDeliveryStatus,
		Title:   "テスト通知",
		Message: "これはテスト通知です",
		Data: map[string]interface{}{
			"test": "data",
		},
		UserID: 1,
	}

	notification := &models.Notification{
		ID:      1,
		Type:    req.Type,
		Title:   req.Title,
		Message: req.Message,
		Data:    req.Data,
		UserID:  req.UserID,
	}

	mockService.On("CreateNotification", mock.Anything, req).Return(notification, nil)

	body, _ := json.Marshal(req)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/notifications", bytes.NewBuffer(body))

	handler.CreateNotification(w, r)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response models.Notification
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, notification.ID, response.ID)
	assert.Equal(t, notification.Title, response.Title)

	mockService.AssertExpectations(t)
}

func TestGetNotification(t *testing.T) {
	mockService := new(mocks.MockNotificationService)
	handler := NewNotificationHandler(mockService)

	notification := &models.Notification{
		ID:      1,
		Type:    models.NotificationTypeDeliveryStatus,
		Title:   "テスト通知",
		Message: "これはテスト通知です",
		UserID:  1,
	}

	mockService.On("GetNotification", mock.Anything, int64(1)).Return(notification, nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/notifications/1", nil)
	r = mux.SetURLVars(r, map[string]string{"id": "1"})

	handler.GetNotification(w, r)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.Notification
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, notification.ID, response.ID)
	assert.Equal(t, notification.Title, response.Title)

	mockService.AssertExpectations(t)
}

func TestListNotifications(t *testing.T) {
	mockService := new(mocks.MockNotificationService)
	handler := NewNotificationHandler(mockService)

	notifications := []*models.Notification{
		{
			ID:      1,
			Type:    models.NotificationTypeDeliveryStatus,
			Title:   "テスト通知1",
			Message: "これはテスト通知1です",
			UserID:  1,
		},
		{
			ID:      2,
			Type:    models.NotificationTypeDeliveryComplete,
			Title:   "テスト通知2",
			Message: "これはテスト通知2です",
			UserID:  1,
		},
	}

	mockService.On("ListNotifications", mock.Anything, int64(1)).Return(notifications, nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/notifications/user/1", nil)
	r = mux.SetURLVars(r, map[string]string{"userID": "1"})

	handler.ListNotifications(w, r)

	assert.Equal(t, http.StatusOK, w.Code)

	var response []*models.Notification
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, len(notifications), len(response))
	assert.Equal(t, notifications[0].ID, response[0].ID)
	assert.Equal(t, notifications[1].ID, response[1].ID)

	mockService.AssertExpectations(t)
}

func TestMarkAsRead(t *testing.T) {
	mockService := new(mocks.MockNotificationService)
	handler := NewNotificationHandler(mockService)

	mockService.On("MarkAsRead", mock.Anything, int64(1)).Return(nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("PUT", "/notifications/1/read", nil)
	r = mux.SetURLVars(r, map[string]string{"id": "1"})

	handler.MarkAsRead(w, r)

	assert.Equal(t, http.StatusNoContent, w.Code)
	mockService.AssertExpectations(t)
}

func TestDeleteNotification(t *testing.T) {
	mockService := new(mocks.MockNotificationService)
	handler := NewNotificationHandler(mockService)

	mockService.On("DeleteNotification", mock.Anything, int64(1)).Return(nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("DELETE", "/notifications/1", nil)
	r = mux.SetURLVars(r, map[string]string{"id": "1"})

	handler.DeleteNotification(w, r)

	assert.Equal(t, http.StatusNoContent, w.Code)
	mockService.AssertExpectations(t)
}

package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"tea-logistics/pkg/models"

	"github.com/gorilla/mux"
)

/*
 * 通知ハンドラー
 * 通知関連のHTTPエンドポイントを実装する
 */

// NotificationService 通知サービスインターフェース
type NotificationService interface {
	CreateNotification(ctx context.Context, req *models.CreateNotificationRequest) (*models.Notification, error)
	GetNotification(ctx context.Context, id int64) (*models.Notification, error)
	ListNotifications(ctx context.Context, userID int64) ([]*models.Notification, error)
	MarkAsRead(ctx context.Context, id int64) error
	DeleteNotification(ctx context.Context, id int64) error
}

// NotificationHandler 通知ハンドラー
type NotificationHandler struct {
	service NotificationService
}

// NewNotificationHandler 通知ハンドラーを作成する
func NewNotificationHandler(service NotificationService) *NotificationHandler {
	return &NotificationHandler{
		service: service,
	}
}

// RegisterRoutes ルートを登録する
func (h *NotificationHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/notifications", h.CreateNotification).Methods("POST")
	router.HandleFunc("/notifications/{id}", h.GetNotification).Methods("GET")
	router.HandleFunc("/notifications/user/{userID}", h.ListNotifications).Methods("GET")
	router.HandleFunc("/notifications/{id}/read", h.MarkAsRead).Methods("PUT")
	router.HandleFunc("/notifications/{id}", h.DeleteNotification).Methods("DELETE")
}

// CreateNotification 通知を作成する
func (h *NotificationHandler) CreateNotification(w http.ResponseWriter, r *http.Request) {
	var req models.CreateNotificationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "リクエストの解析に失敗しました", http.StatusBadRequest)
		return
	}

	notification, err := h.service.CreateNotification(r.Context(), &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(notification)
}

// GetNotification 通知を取得する
func (h *NotificationHandler) GetNotification(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "無効なID", http.StatusBadRequest)
		return
	}

	notification, err := h.service.GetNotification(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notification)
}

// ListNotifications ユーザーの通知一覧を取得する
func (h *NotificationHandler) ListNotifications(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.ParseInt(vars["userID"], 10, 64)
	if err != nil {
		http.Error(w, "無効なユーザーID", http.StatusBadRequest)
		return
	}

	notifications, err := h.service.ListNotifications(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notifications)
}

// MarkAsRead 通知を既読にする
func (h *NotificationHandler) MarkAsRead(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "無効なID", http.StatusBadRequest)
		return
	}

	if err := h.service.MarkAsRead(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DeleteNotification 通知を削除する
func (h *NotificationHandler) DeleteNotification(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "無効なID", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteNotification(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

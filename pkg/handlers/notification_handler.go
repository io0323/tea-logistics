package handlers

import (
	"net/http"
	"strconv"

	"tea-logistics/pkg/models"
	"tea-logistics/pkg/services"

	"github.com/gin-gonic/gin"
)

/*
 * 通知ハンドラ
 * HTTPリクエストを処理し、通知サービスを呼び出す
 */

// NotificationHandler 通知ハンドラ
type NotificationHandler struct {
	service services.NotificationService
}

// NewNotificationHandler 通知ハンドラを作成する
func NewNotificationHandler(service services.NotificationService) *NotificationHandler {
	return &NotificationHandler{
		service: service,
	}
}

// CreateNotification 通知を作成する
func (h *NotificationHandler) CreateNotification(c *gin.Context) {
	var req models.CreateNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効なリクエストです"})
		return
	}

	notification, err := h.service.CreateNotification(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, notification)
}

// GetNotification 通知を取得する
func (h *NotificationHandler) GetNotification(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効な通知IDです"})
		return
	}

	notification, err := h.service.GetNotification(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, notification)
}

// ListNotifications 通知一覧を取得する
func (h *NotificationHandler) ListNotifications(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効なユーザーIDです"})
		return
	}

	notifications, err := h.service.ListNotifications(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, notifications)
}

// MarkAsRead 通知を既読にする
func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効な通知IDです"})
		return
	}

	if err := h.service.MarkAsRead(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// DeleteNotification 通知を削除する
func (h *NotificationHandler) DeleteNotification(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効な通知IDです"})
		return
	}

	if err := h.service.DeleteNotification(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

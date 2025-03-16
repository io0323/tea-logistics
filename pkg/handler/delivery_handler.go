package handler

import (
	"net/http"
	"strconv"

	"tea-logistics/pkg/models"
	"tea-logistics/pkg/services"

	"github.com/gin-gonic/gin"
)

/*
 * 配送ハンドラー
 * 配送関連のHTTP APIエンドポイントを実装する
 */

// DeliveryHandler 配送ハンドラー
type DeliveryHandler struct {
	service *services.DeliveryService
}

// NewDeliveryHandler 配送ハンドラーを作成する
func NewDeliveryHandler(service *services.DeliveryService) *DeliveryHandler {
	return &DeliveryHandler{
		service: service,
	}
}

// RegisterRoutes ルートを登録する
func (h *DeliveryHandler) RegisterRoutes(router *gin.Engine) {
	deliveries := router.Group("/api/deliveries")
	{
		deliveries.POST("", h.CreateDelivery)
		deliveries.GET("", h.ListDeliveries)
		deliveries.GET("/:id", h.GetDelivery)
		deliveries.PUT("/:id/status", h.UpdateDeliveryStatus)
		deliveries.POST("/:id/complete", h.CompleteDelivery)
		deliveries.POST("/:id/tracking", h.CreateDeliveryTracking)
		deliveries.GET("/:id/tracking", h.ListDeliveryTrackings)
	}
}

// CreateDelivery 配送を作成する
func (h *DeliveryHandler) CreateDelivery(c *gin.Context) {
	var req models.CreateDeliveryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効なリクエスト形式です"})
		return
	}

	delivery, err := h.service.CreateDelivery(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, delivery)
}

// ListDeliveries 配送一覧を取得する
func (h *DeliveryHandler) ListDeliveries(c *gin.Context) {
	deliveries, err := h.service.ListDeliveries(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, deliveries)
}

// GetDelivery 配送を取得する
func (h *DeliveryHandler) GetDelivery(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効なID形式です"})
		return
	}

	delivery, err := h.service.GetDelivery(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, delivery)
}

// UpdateDeliveryStatus 配送ステータスを更新する
func (h *DeliveryHandler) UpdateDeliveryStatus(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効なID形式です"})
		return
	}

	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効なリクエスト形式です"})
		return
	}

	if err := h.service.UpdateDeliveryStatus(c.Request.Context(), id, req.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

// CompleteDelivery 配送を完了する
func (h *DeliveryHandler) CompleteDelivery(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効なID形式です"})
		return
	}

	if err := h.service.CompleteDelivery(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

// CreateDeliveryTracking 配送追跡を作成する
func (h *DeliveryHandler) CreateDeliveryTracking(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効なID形式です"})
		return
	}

	var req models.CreateTrackingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効なリクエスト形式です"})
		return
	}

	req.DeliveryID = id
	tracking, err := h.service.CreateDeliveryTracking(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, tracking)
}

// ListDeliveryTrackings 配送追跡履歴を取得する
func (h *DeliveryHandler) ListDeliveryTrackings(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効なID形式です"})
		return
	}

	trackings, err := h.service.ListDeliveryTrackings(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, trackings)
}

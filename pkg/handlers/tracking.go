package handlers

import (
	"net/http"
	"strconv"

	"tea-logistics/pkg/models"
	"tea-logistics/pkg/services"

	"github.com/gin-gonic/gin"
)

/*
 * 配送追跡ハンドラ
 * 配送追跡関連のHTTPリクエストを処理する
 */

// TrackingHandler 配送追跡ハンドラ
type TrackingHandler struct {
	service *services.TrackingService
}

// NewTrackingHandler 配送追跡ハンドラを作成する
func NewTrackingHandler(service *services.TrackingService) *TrackingHandler {
	return &TrackingHandler{service: service}
}

// InitializeTracking 配送追跡を初期化する
func (h *TrackingHandler) InitializeTracking(c *gin.Context) {
	deliveryID, err := strconv.ParseInt(c.Param("delivery_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効な配送IDです"})
		return
	}

	fromLocation := c.Query("from_location")
	if fromLocation == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "出発地点を指定してください"})
		return
	}

	tracking, err := h.service.InitializeTracking(c.Request.Context(), deliveryID, fromLocation)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, tracking)
}

// UpdateTrackingStatus 配送追跡ステータスを更新する
func (h *TrackingHandler) UpdateTrackingStatus(c *gin.Context) {
	type UpdateRequest struct {
		Status      models.TrackingStatus `json:"status" binding:"required"`
		Location    string                `json:"location" binding:"required"`
		Description string                `json:"description" binding:"required"`
	}

	trackingID := c.Param("tracking_id")
	if trackingID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "追跡IDを指定してください"})
		return
	}

	var req UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効なリクエスト形式です"})
		return
	}

	err := h.service.UpdateTrackingStatus(
		c.Request.Context(),
		trackingID,
		req.Status,
		req.Location,
		req.Description,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "追跡ステータスを更新しました"})
}

// AddTrackingEvent 配送追跡イベントを追加する
func (h *TrackingHandler) AddTrackingEvent(c *gin.Context) {
	trackingID := c.Param("tracking_id")
	if trackingID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "追跡IDを指定してください"})
		return
	}

	var event models.TrackingEvent
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効なリクエスト形式です"})
		return
	}

	event.TrackingID = trackingID
	if err := h.service.AddTrackingEvent(c.Request.Context(), &event); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, event)
}

// GetTrackingInfo 配送追跡情報を取得する
func (h *TrackingHandler) GetTrackingInfo(c *gin.Context) {
	trackingID := c.Param("tracking_id")
	if trackingID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "追跡IDを指定してください"})
		return
	}

	tracking, err := h.service.GetTrackingInfo(c.Request.Context(), trackingID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tracking)
}

// SetTrackingCondition 追跡条件を設定する
func (h *TrackingHandler) SetTrackingCondition(c *gin.Context) {
	trackingID := c.Param("tracking_id")
	if trackingID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "追跡IDを指定してください"})
		return
	}

	var condition models.TrackingCondition
	if err := c.ShouldBindJSON(&condition); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効なリクエスト形式です"})
		return
	}

	condition.TrackingID = trackingID
	if err := h.service.SetTrackingCondition(c.Request.Context(), &condition); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, condition)
}

// GetTrackingCondition 追跡条件を取得する
func (h *TrackingHandler) GetTrackingCondition(c *gin.Context) {
	trackingID := c.Param("tracking_id")
	if trackingID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "追跡IDを指定してください"})
		return
	}

	condition, err := h.service.GetTrackingCondition(c.Request.Context(), trackingID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, condition)
}

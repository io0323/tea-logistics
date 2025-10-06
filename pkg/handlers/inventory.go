package handlers

import (
	"net/http"
	"strconv"
	"time"

	"tea-logistics/pkg/models"
	"tea-logistics/pkg/services"

	"github.com/gin-gonic/gin"
)

/*
 * 在庫管理ハンドラ
 * 在庫関連のHTTPリクエストを処理する
 */

// InventoryHandler 在庫管理ハンドラ
type InventoryHandler struct {
	service *services.InventoryService
}

// NewInventoryHandler 在庫管理ハンドラを作成する
func NewInventoryHandler(service *services.InventoryService) *InventoryHandler {
	return &InventoryHandler{service: service}
}

// GetInventory 在庫情報を取得する
func (h *InventoryHandler) GetInventory(c *gin.Context) {
	productID, err := strconv.ParseInt(c.Param("product_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効な商品IDです"})
		return
	}

	location := c.Query("location")
	if location == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ロケーションを指定してください"})
		return
	}

	inventory, err := h.service.GetProductInventory(c.Request.Context(), productID, location)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, inventory)
}

// UpdateInventory 在庫を更新する
func (h *InventoryHandler) UpdateInventory(c *gin.Context) {
	type UpdateRequest struct {
		ProductID int64  `json:"product_id" binding:"required"`
		Location  string `json:"location" binding:"required"`
		Quantity  int    `json:"quantity" binding:"required"`
	}

	var req UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効なリクエスト形式です"})
		return
	}

	err := h.service.UpdateInventoryQuantity(
		c.Request.Context(),
		req.ProductID,
		req.Location,
		req.Quantity,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "在庫を更新しました"})
}

// TransferInventory 在庫を移動する
func (h *InventoryHandler) TransferInventory(c *gin.Context) {
	type TransferRequest struct {
		ProductID    int64  `json:"product_id" binding:"required"`
		FromLocation string `json:"from_location" binding:"required"`
		ToLocation   string `json:"to_location" binding:"required"`
		Quantity     int    `json:"quantity" binding:"required"`
	}

	var req TransferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効なリクエスト形式です"})
		return
	}

	err := h.service.TransferInventory(
		c.Request.Context(),
		req.ProductID,
		req.FromLocation,
		req.ToLocation,
		req.Quantity,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "在庫を移動しました"})
}

// CheckAvailability 在庫の利用可能性をチェックする
func (h *InventoryHandler) CheckAvailability(c *gin.Context) {
	productID, err := strconv.ParseInt(c.Param("product_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効な商品IDです"})
		return
	}

	location := c.Query("location")
	if location == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ロケーションを指定してください"})
		return
	}

	quantity, err := strconv.Atoi(c.Query("quantity"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効な数量です"})
		return
	}

	available, err := h.service.CheckAvailability(
		c.Request.Context(),
		productID,
		location,
		quantity,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	message := "在庫が不足しています"
	if available {
		message = "在庫は利用可能です"
	}

	c.JSON(http.StatusOK, gin.H{
		"available": available,
		"message":   message,
	})
}

// CreateMovement 在庫移動を作成する
func (h *InventoryHandler) CreateMovement(c *gin.Context) {
	type MovementRequest struct {
		ProductID       int64  `json:"product_id" binding:"required"`
		FromLocation    string `json:"from_location" binding:"required"`
		ToLocation      string `json:"to_location" binding:"required"`
		Quantity        int    `json:"quantity" binding:"required"`
		MovementType    string `json:"movement_type" binding:"required"`
		MovementDate    string `json:"movement_date"`
		ReferenceNumber string `json:"reference_number"`
	}

  var req MovementRequest
  if err := c.ShouldBindJSON(&req); err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": "無効なリクエスト形式です"})
    return
  }

  // MovementTypeのバリデーション
  validMovementTypes := map[string]bool{
    string(models.MovementTypeInbound):     true,
    string(models.MovementTypeOutbound):    true,
    string(models.MovementTypeTransfer):    true,
    string(models.MovementTypeAdjustment):  true,
  }
  if !validMovementTypes[req.MovementType] {
    c.JSON(http.StatusBadRequest, gin.H{"error": "無効な移動タイプです。有効な値: inbound, outbound, transfer, adjustment"})
    return
  }

  var movementTime time.Time
  if req.MovementDate != "" {
    // RFC3339想定。異なる形式が必要ならばここで調整
    if t, err := time.Parse(time.RFC3339, req.MovementDate); err == nil {
      movementTime = t
    } else {
      c.JSON(http.StatusBadRequest, gin.H{"error": "movement_dateはRFC3339形式で指定してください"})
      return
    }
  } else {
    movementTime = time.Now()
  }

  movement, err := h.service.CreateMovement(c.Request.Context(), &models.CreateMovementRequest{
    ProductID:       req.ProductID,
    FromLocation:    req.FromLocation,
    ToLocation:      req.ToLocation,
    Quantity:        req.Quantity,
    MovementType:    models.MovementType(req.MovementType),
    MovementDate:    movementTime,
    ReferenceNumber: req.ReferenceNumber,
  })
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, movement)
}

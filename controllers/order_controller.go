package controllers

import (
	"net/http"

	"backend/database"
	"backend/models"

	"github.com/gin-gonic/gin"
)

type OrderController struct{}

func NewOrderController() *OrderController {
	return &OrderController{}
}

type OrderRequest struct {
	LabdipID  *uint   `json:"labdip_id"`
	LabdipNo  string  `json:"labdip_no" binding:"required"`
	PartyName string  `json:"party_name" binding:"required"`
	Quantity  float64 `json:"quantity" binding:"required"`
	Rate      float64 `json:"rate" binding:"required"`
	Status    string  `json:"status"`
	Remarks   string  `json:"remarks"`
}

type OrderStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

func (oc *OrderController) GetOrders(c *gin.Context) {
	var orders []models.Order
	if err := database.DB.Order("created_at desc").Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to fetch orders",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    orders,
	})
}

func (oc *OrderController) GetOrderByID(c *gin.Context) {
	id := c.Param("id")
	var order models.Order
	if err := database.DB.First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Order not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    order,
	})
}

func (oc *OrderController) CreateOrder(c *gin.Context) {
	var req OrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Required fields missing: Labdip No, Party Name, Quantity, and Rate are required",
		})
		return
	}

	totalAmount := req.Quantity * req.Rate
	status := req.Status
	if status == "" {
		status = "Pending"
	}

	order := models.Order{
		LabdipID:    req.LabdipID,
		LabdipNo:    req.LabdipNo,
		PartyName:   req.PartyName,
		Quantity:    req.Quantity,
		Rate:        req.Rate,
		TotalAmount: totalAmount,
		Status:      status,
		Remarks:     req.Remarks,
	}

	if err := database.DB.Create(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to create order",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Order created successfully",
		"data":    order,
	})
}

func (oc *OrderController) UpdateOrder(c *gin.Context) {
	id := c.Param("id")
	var order models.Order
	if err := database.DB.First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Order not found",
		})
		return
	}

	var req OrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Required fields missing",
		})
		return
	}

	order.LabdipID = req.LabdipID
	order.LabdipNo = req.LabdipNo
	order.PartyName = req.PartyName
	order.Quantity = req.Quantity
	order.Rate = req.Rate
	order.TotalAmount = req.Quantity * req.Rate
	if req.Status != "" {
		order.Status = req.Status
	}
	order.Remarks = req.Remarks

	if err := database.DB.Save(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to update order",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Order updated successfully",
		"data":    order,
	})
}

func (oc *OrderController) UpdateOrderStatus(c *gin.Context) {
	id := c.Param("id")
	var order models.Order
	if err := database.DB.First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Order not found",
		})
		return
	}

	var req OrderStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.Status == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Valid status is required",
		})
		return
	}

	validStatuses := map[string]bool{
		"Pending":    true,
		"Processing": true,
		"Dispatched": true,
		"Delivered":  true,
		"Cancelled":  true,
	}

	if !validStatuses[req.Status] {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid status value",
		})
		return
	}

	order.Status = req.Status
	if err := database.DB.Save(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to update order status",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Order status updated successfully",
		"data":    order,
	})
}

func (oc *OrderController) DeleteOrder(c *gin.Context) {
	id := c.Param("id")
	var order models.Order
	if err := database.DB.First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Order not found",
		})
		return
	}

	if err := database.DB.Delete(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to delete order",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Order deleted successfully",
	})
}

package controllers

import (
	"net/http"
	"time"

	"backend/database"
	"backend/models"

	"github.com/gin-gonic/gin"
)

type DashboardController struct{}

func NewDashboardController() *DashboardController {
	return &DashboardController{}
}

type StatusCount struct {
	Status string `json:"status"`
	Count  int64  `json:"count"`
}

type MonthlyTrend struct {
	Month   string  `json:"month"`
	Orders  int64   `json:"orders"`
	Revenue float64 `json:"revenue"`
}

func (dc *DashboardController) GetDashboardStats(c *gin.Context) {
	var totalParties int64
	var totalQualities int64
	var totalLabdips int64
	var totalOrders int64
	var totalRevenue float64
	var totalQuantity float64

	// Counts
	database.DB.Model(&models.Party{}).Count(&totalParties)
	database.DB.Model(&models.Quality{}).Count(&totalQualities)
	database.DB.Model(&models.Labdip{}).Count(&totalLabdips)
	database.DB.Model(&models.Order{}).Count(&totalOrders)

	// Revenue & Quantity sums
	type OrderSums struct {
		TotalRevenue  float64
		TotalQuantity float64
	}
	var sums OrderSums
	database.DB.Model(&models.Order{}).Select("COALESCE(SUM(total_amount), 0) as total_revenue, COALESCE(SUM(quantity), 0) as total_quantity").Scan(&sums)
	totalRevenue = sums.TotalRevenue
	totalQuantity = sums.TotalQuantity

	// Labdips by status
	labdipsByStatus := map[string]int64{
		"Pending":     0,
		"In Progress": 0,
		"Approved":    0,
		"Rejected":    0,
	}
	var labdipStatusCounts []StatusCount
	database.DB.Model(&models.Labdip{}).Select("status, count(*) as count").Group("status").Scan(&labdipStatusCounts)
	for _, sc := range labdipStatusCounts {
		labdipsByStatus[sc.Status] = sc.Count
	}

	// Orders by status
	ordersByStatus := map[string]int64{
		"Pending":    0,
		"Processing": 0,
		"Dispatched": 0,
		"Delivered":  0,
		"Cancelled":  0,
	}
	var orderStatusCounts []StatusCount
	database.DB.Model(&models.Order{}).Select("status, count(*) as count").Group("status").Scan(&orderStatusCounts)
	for _, sc := range orderStatusCounts {
		ordersByStatus[sc.Status] = sc.Count
	}

	// Recent Labdips (latest 5)
	var recentLabdips []models.Labdip
	database.DB.Order("created_at desc").Limit(5).Find(&recentLabdips)

	// Recent Orders (latest 5)
	var recentOrders []models.Order
	database.DB.Order("created_at desc").Limit(5).Find(&recentOrders)

	// Monthly order trends for past 6 months
	monthlyTrends := make([]MonthlyTrend, 6)
	now := time.Now()
	for i := 5; i >= 0; i-- {
		monthDate := now.AddDate(0, -i, 0)
		monthKey := monthDate.Format("Jan 2006")
		startOfMonth := time.Date(monthDate.Year(), monthDate.Month(), 1, 0, 0, 0, 0, time.UTC)
		endOfMonth := startOfMonth.AddDate(0, 1, 0)

		var monthOrders int64
		var monthRevenue float64
		database.DB.Model(&models.Order{}).
			Where("created_at >= ? AND created_at < ?", startOfMonth, endOfMonth).
			Count(&monthOrders)

		type RevResult struct {
			Rev float64
		}
		var revRes RevResult
		database.DB.Model(&models.Order{}).
			Select("COALESCE(SUM(total_amount), 0) as rev").
			Where("created_at >= ? AND created_at < ?", startOfMonth, endOfMonth).
			Scan(&revRes)
		monthRevenue = revRes.Rev

		idx := 5 - i
		monthlyTrends[idx] = MonthlyTrend{
			Month:   monthKey,
			Orders:  monthOrders,
			Revenue: monthRevenue,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"total_parties":        totalParties,
			"total_qualities":      totalQualities,
			"total_labdips":        totalLabdips,
			"total_orders":         totalOrders,
			"total_revenue":        totalRevenue,
			"total_quantity":       totalQuantity,
			"labdips_by_status":    labdipsByStatus,
			"orders_by_status":     ordersByStatus,
			"monthly_order_trends": monthlyTrends,
			"recent_labdips":       recentLabdips,
			"recent_orders":        recentOrders,
		},
	})
}

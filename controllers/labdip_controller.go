package controllers

import (
	"net/http"
	"time"

	"backend/database"
	"backend/models"

	"github.com/gin-gonic/gin"
)

type LabdipController struct{}

func NewLabdipController() *LabdipController {
	return &LabdipController{}
}

type LabdipRequest struct {
	LabdipNo     string  `json:"labdip_no" binding:"required"`
	PartyName    string  `json:"party_name" binding:"required"`
	Status       string  `json:"status" binding:"required"`
	QualityID    *uint   `json:"quality_id"`
	QualityName  string  `json:"quality_name" binding:"required"`
	ColorName    string  `json:"color_name" binding:"required"`
	ReceivedDate *string `json:"received_date"`
	SendingDate  *string `json:"sending_date"`
	Remarks      string  `json:"remarks"`
}

func parseOptionalDate(dateStr *string) *time.Time {
	if dateStr == nil || *dateStr == "" {
		return nil
	}
	// Try parsing YYYY-MM-DD format
	t, err := time.Parse("2006-01-02", *dateStr)
	if err == nil {
		return &t
	}
	// Try RFC3339
	t, err = time.Parse(time.RFC3339, *dateStr)
	if err == nil {
		return &t
	}
	return nil
}

func (lc *LabdipController) GetLabdips(c *gin.Context) {
	var labdips []models.Labdip
	if err := database.DB.Order("created_at desc").Find(&labdips).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to fetch labdips",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    labdips,
	})
}

func (lc *LabdipController) GetLabdipByID(c *gin.Context) {
	id := c.Param("id")
	var labdip models.Labdip
	if err := database.DB.First(&labdip, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Labdip not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    labdip,
	})
}

func (lc *LabdipController) CreateLabdip(c *gin.Context) {
	var req LabdipRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Required fields are missing: Labdip No, Party Name, Status, Quality Name, and Color Name are required",
		})
		return
	}

	labdip := models.Labdip{
		LabdipNo:     req.LabdipNo,
		PartyName:    req.PartyName,
		Status:       req.Status,
		QualityID:    req.QualityID,
		QualityName:  req.QualityName,
		ColorName:    req.ColorName,
		ReceivedDate: parseOptionalDate(req.ReceivedDate),
		SendingDate:  parseOptionalDate(req.SendingDate),
		Remarks:      req.Remarks,
	}

	if err := database.DB.Create(&labdip).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to create labdip entry",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Labdip created successfully",
		"data":    labdip,
	})
}

func (lc *LabdipController) UpdateLabdip(c *gin.Context) {
	id := c.Param("id")
	var labdip models.Labdip
	if err := database.DB.First(&labdip, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Labdip not found",
		})
		return
	}

	var req LabdipRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Required fields are missing",
		})
		return
	}

	labdip.LabdipNo = req.LabdipNo
	labdip.PartyName = req.PartyName
	labdip.Status = req.Status
	labdip.QualityID = req.QualityID
	labdip.QualityName = req.QualityName
	labdip.ColorName = req.ColorName
	labdip.ReceivedDate = parseOptionalDate(req.ReceivedDate)
	labdip.SendingDate = parseOptionalDate(req.SendingDate)
	labdip.Remarks = req.Remarks

	if err := database.DB.Save(&labdip).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to update labdip entry",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Labdip updated successfully",
		"data":    labdip,
	})
}

func (lc *LabdipController) DeleteLabdip(c *gin.Context) {
	id := c.Param("id")
	var labdip models.Labdip
	if err := database.DB.First(&labdip, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Labdip not found",
		})
		return
	}

	if err := database.DB.Delete(&labdip).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to delete labdip entry",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Labdip deleted successfully",
	})
}

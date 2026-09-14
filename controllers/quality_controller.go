package controllers

import (
	"net/http"

	"backend/database"
	"backend/models"

	"github.com/gin-gonic/gin"
)

type QualityController struct{}

func NewQualityController() *QualityController {
	return &QualityController{}
}

type QualityRequest struct {
	QualityName string `json:"quality_name" binding:"required"`
	Code        string `json:"code" binding:"required"`
}

func (qc *QualityController) GetQualities(c *gin.Context) {
	var qualities []models.Quality
	if err := database.DB.Order("created_at desc").Find(&qualities).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to fetch qualities",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    qualities,
	})
}

func (qc *QualityController) GetQualityByID(c *gin.Context) {
	id := c.Param("id")
	var quality models.Quality
	if err := database.DB.First(&quality, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Quality not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    quality,
	})
}

func (qc *QualityController) CreateQuality(c *gin.Context) {
	var req QualityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Quality Name and Code are required",
		})
		return
	}

	quality := models.Quality{
		QualityName: req.QualityName,
		Code:        req.Code,
	}

	if err := database.DB.Create(&quality).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to create quality",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Quality created successfully",
		"data":    quality,
	})
}

func (qc *QualityController) UpdateQuality(c *gin.Context) {
	id := c.Param("id")
	var quality models.Quality
	if err := database.DB.First(&quality, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Quality not found",
		})
		return
	}

	var req QualityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Quality Name and Code are required",
		})
		return
	}

	quality.QualityName = req.QualityName
	quality.Code = req.Code

	if err := database.DB.Save(&quality).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to update quality",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Quality updated successfully",
		"data":    quality,
	})
}

func (qc *QualityController) DeleteQuality(c *gin.Context) {
	id := c.Param("id")
	var quality models.Quality
	if err := database.DB.First(&quality, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Quality not found",
		})
		return
	}

	if err := database.DB.Delete(&quality).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to delete quality",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Quality deleted successfully",
	})
}

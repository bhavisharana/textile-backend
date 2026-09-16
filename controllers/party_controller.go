package controllers

import (
	"net/http"

	"backend/database"
	"backend/models"

	"github.com/gin-gonic/gin"
)

type PartyController struct{}

func NewPartyController() *PartyController {
	return &PartyController{}
}

type PartyRequest struct {
	Name string `json:"name" binding:"required"`
	Code string `json:"code" binding:"required"`
}

func (pc *PartyController) GetParties(c *gin.Context) {
	var parties []models.Party
	if err := database.DB.Order("created_at desc").Find(&parties).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to fetch parties",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    parties,
	})
}

func (pc *PartyController) GetPartyByID(c *gin.Context) {
	id := c.Param("id")
	var party models.Party
	if err := database.DB.First(&party, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Party not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    party,
	})
}

func (pc *PartyController) CreateParty(c *gin.Context) {
	var req PartyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Party Name and Code are required",
		})
		return
	}

	party := models.Party{
		Name: req.Name,
		Code: req.Code,
	}

	if err := database.DB.Create(&party).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to create party. Code must be unique.",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Party created successfully",
		"data":    party,
	})
}

func (pc *PartyController) UpdateParty(c *gin.Context) {
	id := c.Param("id")
	var party models.Party
	if err := database.DB.First(&party, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Party not found",
		})
		return
	}

	var req PartyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Party Name and Code are required",
		})
		return
	}

	party.Name = req.Name
	party.Code = req.Code

	if err := database.DB.Save(&party).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to update party",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Party updated successfully",
		"data":    party,
	})
}

func (pc *PartyController) DeleteParty(c *gin.Context) {
	id := c.Param("id")
	var party models.Party
	if err := database.DB.First(&party, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Party not found",
		})
		return
	}

	if err := database.DB.Delete(&party).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to delete party",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Party deleted successfully",
	})
}

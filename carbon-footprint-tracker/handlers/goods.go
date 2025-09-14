package handlers

import (
	"carbon-footprint-tracker/models"
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func AddGoodsPurchased(c *gin.Context) {
	var req models.GoodsPurchased
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Location == "" {
		req.Location = "Overall"
	}

	if req.Date.IsZero() {
		req.Date = time.Now()
	}

	if err := req.Create(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add goods purchased", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Goods purchased added successfully", "id": req.ID})
}

func GetGoodsPurchased(c *gin.Context) {
	goods, err := models.GetAllGoodsPurchased()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve goods purchased", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, goods)
}

func UpdateGoodsPurchased(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	g, err := models.GetGoodsPurchasedByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Goods purchased entry not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve goods purchased", "details": err.Error()})
		return
	}

	var req models.GoodsPurchased
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.ItemName != "" {
		g.ItemName = req.ItemName
	}
	if req.Quantity != 0 {
		g.Quantity = req.Quantity
	}
	if req.Cost != 0 {
		g.Cost = req.Cost
	}
	if !req.Date.IsZero() {
		g.Date = req.Date
	}
	if req.Location != "" {
		g.Location = req.Location
	}

	if err := g.Update(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update goods purchased", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Goods purchased updated successfully"})
}

func DeleteGoodsPurchased(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := models.DeleteGoodsPurchased(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete goods purchased", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Goods purchased deleted successfully"})
}

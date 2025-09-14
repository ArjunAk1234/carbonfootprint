package handlers

import (
	"carbon-footprint-tracker/models"
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func AddElectricConsumption(c *gin.Context) {
	var req models.ElectricConsumption
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add electric consumption", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Electric consumption added successfully", "id": req.ID})
}

func GetElectricConsumptions(c *gin.Context) {
	consumptions, err := models.GetAllElectricConsumptions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve electric consumptions", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, consumptions)
}

func UpdateElectricConsumption(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	e, err := models.GetElectricConsumptionByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Electric consumption entry not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve electric consumption", "details": err.Error()})
		return
	}

	var req models.ElectricConsumption
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Source != "" {
		e.Source = req.Source
	}
	if req.KWH != 0 {
		e.KWH = req.KWH
	}
	if req.FuelLiters != 0 {
		e.FuelLiters = req.FuelLiters
	}
	if req.Hours != 0 {
		e.Hours = req.Hours
	}
	if !req.Date.IsZero() {
		e.Date = req.Date
	}
	if req.Location != "" {
		e.Location = req.Location
	}

	if err := e.Update(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update electric consumption", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Electric consumption updated successfully"})
}

func DeleteElectricConsumption(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := models.DeleteElectricConsumption(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete electric consumption", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Electric consumption deleted successfully"})
}

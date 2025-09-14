package handlers

import (
	"carbon-footprint-tracker/models"
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func AddWaterReading(c *gin.Context) {
	var req models.WaterConsumption
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add water reading", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Water reading added successfully", "id": req.ID})
}

func GetWaterReadings(c *gin.Context) {
	readings, err := models.GetAllWaterConsumptions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve water readings", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, readings)
}

func UpdateWaterReading(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	w, err := models.GetWaterConsumptionByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Water reading not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve water reading", "details": err.Error()})
		return
	}

	var req models.WaterConsumption
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.MeterReading != 0 {
		w.MeterReading = req.MeterReading
	}
	if !req.Date.IsZero() {
		w.Date = req.Date
	}
	if req.Location != "" {
		w.Location = req.Location
	}

	if err := w.Update(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update water reading", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Water reading updated successfully"})
}

func DeleteWaterReading(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := models.DeleteWaterConsumption(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete water reading", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Water reading deleted successfully"})
}

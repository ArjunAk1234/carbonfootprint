package handlers

import (
	"carbon-footprint-tracker/models"
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func AddTransportData(c *gin.Context) {
	var req models.Transport
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add transport data", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Transport data added successfully", "id": req.ID})
}

func GetTransportData(c *gin.Context) {
	transports, err := models.GetAllTransports()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve transport data", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, transports)
}

func UpdateTransportData(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	t, err := models.GetTransportByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Transport data not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve transport data", "details": err.Error()})
		return
	}

	var req models.Transport
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.VehicleType != "" {
		t.VehicleType = req.VehicleType
	}
	if req.FuelType != "" {
		t.FuelType = req.FuelType
	}
	if req.DistanceKM != 0 {
		t.DistanceKM = req.DistanceKM
	}
	if req.FuelLiters != 0 {
		t.FuelLiters = req.FuelLiters
	}
	if !req.Date.IsZero() {
		t.Date = req.Date
	}
	if req.Location != "" {
		t.Location = req.Location
	}

	if err := t.Update(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update transport data", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transport data updated successfully"})
}

func DeleteTransportData(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := models.DeleteTransport(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete transport data", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transport data deleted successfully"})
}

package handlers

import (
	"carbon-footprint-tracker/models"
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func AddWasteEntry(c *gin.Context) {
	var req models.Waste
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add waste entry", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Waste entry added successfully", "id": req.ID})
}
func GetWasteData(c *gin.Context) {
	wastes, err := models.GetAllWasteEntries()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve waste data", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, wastes)
}

func UpdateWasteEntry(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	w, err := models.GetWasteByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Waste entry not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve waste entry", "details": err.Error()})
		return
	}

	var req models.Waste
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.SpotName != "" {
		w.SpotName = req.SpotName
	}
	if req.WasteType != "" {
		w.WasteType = req.WasteType
	}
	if req.WeightKG != 0 {
		w.WeightKG = req.WeightKG
	}
	if !req.Date.IsZero() {
		w.Date = req.Date
	}
	if req.Location != "" {
		w.Location = req.Location
	}

	if err := w.Update(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update waste entry", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Waste entry updated successfully"})
}

func DeleteWasteEntry(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := models.DeleteWaste(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete waste entry", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Waste entry deleted successfully"})
}

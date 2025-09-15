package handlers

import (
	"carbon-footprint-tracker/models"
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type WasteRequest struct {
	Date               time.Time `json:"date"`
	CollectionLocation string    `json:"collection_location"`
	WasteType          string    `json:"waste_type" binding:"required"`
	SubCategory        string    `json:"sub_category"`
	WeightKG           float64   `json:"weight_kg" binding:"required"`
	CollectionMethod   string    `json:"collection_method"`
	TransportMode      string    `json:"transport_mode"`
	Destination        string    `json:"destination"`
	Remarks            string    `json:"remarks"`
}

func AddWasteEntry(c *gin.Context) {
	var req WasteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.CollectionLocation == "" {
		req.CollectionLocation = "Overall" // Default location
	}
	if req.Date.IsZero() {
		req.Date = time.Now()
	}

	w := models.Waste{
		Date:               req.Date,
		CollectionLocation: req.CollectionLocation,
		WasteType:          req.WasteType,
		SubCategory:        toNullString(req.SubCategory),
		WeightKG:           req.WeightKG,
		CollectionMethod:   toNullString(req.CollectionMethod),
		TransportMode:      toNullString(req.TransportMode),
		Destination:        toNullString(req.Destination),
		Remarks:            toNullString(req.Remarks),
	}

	if err := w.Create(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add waste entry", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Waste entry added successfully", "id": w.ID})
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

	var req WasteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	w.Date = req.Date
	if req.CollectionLocation != "" {
		w.CollectionLocation = req.CollectionLocation
	}
	if req.WasteType != "" {
		w.WasteType = req.WasteType
	}
	w.SubCategory = toNullString(req.SubCategory)
	if req.WeightKG != 0 {
		w.WeightKG = req.WeightKG
	}
	w.CollectionMethod = toNullString(req.CollectionMethod)
	w.TransportMode = toNullString(req.TransportMode)
	w.Destination = toNullString(req.Destination)
	w.Remarks = toNullString(req.Remarks)

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

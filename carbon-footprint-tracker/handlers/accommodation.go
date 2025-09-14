package handlers

import (
	"carbon-footprint-tracker/models"
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func AddAccommodationData(c *gin.Context) {
	var req models.Accommodation
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add accommodation data", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Accommodation data added successfully", "id": req.ID})
}

func GetAccommodationData(c *gin.Context) {
	accommodations, err := models.GetAllAccommodations()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve accommodation data", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, accommodations)
}

func UpdateAccommodationData(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	a, err := models.GetAccommodationByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Accommodation data not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve accommodation data", "details": err.Error()})
		return
	}

	var req models.Accommodation
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.PeopleCount != 0 {
		a.PeopleCount = req.PeopleCount
	}
	if req.Nights != 0 {
		a.Nights = req.Nights
	}
	if !req.Date.IsZero() {
		a.Date = req.Date
	}
	if req.Location != "" {
		a.Location = req.Location
	}

	if err := a.Update(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update accommodation data", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Accommodation data updated successfully"})
}

func DeleteAccommodationData(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := models.DeleteAccommodation(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete accommodation data", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Accommodation data deleted successfully"})
}

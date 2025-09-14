package handlers

import (
	"carbon-footprint-tracker/models"
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func AddPopulation(c *gin.Context) {
	var req models.Population
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add population stats", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Population stats added successfully", "id": req.ID})
}

func GetPopulationStats(c *gin.Context) {
	populations, err := models.GetAllPopulations()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve population stats", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, populations)
}

func UpdatePopulation(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	p, err := models.GetPopulationByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Population stats not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve population stats", "details": err.Error()})
		return
	}

	var req models.Population
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.RegisteredCount != 0 {
		p.RegisteredCount = req.RegisteredCount
	}
	if req.FloatingCount != 0 {
		p.FloatingCount = req.FloatingCount
	}
	if !req.Date.IsZero() {
		p.Date = req.Date
	}
	if req.Location != "" {
		p.Location = req.Location
	}

	if err := p.Update(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update population stats", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Population stats updated successfully"})
}

func DeletePopulation(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := models.DeletePopulation(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete population stats", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Population stats deleted successfully"})
}

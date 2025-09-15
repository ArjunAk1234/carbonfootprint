package handlers

import (
	"carbon-footprint-tracker/models"
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type WaterConsumptionRequest struct {
	Date                    time.Time `json:"date"`
	Location                string    `json:"location"`
	WaterSource             string    `json:"water_source"`
	CumulativeMeterReading  float64   `json:"cumulative_meter_reading"`
	TotalConsumptionKLD     float64   `json:"total_consumption_kld" binding:"required"`
	PerCapitaConsumptionLPD float64   `json:"per_capita_consumption_lpd"`
	UsageType               string    `json:"usage_type"`
	Remarks                 string    `json:"remarks"`
}

func AddWaterConsumption(c *gin.Context) {
	var req WaterConsumptionRequest
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

	w := models.WaterConsumption{
		Date:                    req.Date,
		Location:                req.Location,
		WaterSource:             toNullString(req.WaterSource),
		CumulativeMeterReading:  toNullFloat64(req.CumulativeMeterReading),
		TotalConsumptionKLD:     req.TotalConsumptionKLD,
		PerCapitaConsumptionLPD: toNullFloat64(req.PerCapitaConsumptionLPD),
		UsageType:               toNullString(req.UsageType),
		Remarks:                 toNullString(req.Remarks),
	}

	if err := w.Create(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add water consumption", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Water consumption added successfully", "id": w.ID})
}

func GetWaterConsumptions(c *gin.Context) {
	readings, err := models.GetAllWaterConsumptions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve water consumptions", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, readings)
}

func UpdateWaterConsumption(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	w, err := models.GetWaterConsumptionByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Water consumption reading not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve water consumption", "details": err.Error()})
		return
	}

	var req WaterConsumptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	w.Date = req.Date
	if req.Location != "" {
		w.Location = req.Location
	}
	w.WaterSource = toNullString(req.WaterSource)
	w.CumulativeMeterReading = toNullFloat64(req.CumulativeMeterReading)
	if req.TotalConsumptionKLD != 0 {
		w.TotalConsumptionKLD = req.TotalConsumptionKLD
	}
	w.PerCapitaConsumptionLPD = toNullFloat64(req.PerCapitaConsumptionLPD)
	w.UsageType = toNullString(req.UsageType)
	w.Remarks = toNullString(req.Remarks)

	if err := w.Update(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update water consumption", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Water consumption updated successfully"})
}

func DeleteWaterConsumption(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := models.DeleteWaterConsumption(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete water consumption", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Water consumption deleted successfully"})
}

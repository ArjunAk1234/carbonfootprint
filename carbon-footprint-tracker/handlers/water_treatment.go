package handlers

import (
	"carbon-footprint-tracker/models"
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type WaterTreatmentRequest struct {
	Date                        time.Time `json:"date"`
	Location                    string    `json:"location" binding:"required"`
	TreatedLitersPerDay         float64   `json:"treated_liters_per_day"`
	UltraFiltrationLitersPerDay float64   `json:"ultra_filtration_liters_per_day"`
	PercentageWaterReused       float64   `json:"percentage_water_reused"`
	ElectricityUsedKWH          float64   `json:"electricity_used_kwh"`
	ChemicalsUsedDescription    string    `json:"chemicals_used_description"`
	ChemicalsUsedQuantityKG     float64   `json:"chemicals_used_quantity_kg"`
	Remarks                     string    `json:"remarks"`
}

func AddWaterTreatment(c *gin.Context) {
	var req WaterTreatmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Date.IsZero() {
		req.Date = time.Now()
	}

	wt := models.WaterTreatment{
		Date:                        req.Date,
		Location:                    req.Location,
		TreatedLitersPerDay:         toNullFloat64(req.TreatedLitersPerDay),
		UltraFiltrationLitersPerDay: toNullFloat64(req.UltraFiltrationLitersPerDay),
		PercentageWaterReused:       toNullFloat64(req.PercentageWaterReused),
		ElectricityUsedKWH:          toNullFloat64(req.ElectricityUsedKWH),
		ChemicalsUsedDescription:    toNullString(req.ChemicalsUsedDescription),
		ChemicalsUsedQuantityKG:     toNullFloat64(req.ChemicalsUsedQuantityKG),
		Remarks:                     toNullString(req.Remarks),
	}

	if err := wt.Create(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add water treatment", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Water treatment added successfully", "id": wt.ID})
}

func GetWaterTreatments(c *gin.Context) {
	treatments, err := models.GetAllWaterTreatments()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve water treatments", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, treatments)
}

func UpdateWaterTreatment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	wt, err := models.GetWaterTreatmentByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Water treatment entry not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve water treatment", "details": err.Error()})
		return
	}

	var req WaterTreatmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	wt.Date = req.Date
	if req.Location != "" {
		wt.Location = req.Location
	}
	wt.TreatedLitersPerDay = toNullFloat64(req.TreatedLitersPerDay)
	wt.UltraFiltrationLitersPerDay = toNullFloat64(req.UltraFiltrationLitersPerDay)
	wt.PercentageWaterReused = toNullFloat64(req.PercentageWaterReused)
	wt.ElectricityUsedKWH = toNullFloat64(req.ElectricityUsedKWH)
	wt.ChemicalsUsedDescription = toNullString(req.ChemicalsUsedDescription)
	wt.ChemicalsUsedQuantityKG = toNullFloat64(req.ChemicalsUsedQuantityKG)
	wt.Remarks = toNullString(req.Remarks)

	if err := wt.Update(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update water treatment", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Water treatment updated successfully"})
}

func DeleteWaterTreatment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := models.DeleteWaterTreatment(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete water treatment", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Water treatment deleted successfully"})
}

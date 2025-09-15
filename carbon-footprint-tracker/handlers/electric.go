package handlers

import (
	"carbon-footprint-tracker/models"
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type ElectricConsumptionRequest struct {
	Date     time.Time `json:"date"`
	Location string    `json:"location"`
	Source   string    `json:"source" binding:"required"`

	DgCapacityKVA        float64 `json:"dg_capacity_kva"`
	RunningTimeHours     float64 `json:"running_time_hours"`
	FuelConsumedLiters   float64 `json:"fuel_consumed_liters"`
	FuelType             string  `json:"fuel_type"`
	EnergyGeneratedDgKWH float64 `json:"energy_generated_dg_kwh"`

	GridElectricityUsedKWH    float64 `json:"grid_electricity_used_kwh"`
	ElectricityBillKWH        float64 `json:"electricity_bill_kwh"`
	ElectricityBillCostINR    float64 `json:"electricity_bill_cost_inr"`
	ElectricalAppliancesCount int     `json:"electrical_appliances_count"`

	SolarGeneratedKWH float64 `json:"solar_generated_kwh"`
	Remarks           string  `json:"remarks"`
}

func toNullFloat64(f float64) sql.NullFloat64 {
	if f != 0 {
		return sql.NullFloat64{Float64: f, Valid: true}
	}
	return sql.NullFloat64{}
}

func toNullInt32(i int) sql.NullInt32 {
	if i != 0 {
		return sql.NullInt32{Int32: int32(i), Valid: true}
	}
	return sql.NullInt32{}
}

func toNullString(s string) sql.NullString {
	if s != "" {
		return sql.NullString{String: s, Valid: true}
	}
	return sql.NullString{}
}

func toNullBool(b bool) sql.NullBool {
	return sql.NullBool{Bool: b, Valid: true}
}

func AddElectricConsumption(c *gin.Context) {
	var req ElectricConsumptionRequest
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

	e := models.ElectricConsumption{
		Date:                      req.Date,
		Location:                  req.Location,
		Source:                    req.Source,
		DgCapacityKVA:             toNullFloat64(req.DgCapacityKVA),
		RunningTimeHours:          toNullFloat64(req.RunningTimeHours),
		FuelConsumedLiters:        toNullFloat64(req.FuelConsumedLiters),
		FuelType:                  toNullString(req.FuelType),
		EnergyGeneratedDgKWH:      toNullFloat64(req.EnergyGeneratedDgKWH),
		GridElectricityUsedKWH:    toNullFloat64(req.GridElectricityUsedKWH),
		ElectricityBillKWH:        toNullFloat64(req.ElectricityBillKWH),
		ElectricityBillCostINR:    toNullFloat64(req.ElectricityBillCostINR),
		ElectricalAppliancesCount: toNullInt32(req.ElectricalAppliancesCount),
		SolarGeneratedKWH:         toNullFloat64(req.SolarGeneratedKWH),
		Remarks:                   toNullString(req.Remarks),
	}

	if err := e.Create(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add electric consumption", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Electric consumption added successfully", "id": e.ID})
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

	var req ElectricConsumptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	e.Date = req.Date
	if req.Location != "" {
		e.Location = req.Location
	}
	if req.Source != "" {
		e.Source = req.Source
	}
	e.DgCapacityKVA = toNullFloat64(req.DgCapacityKVA)
	e.RunningTimeHours = toNullFloat64(req.RunningTimeHours)
	e.FuelConsumedLiters = toNullFloat64(req.FuelConsumedLiters)
	e.FuelType = toNullString(req.FuelType)
	e.EnergyGeneratedDgKWH = toNullFloat64(req.EnergyGeneratedDgKWH)
	e.GridElectricityUsedKWH = toNullFloat64(req.GridElectricityUsedKWH)
	e.ElectricityBillKWH = toNullFloat64(req.ElectricityBillKWH)
	e.ElectricityBillCostINR = toNullFloat64(req.ElectricityBillCostINR)
	e.ElectricalAppliancesCount = toNullInt32(req.ElectricalAppliancesCount)
	e.SolarGeneratedKWH = toNullFloat64(req.SolarGeneratedKWH)
	e.Remarks = toNullString(req.Remarks)

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

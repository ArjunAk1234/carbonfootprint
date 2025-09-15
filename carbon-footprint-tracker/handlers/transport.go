package handlers

import (
	"carbon-footprint-tracker/models"
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type TransportRequest struct {
	Date                     time.Time `json:"date"`
	EventAreaLocation        string    `json:"event_area_location"`
	VehicleType              string    `json:"vehicle_type" binding:"required"`
	FuelType                 string    `json:"fuel_type" binding:"required"`
	VehicleNumber            string    `json:"vehicle_number"`
	StartLocation            string    `json:"start_location"`
	EndLocation              string    `json:"end_location"`
	DistanceKM               float64   `json:"distance_km" binding:"required"`
	FuelLiters               float64   `json:"fuel_liters" binding:"required"`
	PeopleTravelledCount     int       `json:"people_travelled_count"`
	FuelEfficiencyKMPerLiter float64   `json:"fuel_efficiency_km_per_liter"`
	Remarks                  string    `json:"remarks"`
}

func AddTransportData(c *gin.Context) {
	var req TransportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.EventAreaLocation == "" {
		req.EventAreaLocation = "Overall"
	}
	if req.Date.IsZero() {
		req.Date = time.Now()
	}

	t := models.Transport{
		Date:                     req.Date,
		EventAreaLocation:        req.EventAreaLocation,
		VehicleType:              req.VehicleType,
		FuelType:                 req.FuelType,
		VehicleNumber:            toNullString(req.VehicleNumber),
		StartLocation:            toNullString(req.StartLocation),
		EndLocation:              toNullString(req.EndLocation),
		DistanceKM:               req.DistanceKM,
		FuelLiters:               req.FuelLiters,
		PeopleTravelledCount:     toNullInt32(req.PeopleTravelledCount),
		FuelEfficiencyKMPerLiter: toNullFloat64(req.FuelEfficiencyKMPerLiter),
		Remarks:                  toNullString(req.Remarks),
	}

	if err := t.Create(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add transport data", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Transport data added successfully", "id": t.ID})
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

	var req TransportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	t.Date = req.Date
	if req.EventAreaLocation != "" {
		t.EventAreaLocation = req.EventAreaLocation
	}
	if req.VehicleType != "" {
		t.VehicleType = req.VehicleType
	}
	if req.FuelType != "" {
		t.FuelType = req.FuelType
	}
	t.VehicleNumber = toNullString(req.VehicleNumber)
	t.StartLocation = toNullString(req.StartLocation)
	t.EndLocation = toNullString(req.EndLocation)
	if req.DistanceKM != 0 {
		t.DistanceKM = req.DistanceKM
	}
	if req.FuelLiters != 0 {
		t.FuelLiters = req.FuelLiters
	}
	t.PeopleTravelledCount = toNullInt32(req.PeopleTravelledCount)
	t.FuelEfficiencyKMPerLiter = toNullFloat64(req.FuelEfficiencyKMPerLiter)
	t.Remarks = toNullString(req.Remarks)

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

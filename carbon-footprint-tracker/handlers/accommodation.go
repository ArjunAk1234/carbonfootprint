package handlers

import (
	"carbon-footprint-tracker/models"
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type AccommodationRequest struct {
	Date                      time.Time `json:"date"`
	ParticipantGuestName      string    `json:"participant_guest_name"`
	Category                  string    `json:"category"`
	PeopleCount               int       `json:"people_count" binding:"required"`
	AccommodationFacilityName string    `json:"accommodation_facility_name"`
	AccommodationType         string    `json:"accommodation_type"`
	RoomType                  string    `json:"room_type"`
	NoOfRooms                 int       `json:"no_of_rooms"`
	Nights                    int       `json:"nights" binding:"required"`
	ElectricityConsumptionKWH float64   `json:"electricity_consumption_kwh"`
	WaterConsumptionLPD       float64   `json:"water_consumption_lpd"`
	MealsProvided             bool      `json:"meals_provided"`
	TransportModeToVenue      string    `json:"transport_mode_to_venue"`
	Remarks                   string    `json:"remarks"`
}

// AddAccommodationData handles adding new accommodation data
func AddAccommodationData(c *gin.Context) {
	var req AccommodationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.AccommodationFacilityName == "" {
		req.AccommodationFacilityName = "Overall"
	}
	if req.Date.IsZero() {
		req.Date = time.Now()
	}

	a := models.Accommodation{
		Date:                      req.Date,
		ParticipantGuestName:      toNullString(req.ParticipantGuestName),
		Category:                  toNullString(req.Category),
		PeopleCount:               req.PeopleCount,
		AccommodationFacilityName: req.AccommodationFacilityName,
		AccommodationType:         toNullString(req.AccommodationType),
		RoomType:                  toNullString(req.RoomType),
		NoOfRooms:                 toNullInt32(req.NoOfRooms),
		Nights:                    req.Nights,
		ElectricityConsumptionKWH: toNullFloat64(req.ElectricityConsumptionKWH),
		WaterConsumptionLPD:       toNullFloat64(req.WaterConsumptionLPD),
		MealsProvided:             toNullBool(req.MealsProvided),
		TransportModeToVenue:      toNullString(req.TransportModeToVenue),
		Remarks:                   toNullString(req.Remarks),
	}

	if err := a.Create(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add accommodation data", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Accommodation data added successfully", "id": a.ID})
}

// GetAccommodationData retrieves all accommodation data
func GetAccommodationData(c *gin.Context) {
	accommodations, err := models.GetAllAccommodations()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve accommodation data", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, accommodations)
}

// UpdateAccommodationData updates existing accommodation data
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

	var req AccommodationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	a.Date = req.Date
	a.ParticipantGuestName = toNullString(req.ParticipantGuestName)
	a.Category = toNullString(req.Category)
	if req.PeopleCount != 0 {
		a.PeopleCount = req.PeopleCount
	}
	if req.AccommodationFacilityName != "" {
		a.AccommodationFacilityName = req.AccommodationFacilityName
	}
	a.AccommodationType = toNullString(req.AccommodationType)
	a.RoomType = toNullString(req.RoomType)
	a.NoOfRooms = toNullInt32(req.NoOfRooms)
	if req.Nights != 0 {
		a.Nights = req.Nights
	}
	a.ElectricityConsumptionKWH = toNullFloat64(req.ElectricityConsumptionKWH)
	a.WaterConsumptionLPD = toNullFloat64(req.WaterConsumptionLPD)
	a.MealsProvided = toNullBool(req.MealsProvided)
	a.TransportModeToVenue = toNullString(req.TransportModeToVenue)
	a.Remarks = toNullString(req.Remarks)

	if err := a.Update(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update accommodation data", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Accommodation data updated successfully"})
}

// DeleteAccommodationData deletes accommodation data
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

package handlers

import (
	"carbon-footprint-tracker/models"
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type FoodConsumptionRequest struct {
	Date                     time.Time `json:"date"`
	Location                 string    `json:"location" binding:"required"`
	FoodItem                 string    `json:"food_item" binding:"required"`
	QuantityCookedKgLiter    float64   `json:"quantity_cooked_kg_liter" binding:"required"`
	NoOfMealsServed          int       `json:"no_of_meals_served"`
	RawMaterialSource        string    `json:"raw_material_source"`
	WaterUsedLWashingCooking float64   `json:"water_used_l_washing_cooking"`
	FuelUsedType             string    `json:"fuel_used_type"`
	FuelUsedQuantity         float64   `json:"fuel_used_quantity"`
	Remarks                  string    `json:"remarks"`
}

func AddFoodConsumption(c *gin.Context) {
	var req FoodConsumptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Date.IsZero() {
		req.Date = time.Now()
	}

	f := models.FoodConsumption{
		Date:                     req.Date,
		Location:                 req.Location,
		FoodItem:                 req.FoodItem,
		QuantityCookedKgLiter:    req.QuantityCookedKgLiter,
		NoOfMealsServed:          toNullInt32(req.NoOfMealsServed),
		RawMaterialSource:        toNullString(req.RawMaterialSource),
		WaterUsedLWashingCooking: toNullFloat64(req.WaterUsedLWashingCooking),
		FuelUsedType:             toNullString(req.FuelUsedType),
		FuelUsedQuantity:         toNullFloat64(req.FuelUsedQuantity),
		Remarks:                  toNullString(req.Remarks),
	}

	if err := f.Create(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add food consumption", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Food consumption added successfully", "id": f.ID})
}

func GetFoodConsumptions(c *gin.Context) {
	consumptions, err := models.GetAllFoodConsumptions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve food consumptions", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, consumptions)
}

func UpdateFoodConsumption(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	f, err := models.GetFoodConsumptionByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Food consumption entry not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve food consumption", "details": err.Error()})
		return
	}

	var req FoodConsumptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	f.Date = req.Date
	if req.Location != "" {
		f.Location = req.Location
	}
	if req.FoodItem != "" {
		f.FoodItem = req.FoodItem
	}
	if req.QuantityCookedKgLiter != 0 {
		f.QuantityCookedKgLiter = req.QuantityCookedKgLiter
	}
	f.NoOfMealsServed = toNullInt32(req.NoOfMealsServed)
	f.RawMaterialSource = toNullString(req.RawMaterialSource)
	f.WaterUsedLWashingCooking = toNullFloat64(req.WaterUsedLWashingCooking)
	f.FuelUsedType = toNullString(req.FuelUsedType)
	f.FuelUsedQuantity = toNullFloat64(req.FuelUsedQuantity)
	f.Remarks = toNullString(req.Remarks)

	if err := f.Update(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update food consumption", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Food consumption updated successfully"})
}

func DeleteFoodConsumption(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := models.DeleteFoodConsumption(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete food consumption", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Food consumption deleted successfully"})
}

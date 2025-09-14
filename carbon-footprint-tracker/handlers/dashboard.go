package handlers

import (
	"carbon-footprint-tracker/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Simplified Emission Factors (PLACEHOLDER VALUES - USE ACTUAL FACTORS FOR ACCURATE CALCULATIONS)
const (
	// Electricity grid emissions (kg CO2e / kWh) - Varies heavily by region
	EmissionFactorElectricity = 0.4 // Example: average for some grids
	// Diesel fuel (kg CO2e / liter)
	EmissionFactorDiesel = 2.68
	// Petrol/Gasoline fuel (kg CO2e / liter)
	EmissionFactorPetrol = 2.31
	// Water consumption (kg CO2e / liter) - for treatment and distribution
	EmissionFactorWater = 0.00034 // Example: 0.34 kg CO2e per m³ (1000 liters)
	// Waste (kg CO2e / kg) - highly dependent on waste type and disposal method
	EmissionFactorWasteBiodegradable = 0.1  // Example: simple composting
	EmissionFactorWasteRecyclable    = -0.1 // Example: positive impact from recycling
	EmissionFactorWasteLandfill      = 0.5  // Example: includes methane from decomposition
)

type DashboardData struct {
	TotalCarbonFootprint float64                `json:"total_carbon_footprint_co2e"`
	ComponentBreakdown   map[string]float64     `json:"component_breakdown"`
	PerCapitaFootprint   float64                `json:"per_capita_footprint_co2e"`
	TotalPopulation      int                    `json:"total_population"`
	Trends               map[string]interface{} `json:"trends,omitempty"` // Placeholder for trends
}

// GetDashboardSummary provides a summarized view of the carbon footprint
func GetDashboardSummary(c *gin.Context) {
	totalCarbonFootprint := 0.0
	componentBreakdown := make(map[string]float64)
	totalPopulation := 0

	// 1. Electrical Consumption
	electricConsumptions, err := models.GetAllElectricConsumptions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get electric consumption data for dashboard", "details": err.Error()})
		return
	}
	electricFootprint := 0.0
	for _, ec := range electricConsumptions {
		if ec.Source == "main board" && ec.KWH > 0 {
			electricFootprint += ec.KWH * EmissionFactorElectricity
		} else if ec.Source == "generator" && ec.FuelLiters > 0 {
			// Assuming diesel generators for simplicity
			electricFootprint += ec.FuelLiters * EmissionFactorDiesel
		}
	}
	componentBreakdown["Electrical"] = electricFootprint
	totalCarbonFootprint += electricFootprint

	// 2. Population (for per-capita)
	populations, err := models.GetAllPopulations()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get population data for dashboard", "details": err.Error()})
		return
	}
	// For simplicity, aggregate all population entries for total event population.
	// A real-world scenario might take the max, or sum over a period for duration-based events.
	for _, p := range populations {
		totalPopulation += p.RegisteredCount + p.FloatingCount
	}

	// 3. Transport
	transports, err := models.GetAllTransports()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get transport data for dashboard", "details": err.Error()})
		return
	}
	transportFootprint := 0.0
	for _, t := range transports {
		if t.FuelLiters > 0 {
			if t.FuelType == "Diesel" {
				transportFootprint += t.FuelLiters * EmissionFactorDiesel
			} else if t.FuelType == "Petrol" || t.FuelType == "Gasoline" {
				transportFootprint += t.FuelLiters * EmissionFactorPetrol
			}
			// Add more fuel types as needed
		}
	}
	componentBreakdown["Transport"] = transportFootprint
	totalCarbonFootprint += transportFootprint

	// 4. Water Consumption
	waterConsumptions, err := models.GetAllWaterConsumptions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get water consumption data for dashboard", "details": err.Error()})
		return
	}
	waterFootprint := 0.0
	for _, wc := range waterConsumptions {
		if wc.MeterReading > 0 {
			waterFootprint += wc.MeterReading * EmissionFactorWater
		}
	}
	componentBreakdown["Water"] = waterFootprint
	totalCarbonFootprint += waterFootprint

	// 5. Waste Generation
	wasteEntries, err := models.GetAllWasteEntries()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get waste data for dashboard", "details": err.Error()})
		return
	}
	wasteFootprint := 0.0
	for _, w := range wasteEntries {
		if w.WeightKG > 0 {
			switch w.WasteType {
			case "biodegradable":
				wasteFootprint += w.WeightKG * EmissionFactorWasteBiodegradable
			case "recyclable":
				wasteFootprint += w.WeightKG * EmissionFactorWasteRecyclable
			case "landfill":
				wasteFootprint += w.WeightKG * EmissionFactorWasteLandfill
			}
		}
	}
	componentBreakdown["Waste"] = wasteFootprint
	totalCarbonFootprint += wasteFootprint

	// 6. Accommodation (simplified - assuming direct energy/water already covered or negligible separate footprint)
	// For simplicity, accommodation itself is not directly calculating CO2e here, as its impact
	// largely comes from electricity, water, and waste, which are covered above.
	// If specific factors for accommodation (e.g., heating per night) were available, they would be added here.
	componentBreakdown["Accommodation"] = 0.0

	// 7. Goods Purchased (embodied carbon - highly complex, often estimated)
	// This is a placeholder as embodied carbon calculation is typically complex and requires detailed product LCI data.
	// For a simple example, we might assume a very rough factor per cost or quantity if no other data exists.
	// goodsPurchased, err := models.GetAllGoodsPurchased()
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get goods purchased data for dashboard", "details": err.Error()})
	// 	return
	// }
	componentBreakdown["Goods Purchased"] = 0.0 // Placeholder for now

	perCapitaFootprint := 0.0
	if totalPopulation > 0 {
		perCapitaFootprint = totalCarbonFootprint / float64(totalPopulation)
	}

	// Placeholder for trends over time
	// This would involve querying data for specific date ranges and aggregating.
	trends := map[string]interface{}{
		"daily_carbon": []interface{}{}, // e.g., [{"date": "2024-01-01", "carbon": 100}, ...]
	}

	c.JSON(http.StatusOK, DashboardData{
		TotalCarbonFootprint: totalCarbonFootprint,
		ComponentBreakdown:   componentBreakdown,
		PerCapitaFootprint:   perCapitaFootprint,
		TotalPopulation:      totalPopulation,
		Trends:               trends,
	})
}

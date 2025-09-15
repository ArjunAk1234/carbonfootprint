package handlers

import (
	"carbon-footprint-tracker/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	EmissionFactorGridElectricity = 0.45

	EmissionFactorDiesel = 2.68

	EmissionFactorPetrol = 2.31

	EmissionFactorBiofuel = 1.0

	EmissionFactorWaterTreatmentDist = 0.00034

	EmissionFactorWasteBiodegradable = 0.1
	EmissionFactorWasteRecyclable    = -0.1
	EmissionFactorWasteLandfill      = 0.5
	EmissionFactorWasteEwaste        = 2.0
	EmissionFactorGoodsCost          = 0.05

	EmissionFactorFoodWater = 0.0001

	EmissionFactorLPG = 2.98

	EmissionFactorFirewood  = 1.8
	EmissionFactorChemicals = 1.0
)

type DashboardData struct {
	TotalCarbonFootprint float64                `json:"total_carbon_footprint_co2e"`
	ComponentBreakdown   map[string]float64     `json:"component_breakdown"`
	PerCapitaFootprint   float64                `json:"per_capita_footprint_co2e"`
	TotalPopulation      int                    `json:"total_population"`
	Trends               map[string]interface{} `json:"trends,omitempty"`
}

func GetDashboardSummary(c *gin.Context) {
	totalCarbonFootprint := 0.0
	componentBreakdown := make(map[string]float64)
	totalPopulation := 0

	//1. Electrical Consumption
	electricConsumptions, err := models.GetAllElectricConsumptions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get electric consumption data for dashboard", "details": err.Error()})
		return
	}
	electricFootprint := 0.0
	for _, ec := range electricConsumptions {
		switch ec.Source {
		case "Main Board":
			if ec.GridElectricityUsedKWH.Valid {
				electricFootprint += ec.GridElectricityUsedKWH.Float64 * EmissionFactorGridElectricity
			}
			// Solar generated is an offset, not an emission
			if ec.SolarGeneratedKWH.Valid {
				// This assumes solar generation directly offsets grid consumption.
				// For net calculation: electricFootprint -= ec.SolarGeneratedKWH.Float64 * EmissionFactorGridElectricity
				// For gross accounting, it's typically separate. For simplicity, just track grid usage.
			}
		case "Diesel Generator":
			if ec.FuelConsumedLiters.Valid {
				electricFootprint += ec.FuelConsumedLiters.Float64 * EmissionFactorDiesel
			}
		case "Biofuel Generator":
			if ec.FuelConsumedLiters.Valid {
				electricFootprint += ec.FuelConsumedLiters.Float64 * EmissionFactorBiofuel // Assumed lower factor
			}
			// "Solar Generation" source itself doesn't cause emissions, but represents offset potential
		}
	}
	componentBreakdown["Electrical"] = electricFootprint
	totalCarbonFootprint += electricFootprint

	//2. Population (for per-capita)
	populations, err := models.GetAllPopulations()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get population data for dashboard", "details": err.Error()})
		return
	}
	for _, p := range populations {
		totalPopulation += p.RegisteredCount + p.FloatingCount
	}

	//3. Transport
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
			} else if t.FuelType == "Biofuel" {
				transportFootprint += t.FuelLiters * EmissionFactorBiofuel
			}
		}
	}
	componentBreakdown["Transport"] = transportFootprint
	totalCarbonFootprint += transportFootprint

	//  4. Water Consumption (from usage)
	waterConsumptions, err := models.GetAllWaterConsumptions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get water consumption data for dashboard", "details": err.Error()})
		return
	}
	waterConsumptionFootprint := 0.0
	for _, wc := range waterConsumptions {
		// Assuming total_consumption_kld is daily in KL (1 KL = 1000 L)
		waterConsumptionFootprint += wc.TotalConsumptionKLD * 1000 * EmissionFactorWaterTreatmentDist
	}
	componentBreakdown["Water Consumption"] = waterConsumptionFootprint
	totalCarbonFootprint += waterConsumptionFootprint

	// 5. Water Treatment (from treatment processes)
	waterTreatments, err := models.GetAllWaterTreatments()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get water treatment data for dashboard", "details": err.Error()})
		return
	}
	waterTreatmentFootprint := 0.0
	for _, wt := range waterTreatments {
		if wt.ElectricityUsedKWH.Valid {
			waterTreatmentFootprint += wt.ElectricityUsedKWH.Float64 * EmissionFactorGridElectricity
		}
		if wt.ChemicalsUsedQuantityKG.Valid {
			waterTreatmentFootprint += wt.ChemicalsUsedQuantityKG.Float64 * EmissionFactorChemicals
		}
	}
	componentBreakdown["Water Treatment"] = waterTreatmentFootprint
	totalCarbonFootprint += waterTreatmentFootprint

	//6. Waste Generation
	wasteEntries, err := models.GetAllWasteEntries()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get waste data for dashboard", "details": err.Error()})
		return
	}
	wasteFootprint := 0.0
	for _, w := range wasteEntries {
		if w.WeightKG > 0 {
			switch w.WasteType {
			case "Biodegradable":
				wasteFootprint += w.WeightKG * EmissionFactorWasteBiodegradable
			case "Recyclable":
				wasteFootprint += w.WeightKG * EmissionFactorWasteRecyclable
			case "Landfill":
				wasteFootprint += w.WeightKG * EmissionFactorWasteLandfill
			case "Non-Biodegradable":

				if w.SubCategory.String == "E-waste" {
					wasteFootprint += w.WeightKG * EmissionFactorWasteEwaste
				} else {
					wasteFootprint += w.WeightKG * EmissionFactorWasteLandfill
				}
			}
		}
	}
	componentBreakdown["Waste"] = wasteFootprint
	totalCarbonFootprint += wasteFootprint

	// 7. Accommodation
	accommodations, err := models.GetAllAccommodations()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get accommodation data for dashboard", "details": err.Error()})
		return
	}
	accommodationFootprint := 0.0
	for _, acc := range accommodations {
		if acc.ElectricityConsumptionKWH.Valid {
			accommodationFootprint += acc.ElectricityConsumptionKWH.Float64 * EmissionFactorGridElectricity
		}
		// Assuming water consumption is L/person/day for the duration.
		if acc.WaterConsumptionLPD.Valid {
			// Total L = people_count * nights * L/person/day
			accommodationFootprint += float64(acc.PeopleCount) * float64(acc.Nights) * acc.WaterConsumptionLPD.Float64 * EmissionFactorWaterTreatmentDist
		}

	}
	componentBreakdown["Accommodation"] = accommodationFootprint
	totalCarbonFootprint += accommodationFootprint

	goodsPurchased, err := models.GetAllGoodsPurchased()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get goods purchased data for dashboard", "details": err.Error()})
		return
	}
	goodsFootprint := 0.0
	for _, gp := range goodsPurchased {

		goodsFootprint += gp.BillAmountINR * EmissionFactorGoodsCost

	}
	componentBreakdown["Goods Purchased"] = goodsFootprint
	totalCarbonFootprint += goodsFootprint

	foodConsumptions, err := models.GetAllFoodConsumptions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get food consumption data for dashboard", "details": err.Error()})
		return
	}
	foodFootprint := 0.0
	for _, fc := range foodConsumptions {
		if fc.WaterUsedLWashingCooking.Valid {
			foodFootprint += fc.WaterUsedLWashingCooking.Float64 * EmissionFactorFoodWater
		}
		if fc.FuelUsedType.Valid {
			switch fc.FuelUsedType.String {
			case "LPG":
				if fc.FuelUsedQuantity.Valid {
					foodFootprint += fc.FuelUsedQuantity.Float64 * EmissionFactorLPG
				}
			case "Firewood":
				if fc.FuelUsedQuantity.Valid {
					foodFootprint += fc.FuelUsedQuantity.Float64 * EmissionFactorFirewood
				}
			case "Electricity":
				if fc.FuelUsedQuantity.Valid {
					foodFootprint += fc.FuelUsedQuantity.Float64 * EmissionFactorGridElectricity
				}
			}
		}
		// Embodied carbon of food items themselves is a huge factor but very complex.
		// fc.QuantityCookedKgLiter could be used with specific food-type emission factors.
	}
	componentBreakdown["Food Consumption"] = foodFootprint
	totalCarbonFootprint += foodFootprint

	perCapitaFootprint := 0.0
	if totalPopulation > 0 {
		perCapitaFootprint = totalCarbonFootprint / float64(totalPopulation)
	}

	trends := map[string]interface{}{
		"daily_carbon": []interface{}{},
	}

	c.JSON(http.StatusOK, DashboardData{
		TotalCarbonFootprint: totalCarbonFootprint,
		ComponentBreakdown:   componentBreakdown,
		PerCapitaFootprint:   perCapitaFootprint,
		TotalPopulation:      totalPopulation,
		Trends:               trends,
	})
}

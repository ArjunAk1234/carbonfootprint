package models

import (
	"carbon-footprint-tracker/config"
	"database/sql"
	"time"
)

type FoodConsumption struct {
	ID                       int             `json:"id"`
	Date                     time.Time       `json:"date"`
	Location                 string          `json:"location"` // 'Canteen', 'Hostel', 'Event'
	FoodItem                 string          `json:"food_item"`
	QuantityCookedKgLiter    float64         `json:"quantity_cooked_kg_liter"`
	NoOfMealsServed          sql.NullInt32   `json:"no_of_meals_served,omitempty"`
	RawMaterialSource        sql.NullString  `json:"raw_material_source,omitempty"` // 'Local', 'Market', 'Imported'
	WaterUsedLWashingCooking sql.NullFloat64 `json:"water_used_l_washing_cooking,omitempty"`
	FuelUsedType             sql.NullString  `json:"fuel_used_type,omitempty"` // 'LPG', 'Firewood', 'Electricity'
	FuelUsedQuantity         sql.NullFloat64 `json:"fuel_used_quantity,omitempty"`
	Remarks                  sql.NullString  `json:"remarks,omitempty"`
}

func (f *FoodConsumption) Create() error {
	query := `INSERT INTO food_consumption (
		date, location, food_item, quantity_cooked_kg_liter, no_of_meals_served,
		raw_material_source, water_used_l_washing_cooking, fuel_used_type,
		fuel_used_quantity, remarks
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id`

	return config.DB.QueryRow(query,
		f.Date, f.Location, f.FoodItem, f.QuantityCookedKgLiter, f.NoOfMealsServed,
		f.RawMaterialSource, f.WaterUsedLWashingCooking, f.FuelUsedType,
		f.FuelUsedQuantity, f.Remarks,
	).Scan(&f.ID)
}

func GetAllFoodConsumptions() ([]FoodConsumption, error) {
	rows, err := config.DB.Query(`SELECT
		id, date, location, food_item, quantity_cooked_kg_liter, no_of_meals_served,
		raw_material_source, water_used_l_washing_cooking, fuel_used_type,
		fuel_used_quantity, remarks
		FROM food_consumption ORDER BY date DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var consumptions []FoodConsumption
	for rows.Next() {
		f := FoodConsumption{}
		err := rows.Scan(
			&f.ID, &f.Date, &f.Location, &f.FoodItem, &f.QuantityCookedKgLiter, &f.NoOfMealsServed,
			&f.RawMaterialSource, &f.WaterUsedLWashingCooking, &f.FuelUsedType,
			&f.FuelUsedQuantity, &f.Remarks,
		)
		if err != nil {
			return nil, err
		}
		consumptions = append(consumptions, f)
	}
	return consumptions, nil
}

func GetFoodConsumptionByID(id int) (*FoodConsumption, error) {
	f := &FoodConsumption{}
	query := `SELECT
		id, date, location, food_item, quantity_cooked_kg_liter, no_of_meals_served,
		raw_material_source, water_used_l_washing_cooking, fuel_used_type,
		fuel_used_quantity, remarks
		FROM food_consumption WHERE id = $1`
	err := config.DB.QueryRow(query, id).Scan(
		&f.ID, &f.Date, &f.Location, &f.FoodItem, &f.QuantityCookedKgLiter, f.NoOfMealsServed,
		&f.RawMaterialSource, &f.WaterUsedLWashingCooking, &f.FuelUsedType,
		&f.FuelUsedQuantity, &f.Remarks,
	)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (f *FoodConsumption) Update() error {
	query := `UPDATE food_consumption SET
		date=$1, location=$2, food_item=$3, quantity_cooked_kg_liter=$4, no_of_meals_served=$5,
		raw_material_source=$6, water_used_l_washing_cooking=$7, fuel_used_type=$8,
		fuel_used_quantity=$9, remarks=$10
		WHERE id=$11`
	_, err := config.DB.Exec(query,
		f.Date, f.Location, f.FoodItem, f.QuantityCookedKgLiter, f.NoOfMealsServed,
		f.RawMaterialSource, f.WaterUsedLWashingCooking, f.FuelUsedType,
		f.FuelUsedQuantity, f.Remarks, f.ID,
	)
	return err
}

func DeleteFoodConsumption(id int) error {
	query := `DELETE FROM food_consumption WHERE id=$1`
	_, err := config.DB.Exec(query, id)
	return err
}

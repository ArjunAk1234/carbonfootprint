package models

import (
	"carbon-footprint-tracker/config"
	"database/sql"
	"time"
)

type WaterTreatment struct {
	ID                          int             `json:"id"`
	Date                        time.Time       `json:"date"`
	Location                    string          `json:"location"`
	TreatedLitersPerDay         sql.NullFloat64 `json:"treated_liters_per_day,omitempty"`
	UltraFiltrationLitersPerDay sql.NullFloat64 `json:"ultra_filtration_liters_per_day,omitempty"`
	PercentageWaterReused       sql.NullFloat64 `json:"percentage_water_reused,omitempty"`
	ElectricityUsedKWH          sql.NullFloat64 `json:"electricity_used_kwh,omitempty"`
	ChemicalsUsedDescription    sql.NullString  `json:"chemicals_used_description,omitempty"`
	ChemicalsUsedQuantityKG     sql.NullFloat64 `json:"chemicals_used_quantity_kg,omitempty"`
	Remarks                     sql.NullString  `json:"remarks,omitempty"`
}

func (wt *WaterTreatment) Create() error {
	query := `INSERT INTO water_treatment (
		date, location, treated_liters_per_day, ultra_filtration_liters_per_day,
		percentage_water_reused, electricity_used_kwh, chemicals_used_description,
		chemicals_used_quantity_kg, remarks
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`

	return config.DB.QueryRow(query,
		wt.Date, wt.Location, wt.TreatedLitersPerDay, wt.UltraFiltrationLitersPerDay,
		wt.PercentageWaterReused, wt.ElectricityUsedKWH, wt.ChemicalsUsedDescription,
		wt.ChemicalsUsedQuantityKG, wt.Remarks,
	).Scan(&wt.ID)
}

func GetAllWaterTreatments() ([]WaterTreatment, error) {
	rows, err := config.DB.Query(`SELECT
		id, date, location, treated_liters_per_day, ultra_filtration_liters_per_day,
		percentage_water_reused, electricity_used_kwh, chemicals_used_description,
		chemicals_used_quantity_kg, remarks
		FROM water_treatment ORDER BY date DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var treatments []WaterTreatment
	for rows.Next() {
		wt := WaterTreatment{}
		err := rows.Scan(
			&wt.ID, &wt.Date, &wt.Location, &wt.TreatedLitersPerDay, &wt.UltraFiltrationLitersPerDay,
			&wt.PercentageWaterReused, &wt.ElectricityUsedKWH, &wt.ChemicalsUsedDescription,
			&wt.ChemicalsUsedQuantityKG, &wt.Remarks,
		)
		if err != nil {
			return nil, err
		}
		treatments = append(treatments, wt)
	}
	return treatments, nil
}

func GetWaterTreatmentByID(id int) (*WaterTreatment, error) {
	wt := &WaterTreatment{}
	query := `SELECT
		id, date, location, treated_liters_per_day, ultra_filtration_liters_per_day,
		percentage_water_reused, electricity_used_kwh, chemicals_used_description,
		chemicals_used_quantity_kg, remarks
		FROM water_treatment WHERE id = $1`
	err := config.DB.QueryRow(query, id).Scan(
		&wt.ID, &wt.Date, &wt.Location, &wt.TreatedLitersPerDay, &wt.UltraFiltrationLitersPerDay,
		&wt.PercentageWaterReused, &wt.ElectricityUsedKWH, &wt.ChemicalsUsedDescription,
		&wt.ChemicalsUsedQuantityKG, &wt.Remarks,
	)
	if err != nil {
		return nil, err
	}
	return wt, nil
}

func (wt *WaterTreatment) Update() error {
	query := `UPDATE water_treatment SET
		date=$1, location=$2, treated_liters_per_day=$3, ultra_filtration_liters_per_day=$4,
		percentage_water_reused=$5, electricity_used_kwh=$6, chemicals_used_description=$7,
		chemicals_used_quantity_kg=$8, remarks=$9
		WHERE id=$10`
	_, err := config.DB.Exec(query,
		wt.Date, wt.Location, wt.TreatedLitersPerDay, wt.UltraFiltrationLitersPerDay,
		wt.PercentageWaterReused, wt.ElectricityUsedKWH, wt.ChemicalsUsedDescription,
		wt.ChemicalsUsedQuantityKG, wt.Remarks, wt.ID,
	)
	return err
}

func DeleteWaterTreatment(id int) error {
	query := `DELETE FROM water_treatment WHERE id=$1`
	_, err := config.DB.Exec(query, id)
	return err
}

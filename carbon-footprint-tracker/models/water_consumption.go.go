package models

import (
	"carbon-footprint-tracker/config"
	"database/sql"
	"time"
)

type WaterConsumption struct {
	ID                      int             `json:"id"`
	Date                    time.Time       `json:"date"`
	Location                string          `json:"location"`
	WaterSource             sql.NullString  `json:"water_source,omitempty"`
	CumulativeMeterReading  sql.NullFloat64 `json:"cumulative_meter_reading,omitempty"`
	TotalConsumptionKLD     float64         `json:"total_consumption_kld"`
	PerCapitaConsumptionLPD sql.NullFloat64 `json:"per_capita_consumption_lpd,omitempty"`
	UsageType               sql.NullString  `json:"usage_type,omitempty"`
	Remarks                 sql.NullString  `json:"remarks,omitempty"`
}

func (w *WaterConsumption) Create() error {
	query := `INSERT INTO water_consumption (
		date, location, water_source, cumulative_meter_reading, total_consumption_kld,
		per_capita_consumption_lpd, usage_type, remarks
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`
	return config.DB.QueryRow(query,
		w.Date, w.Location, w.WaterSource, w.CumulativeMeterReading, w.TotalConsumptionKLD,
		w.PerCapitaConsumptionLPD, w.UsageType, w.Remarks,
	).Scan(&w.ID)
}

func GetAllWaterConsumptions() ([]WaterConsumption, error) {
	rows, err := config.DB.Query(`SELECT
		id, date, location, water_source, cumulative_meter_reading, total_consumption_kld,
		per_capita_consumption_lpd, usage_type, remarks
		FROM water_consumption ORDER BY date DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var consumptions []WaterConsumption
	for rows.Next() {
		w := WaterConsumption{}
		err := rows.Scan(
			&w.ID, &w.Date, &w.Location, &w.WaterSource, &w.CumulativeMeterReading, &w.TotalConsumptionKLD,
			&w.PerCapitaConsumptionLPD, &w.UsageType, &w.Remarks,
		)
		if err != nil {
			return nil, err
		}
		consumptions = append(consumptions, w)
	}
	return consumptions, nil
}

func GetWaterConsumptionByID(id int) (*WaterConsumption, error) {
	w := &WaterConsumption{}
	query := `SELECT
		id, date, location, water_source, cumulative_meter_reading, total_consumption_kld,
		per_capita_consumption_lpd, usage_type, remarks
		FROM water_consumption WHERE id = $1`
	err := config.DB.QueryRow(query, id).Scan(
		&w.ID, &w.Date, &w.Location, &w.WaterSource, &w.CumulativeMeterReading, &w.TotalConsumptionKLD,
		&w.PerCapitaConsumptionLPD, &w.UsageType, &w.Remarks,
	)
	if err != nil {
		return nil, err
	}
	return w, nil
}

func (w *WaterConsumption) Update() error {
	query := `UPDATE water_consumption SET
		date=$1, location=$2, water_source=$3, cumulative_meter_reading=$4, total_consumption_kld=$5,
		per_capita_consumption_lpd=$6, usage_type=$7, remarks=$8
		WHERE id=$9`
	_, err := config.DB.Exec(query,
		w.Date, w.Location, w.WaterSource, w.CumulativeMeterReading, w.TotalConsumptionKLD,
		w.PerCapitaConsumptionLPD, w.UsageType, w.Remarks, w.ID,
	)
	return err
}

func DeleteWaterConsumption(id int) error {
	query := `DELETE FROM water_consumption WHERE id=$1`
	_, err := config.DB.Exec(query, id)
	return err
}

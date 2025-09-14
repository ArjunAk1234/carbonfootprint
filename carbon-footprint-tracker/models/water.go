package models

import (
	"carbon-footprint-tracker/config"
	"time"
)

type WaterConsumption struct {
	ID           int       `json:"id"`
	MeterReading float64   `json:"meter_reading"`
	Date         time.Time `json:"date"`
	Location     string    `json:"location"`
}

func (w *WaterConsumption) Create() error {
	query := `INSERT INTO water_consumption (meter_reading, date, location) VALUES ($1, $2, $3) RETURNING id`
	return config.DB.QueryRow(query, w.MeterReading, w.Date, w.Location).Scan(&w.ID)
}

func GetAllWaterConsumptions() ([]WaterConsumption, error) {
	rows, err := config.DB.Query(`SELECT id, meter_reading, date, location FROM water_consumption ORDER BY date DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var consumptions []WaterConsumption
	for rows.Next() {
		w := WaterConsumption{}
		err := rows.Scan(&w.ID, &w.MeterReading, &w.Date, &w.Location)
		if err != nil {
			return nil, err
		}
		consumptions = append(consumptions, w)
	}
	return consumptions, nil
}

func GetWaterConsumptionByID(id int) (*WaterConsumption, error) {
	w := &WaterConsumption{}
	query := `SELECT id, meter_reading, date, location FROM water_consumption WHERE id = $1`
	err := config.DB.QueryRow(query, id).Scan(&w.ID, &w.MeterReading, &w.Date, &w.Location)
	if err != nil {
		return nil, err
	}
	return w, nil
}

func (w *WaterConsumption) Update() error {
	query := `UPDATE water_consumption SET meter_reading=$1, date=$2, location=$3 WHERE id=$4`
	_, err := config.DB.Exec(query, w.MeterReading, w.Date, w.Location, w.ID)
	return err
}

func DeleteWaterConsumption(id int) error {
	query := `DELETE FROM water_consumption WHERE id=$1`
	_, err := config.DB.Exec(query, id)
	return err
}

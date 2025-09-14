package models

import (
	"carbon-footprint-tracker/config"
	"time"
)

type ElectricConsumption struct {
	ID         int       `json:"id"`
	Source     string    `json:"source"`
	KWH        float64   `json:"kwh,omitempty"`
	FuelLiters float64   `json:"fuel_liters,omitempty"`
	Hours      float64   `json:"hours,omitempty"`
	Date       time.Time `json:"date"`
	Location   string    `json:"location"`
}

func (e *ElectricConsumption) Create() error {
	query := `INSERT INTO electric_consumption (source, kwh, fuel_liters, hours, date, location) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	return config.DB.QueryRow(query, e.Source, e.KWH, e.FuelLiters, e.Hours, e.Date, e.Location).Scan(&e.ID)
}

func GetAllElectricConsumptions() ([]ElectricConsumption, error) {
	rows, err := config.DB.Query(`SELECT id, source, kwh, fuel_liters, hours, date, location FROM electric_consumption ORDER BY date DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var consumptions []ElectricConsumption
	for rows.Next() {
		e := ElectricConsumption{}
		err := rows.Scan(&e.ID, &e.Source, &e.KWH, &e.FuelLiters, &e.Hours, &e.Date, &e.Location)
		if err != nil {
			return nil, err
		}
		consumptions = append(consumptions, e)
	}
	return consumptions, nil
}

func GetElectricConsumptionByID(id int) (*ElectricConsumption, error) {
	e := &ElectricConsumption{}
	query := `SELECT id, source, kwh, fuel_liters, hours, date, location FROM electric_consumption WHERE id = $1`
	err := config.DB.QueryRow(query, id).Scan(&e.ID, &e.Source, &e.KWH, &e.FuelLiters, &e.Hours, &e.Date, &e.Location)
	if err != nil {
		return nil, err
	}
	return e, nil
}

func (e *ElectricConsumption) Update() error {
	query := `UPDATE electric_consumption SET source=$1, kwh=$2, fuel_liters=$3, hours=$4, date=$5, location=$6 WHERE id=$7`
	_, err := config.DB.Exec(query, e.Source, e.KWH, e.FuelLiters, e.Hours, e.Date, e.Location, e.ID)
	return err
}

func DeleteElectricConsumption(id int) error {
	query := `DELETE FROM electric_consumption WHERE id=$1`
	_, err := config.DB.Exec(query, id)
	return err
}

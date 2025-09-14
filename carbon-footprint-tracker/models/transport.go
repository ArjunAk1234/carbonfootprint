package models

import (
	"carbon-footprint-tracker/config"
	"time"
)

type Transport struct {
	ID          int       `json:"id"`
	VehicleType string    `json:"vehicle_type"`
	FuelType    string    `json:"fuel_type"`
	DistanceKM  float64   `json:"distance_km"`
	FuelLiters  float64   `json:"fuel_liters"`
	Date        time.Time `json:"date"`
	Location    string    `json:"location"`
}

func (t *Transport) Create() error {
	query := `INSERT INTO transport (vehicle_type, fuel_type, distance_km, fuel_liters, date, location) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	return config.DB.QueryRow(query, t.VehicleType, t.FuelType, t.DistanceKM, t.FuelLiters, t.Date, t.Location).Scan(&t.ID)
}

func GetAllTransports() ([]Transport, error) {
	rows, err := config.DB.Query(`SELECT id, vehicle_type, fuel_type, distance_km, fuel_liters, date, location FROM transport ORDER BY date DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transports []Transport
	for rows.Next() {
		t := Transport{}
		err := rows.Scan(&t.ID, &t.VehicleType, &t.FuelType, &t.DistanceKM, &t.FuelLiters, &t.Date, &t.Location)
		if err != nil {
			return nil, err
		}
		transports = append(transports, t)
	}
	return transports, nil
}

func GetTransportByID(id int) (*Transport, error) {
	t := &Transport{}
	query := `SELECT id, vehicle_type, fuel_type, distance_km, fuel_liters, date, location FROM transport WHERE id = $1`
	err := config.DB.QueryRow(query, id).Scan(&t.ID, &t.VehicleType, &t.FuelType, &t.DistanceKM, &t.FuelLiters, &t.Date, &t.Location)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (t *Transport) Update() error {
	query := `UPDATE transport SET vehicle_type=$1, fuel_type=$2, distance_km=$3, fuel_liters=$4, date=$5, location=$6 WHERE id=$7`
	_, err := config.DB.Exec(query, t.VehicleType, t.FuelType, t.DistanceKM, t.FuelLiters, t.Date, t.Location, t.ID)
	return err
}

func DeleteTransport(id int) error {
	query := `DELETE FROM transport WHERE id=$1`
	_, err := config.DB.Exec(query, id)
	return err
}

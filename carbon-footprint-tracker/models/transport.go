package models

import (
	"carbon-footprint-tracker/config"
	"database/sql"
	"time"
)

type Transport struct {
	ID                       int             `json:"id"`
	Date                     time.Time       `json:"date"`
	EventAreaLocation        string          `json:"event_area_location"`
	VehicleType              string          `json:"vehicle_type"`
	FuelType                 string          `json:"fuel_type"`
	VehicleNumber            sql.NullString  `json:"vehicle_number,omitempty"`
	StartLocation            sql.NullString  `json:"start_location,omitempty"`
	EndLocation              sql.NullString  `json:"end_location,omitempty"`
	DistanceKM               float64         `json:"distance_km"`
	FuelLiters               float64         `json:"fuel_liters"`
	PeopleTravelledCount     sql.NullInt32   `json:"people_travelled_count,omitempty"`
	FuelEfficiencyKMPerLiter sql.NullFloat64 `json:"fuel_efficiency_km_per_liter,omitempty"`
	Remarks                  sql.NullString  `json:"remarks,omitempty"`
}

func (t *Transport) Create() error {
	query := `INSERT INTO transport (
		date, event_area_location, vehicle_type, fuel_type, vehicle_number, start_location,
		end_location, distance_km, fuel_liters, people_travelled_count, fuel_efficiency_km_per_liter, remarks
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) RETURNING id`
	return config.DB.QueryRow(query,
		t.Date, t.EventAreaLocation, t.VehicleType, t.FuelType, t.VehicleNumber, t.StartLocation,
		t.EndLocation, t.DistanceKM, t.FuelLiters, t.PeopleTravelledCount, t.FuelEfficiencyKMPerLiter, t.Remarks,
	).Scan(&t.ID)
}

func GetAllTransports() ([]Transport, error) {
	rows, err := config.DB.Query(`SELECT
		id, date, event_area_location, vehicle_type, fuel_type, vehicle_number, start_location,
		end_location, distance_km, fuel_liters, people_travelled_count, fuel_efficiency_km_per_liter, remarks
		FROM transport ORDER BY date DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transports []Transport
	for rows.Next() {
		t := Transport{}
		err := rows.Scan(
			&t.ID, &t.Date, &t.EventAreaLocation, &t.VehicleType, &t.FuelType, &t.VehicleNumber,
			&t.StartLocation, &t.EndLocation, &t.DistanceKM, &t.FuelLiters, &t.PeopleTravelledCount,
			&t.FuelEfficiencyKMPerLiter, &t.Remarks,
		)
		if err != nil {
			return nil, err
		}
		transports = append(transports, t)
	}
	return transports, nil
}

func GetTransportByID(id int) (*Transport, error) {
	t := &Transport{}
	query := `SELECT
		id, date, event_area_location, vehicle_type, fuel_type, vehicle_number, start_location,
		end_location, distance_km, fuel_liters, people_travelled_count, fuel_efficiency_km_per_liter, remarks
		FROM transport WHERE id = $1`
	err := config.DB.QueryRow(query, id).Scan(
		&t.ID, &t.Date, &t.EventAreaLocation, &t.VehicleType, &t.FuelType, &t.VehicleNumber,
		&t.StartLocation, &t.EndLocation, &t.DistanceKM, &t.FuelLiters, &t.PeopleTravelledCount,
		&t.FuelEfficiencyKMPerLiter, &t.Remarks,
	)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (t *Transport) Update() error {
	query := `UPDATE transport SET
		date=$1, event_area_location=$2, vehicle_type=$3, fuel_type=$4, vehicle_number=$5, start_location=$6,
		end_location=$7, distance_km=$8, fuel_liters=$9, people_travelled_count=$10, fuel_efficiency_km_per_liter=$11, remarks=$12
		WHERE id=$13`
	_, err := config.DB.Exec(query,
		t.Date, t.EventAreaLocation, t.VehicleType, t.FuelType, t.VehicleNumber, t.StartLocation,
		t.EndLocation, t.DistanceKM, t.FuelLiters, t.PeopleTravelledCount, t.FuelEfficiencyKMPerLiter, t.Remarks, t.ID,
	)
	return err
}

func DeleteTransport(id int) error {
	query := `DELETE FROM transport WHERE id=$1`
	_, err := config.DB.Exec(query, id)
	return err
}

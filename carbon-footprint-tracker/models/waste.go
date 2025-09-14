package models

import (
	"carbon-footprint-tracker/config"
	"time"
)

type Waste struct {
	ID        int       `json:"id"`
	SpotName  string    `json:"spot_name"`
	WasteType string    `json:"waste_type"`
	WeightKG  float64   `json:"weight_kg"`
	Date      time.Time `json:"date"`
	Location  string    `json:"location"`
}

func (w *Waste) Create() error {
	query := `INSERT INTO waste (spot_name, waste_type, weight_kg, date, location) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	return config.DB.QueryRow(query, w.SpotName, w.WasteType, w.WeightKG, w.Date, w.Location).Scan(&w.ID)
}

func GetAllWasteEntries() ([]Waste, error) {
	rows, err := config.DB.Query(`SELECT id, spot_name, waste_type, weight_kg, date, location FROM waste ORDER BY date DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var wastes []Waste
	for rows.Next() {
		w := Waste{}
		err := rows.Scan(&w.ID, &w.SpotName, &w.WasteType, &w.WeightKG, &w.Date, &w.Location)
		if err != nil {
			return nil, err
		}
		wastes = append(wastes, w)
	}
	return wastes, nil
}

func GetWasteByID(id int) (*Waste, error) {
	w := &Waste{}
	query := `SELECT id, spot_name, waste_type, weight_kg, date, location FROM waste WHERE id = $1`
	err := config.DB.QueryRow(query, id).Scan(&w.ID, &w.SpotName, &w.WasteType, &w.WeightKG, &w.Date, &w.Location)
	if err != nil {
		return nil, err
	}
	return w, nil
}

func (w *Waste) Update() error {
	query := `UPDATE waste SET spot_name=$1, waste_type=$2, weight_kg=$3, date=$4, location=$5 WHERE id=$6`
	_, err := config.DB.Exec(query, w.SpotName, w.WasteType, w.WeightKG, w.Date, w.Location, w.ID)
	return err
}

func DeleteWaste(id int) error {
	query := `DELETE FROM waste WHERE id=$1`
	_, err := config.DB.Exec(query, id)
	return err
}

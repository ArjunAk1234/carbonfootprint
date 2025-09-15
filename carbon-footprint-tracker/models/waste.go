package models

import (
	"carbon-footprint-tracker/config"
	"database/sql"
	"time"
)

type Waste struct {
	ID                 int            `json:"id"`
	Date               time.Time      `json:"date"`
	CollectionLocation string         `json:"collection_location"`    // Renamed from Location, from OCR "Location/Building"
	WasteType          string         `json:"waste_type"`             // 'Biodegradable', 'Non-Biodegradable', 'Recyclable', 'Landfill'
	SubCategory        sql.NullString `json:"sub_category,omitempty"` // 'Food', 'Garden', 'Plastic', 'Paper', 'Glass', 'Metal', 'E-waste'
	WeightKG           float64        `json:"weight_kg"`
	CollectionMethod   sql.NullString `json:"collection_method,omitempty"` // 'Bins', 'Direct', 'Vehicle'
	TransportMode      sql.NullString `json:"transport_mode,omitempty"`
	Destination        sql.NullString `json:"destination,omitempty"` // 'Composting', 'Recycler', 'Landfill', 'Incinerator', 'OWC', 'STP'
	Remarks            sql.NullString `json:"remarks,omitempty"`
}

func (w *Waste) Create() error {
	query := `INSERT INTO waste (
		date, collection_location, waste_type, sub_category, weight_kg,
		collection_method, transport_mode, destination, remarks
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`
	return config.DB.QueryRow(query,
		w.Date, w.CollectionLocation, w.WasteType, w.SubCategory, w.WeightKG,
		w.CollectionMethod, w.TransportMode, w.Destination, w.Remarks,
	).Scan(&w.ID)
}

func GetAllWasteEntries() ([]Waste, error) {
	rows, err := config.DB.Query(`SELECT
		id, date, collection_location, waste_type, sub_category, weight_kg,
		collection_method, transport_mode, destination, remarks
		FROM waste ORDER BY date DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var wastes []Waste
	for rows.Next() {
		w := Waste{}
		err := rows.Scan(
			&w.ID, &w.Date, &w.CollectionLocation, &w.WasteType, &w.SubCategory, &w.WeightKG,
			&w.CollectionMethod, &w.TransportMode, &w.Destination, &w.Remarks,
		)
		if err != nil {
			return nil, err
		}
		wastes = append(wastes, w)
	}
	return wastes, nil
}

func GetWasteByID(id int) (*Waste, error) {
	w := &Waste{}
	query := `SELECT
		id, date, collection_location, waste_type, sub_category, weight_kg,
		collection_method, transport_mode, destination, remarks
		FROM waste WHERE id = $1`
	err := config.DB.QueryRow(query, id).Scan(
		&w.ID, &w.Date, &w.CollectionLocation, &w.WasteType, &w.SubCategory, &w.WeightKG,
		&w.CollectionMethod, &w.TransportMode, &w.Destination, &w.Remarks,
	)
	if err != nil {
		return nil, err
	}
	return w, nil
}

func (w *Waste) Update() error {
	query := `UPDATE waste SET
		date=$1, collection_location=$2, waste_type=$3, sub_category=$4, weight_kg=$5,
		collection_method=$6, transport_mode=$7, destination=$8, remarks=$9
		WHERE id=$10`
	_, err := config.DB.Exec(query,
		w.Date, w.CollectionLocation, w.WasteType, w.SubCategory, w.WeightKG,
		w.CollectionMethod, w.TransportMode, w.Destination, w.Remarks, w.ID,
	)
	return err
}

func DeleteWaste(id int) error {
	query := `DELETE FROM waste WHERE id=$1`
	_, err := config.DB.Exec(query, id)
	return err
}
